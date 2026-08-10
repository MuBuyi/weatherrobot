package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"wechatrobot/internal/config"
	"wechatrobot/internal/log"
)

type ForecastResponse struct {
	Error   bool   `json:"error"`
	Reason  string `json:"reason"`
	Current struct {
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		Humidity            int     `json:"relative_humidity_2m"`
		WeatherCode         int     `json:"weather_code"`
		WindSpeed           float64 `json:"wind_speed_10m"`
		WindDirection       int     `json:"wind_direction_10m"`
	} `json:"current"`
	Daily struct {
		Time                     []string  `json:"time"`
		WeatherCode              []int     `json:"weather_code"`
		TempMax                  []float64 `json:"temperature_2m_max"`
		TempMin                  []float64 `json:"temperature_2m_min"`
		PrecipitationProbability []int     `json:"precipitation_probability_max"`
	} `json:"daily"`
}

func weatherHTTPClient() *http.Client {
	return &http.Client{Timeout: time.Duration(config.Cfg.HTTPTimeoutSeconds) * time.Second}
}

func GetForecast() (*ForecastResponse, error) {
	endpoint, err := url.Parse(config.Cfg.WeatherBaseURL + "/v1/forecast")
	if err != nil {
		return nil, fmt.Errorf("天气接口地址无效: %w", err)
	}
	query := endpoint.Query()
	query.Set("latitude", strconv.FormatFloat(config.Cfg.Latitude, 'f', 6, 64))
	query.Set("longitude", strconv.FormatFloat(config.Cfg.Longitude, 'f', 6, 64))
	query.Set("current", "temperature_2m,apparent_temperature,relative_humidity_2m,weather_code,wind_speed_10m,wind_direction_10m")
	query.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max")
	// Daily data starts with today. Request one extra day so the future forecast excludes today.
	query.Set("forecast_days", strconv.Itoa(config.Cfg.ForecastDays+1))
	query.Set("timezone", config.Cfg.Timezone)
	endpoint.RawQuery = query.Encode()

	resp, err := weatherHTTPClient().Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("请求天气接口失败: %w", err)
	}
	defer resp.Body.Close()

	var result ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析天气响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK || result.Error {
		return nil, fmt.Errorf("天气接口返回错误: HTTP %d, %s", resp.StatusCode, result.Reason)
	}
	if len(result.Daily.Time) == 0 {
		return nil, fmt.Errorf("天气接口未返回预报数据")
	}
	return &result, nil
}

func WeatherDescription(code int) string {
	switch {
	case code == 0:
		return "晴"
	case code == 1:
		return "大部晴朗"
	case code == 2:
		return "多云"
	case code == 3:
		return "阴"
	case code == 45 || code == 48:
		return "雾"
	case code >= 51 && code <= 57:
		return "毛毛雨"
	case code >= 61 && code <= 67:
		return "雨"
	case code >= 71 && code <= 77:
		return "雪"
	case code >= 80 && code <= 82:
		return "阵雨"
	case code == 85 || code == 86:
		return "阵雪"
	case code >= 95:
		return "雷暴"
	default:
		return "未知天气"
	}
}

func WindDirection(degrees int) string {
	directions := []string{"北风", "东北风", "东风", "东南风", "南风", "西南风", "西风", "西北风"}
	index := ((degrees + 22) % 360) / 45
	return directions[index]
}

func SendErrorAlert(err error) { log.Error("错误警报: ", err) }
