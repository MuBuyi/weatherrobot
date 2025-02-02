package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wechatrobot/internal/config"
	"wechatrobot/internal/log"
)

type WeatherResponse struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Wind360   string `json:"wind360"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		WindSpeed string `json:"windSpeed"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
		Dew       string `json:"dew"`
	} `json:"now"`
	Daily []struct {
		FxDate         string `json:"fxDate"`
		Sunrise        string `json:"sunrise"`
		Sunset         string `json:"sunset"`
		Moonrise       string `json:"moonrise"`
		Moonset        string `json:"moonset"`
		MoonPhase      string `json:"moonPhase"`
		MoonPhaseIcon  string `json:"moonPhaseIcon"`
		TempMax        string `json:"tempMax"`
		TempMin        string `json:"tempMin"`
		IconDay        string `json:"iconDay"`
		TextDay        string `json:"textDay"`
		IconNight      string `json:"iconNight"`
		TextNight      string `json:"textNight"`
		Wind360Day     string `json:"wind360Day"`
		WindDirDay     string `json:"windDirDay"`
		WindScaleDay   string `json:"windScaleDay"`
		WindSpeedDay   string `json:"windSpeedDay"`
		Wind360Night   string `json:"wind360Night"`
		WindDirNight   string `json:"windDirNight"`
		WindScaleNight string `json:"windScaleNight"`
		WindSpeedNight string `json:"windSpeedNight"`
		Humidity       string `json:"humidity"`
		Precip         string `json:"precip"`
		Pressure       string `json:"pressure"`
		Vis            string `json:"vis"`
		Cloud          string `json:"cloud"`
		UvIndex        string `json:"uvIndex"`
	} `json:"daily"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}

type LivingIndicesResponse struct {
	Daily []struct {
		Name string `json:"name"`
		Text string `json:"text"`
	} `json:"daily"`
}

func GetWeather(location, weatherType string) (*WeatherResponse, error) {

	fmt.Println("开始请求天气接口")

	url := fmt.Sprintf("https://api.qweather.com/v7/weather/now?location=%s&key=%s", location, config.Cfg.WeatherAPIKey)

	if weatherType == "7d" {
		url = fmt.Sprintf("https://api.qweather.com/v7/weather/7d?location=%s&key=%s", location, config.Cfg.WeatherAPIKey)

	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 打印 API 响应
	fmt.Println("===API 响应===: ", string(body))

	var weatherResponse WeatherResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return nil, err
	}

	fmt.Println("天气获取成功")

	return &weatherResponse, nil
}

func GetLivingIndices(location string) (*LivingIndicesResponse, error) {

	// https://api.qweather.com/v7/indices/1d?type=1,2&location=101010100
	url := fmt.Sprintf("https://api.qweather.com/v7/indices/1d?type=1,2&key=%s&location=%s", config.Cfg.WeatherAPIKey, location)
	fmt.Println("开始请求生活指数接口")
	fmt.Println("==url===: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var indicesResponse LivingIndicesResponse
	if err := json.Unmarshal(body, &indicesResponse); err != nil {
		return nil, err
	}
	fmt.Println("生活指数获取成功")
	fmt.Println("===indicesResponse===: ", indicesResponse)
	return &indicesResponse, nil
}

func GetCityName(location string) string {
	// 这里可以根据 location 返回对应的城市名称
	// 例如：101020100 -> 上海
	switch location {
	case "101200101":
		return "武汉"
	case "101200805":
		return "监利"
	default:
		return "未知城市"
	}
}

func SendErrorAlert(err error) {
	// 这里可以实现发送错误警报的逻辑，例如通过企业微信发送
	log.Error("错误警报: ", err)
}
