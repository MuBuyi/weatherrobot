package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"wechatrobot/internal/holiday"
)

var Cfg Config

type Config struct {
	WecomWebhook         string            `mapstructure:"wecom_webhook"`
	CityName             string            `mapstructure:"city_name"`
	Latitude             float64           `mapstructure:"latitude"`
	Longitude            float64           `mapstructure:"longitude"`
	MentionUsers         []string          `mapstructure:"mention_users"`
	Holidays             []holiday.Holiday `mapstructure:"holidays"`
	WeatherCron          string            `mapstructure:"weather_cron"`
	LunchCron            string            `mapstructure:"lunch_cron"`
	OffWorkCron          string            `mapstructure:"off_work_cron"`
	Timezone             string            `mapstructure:"timezone"`
	WeatherBaseURL       string            `mapstructure:"weather_base_url"`
	HTTPTimeoutSeconds   int               `mapstructure:"http_timeout_seconds"`
	ForecastDays         int               `mapstructure:"forecast_days"`
	ReminderTitle        string            `mapstructure:"reminder_title"`
	ReminderIntro        string            `mapstructure:"reminder_intro"`
	LunchReminderTitle   string            `mapstructure:"lunch_reminder_title"`
	LunchReminderText    string            `mapstructure:"lunch_reminder_text"`
	OffWorkReminderTitle string            `mapstructure:"off_work_reminder_title"`
	OffWorkReminderText  string            `mapstructure:"off_work_reminder_text"`
	WecomCorpID          string            `mapstructure:"wecom_corp_id"`
	WecomCheckinSecret   string            `mapstructure:"wecom_checkin_secret"`
	WecomCheckinUserIDs  []string          `mapstructure:"wecom_checkin_userids"`
	AttendancePollCron   string            `mapstructure:"attendance_poll_cron"`
	AttendanceStateFile  string            `mapstructure:"attendance_state_file"`
	Footer               string            `mapstructure:"footer"`
}

func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetDefault("weather_cron", "0 8 * * *")
	viper.SetDefault("lunch_cron", "30 11 * * *")
	viper.SetDefault("off_work_cron", "0 18 * * *")
	viper.SetDefault("attendance_poll_cron", "*/1 * * * *")
	viper.SetDefault("attendance_state_file", "data/attendance-state.json")
	viper.SetDefault("timezone", "Asia/Shanghai")
	viper.SetDefault("weather_base_url", "https://api.open-meteo.com")
	viper.SetDefault("http_timeout_seconds", 10)
	viper.SetDefault("forecast_days", 3)
	viper.SetDefault("reminder_title", "🌦️ 武汉市今日天气提醒")
	viper.SetDefault("reminder_intro", "早上好！新的一天开始了，请大家及时关注天气变化，合理安排出行与工作。")
	viper.SetDefault("lunch_reminder_title", "🍱 午餐点餐提醒")
	viper.SetDefault("lunch_reminder_text", "午餐时间快到啦！想好今天吃什么了吗？请大家及时完成点餐，祝大家吃得开心，下午元气满满！")
	viper.SetDefault("off_work_reminder_title", "⏰ 下班打卡提醒")
	viper.SetDefault("off_work_reminder_text", "下班时间到啦，请大家记得完成下班打卡，路上注意安全！")
	viper.SetDefault("footer", "💡 温馨提示：记得关注天气变化哦！")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("读取配置文件失败: ", err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		logrus.Fatal("解析配置失败: ", err)
	}

	Cfg.WecomCorpID = strings.TrimSpace(os.Getenv("WECOM_CORP_ID"))
	Cfg.WecomCheckinSecret = strings.TrimSpace(os.Getenv("WECOM_CHECKIN_SECRET"))
	Cfg.WecomCheckinUserIDs = splitCSV(os.Getenv("WECOM_CHECKIN_USERIDS"))

	Cfg.WeatherBaseURL = strings.TrimRight(Cfg.WeatherBaseURL, "/")
	if Cfg.WecomWebhook == "" || Cfg.CityName == "" {
		logrus.Fatal("wecom_webhook 和 city_name 必须配置")
	}
	if Cfg.Latitude < -90 || Cfg.Latitude > 90 || Cfg.Longitude < -180 || Cfg.Longitude > 180 {
		logrus.Fatal("城市经纬度配置不合法")
	}
	if Cfg.HTTPTimeoutSeconds <= 0 {
		logrus.Fatal("http_timeout_seconds 必须大于 0")
	}
	if Cfg.ForecastDays < 1 || Cfg.ForecastDays > 16 {
		logrus.Fatal("forecast_days 必须在 1 到 16 之间")
	}
	if Cfg.WeatherBaseURL == "" {
		logrus.Fatal("weather_base_url 不能为空")
	}
}

func splitCSV(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	return values
}
