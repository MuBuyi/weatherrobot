package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"wechatrobot/internal/holiday"
)

var Cfg Config

type Config struct {
	WecomWebhook    string           `mapstructure:"wecom_webhook"`
	WeatherAPIKey   string           `mapstructure:"weather_api_key"`
	Locations       []string         `mapstructure:"locations"`
	OpenAIAPIKey    string           `mapstructure:"openai_api_key"`
	UseAIReminder   bool             `mapstructure:"use_ai_reminder"`
	DoubaoURL       string           `mapstructure:"doubao_url"`
	DoubaoAPIKey    string           `mapstructure:"doubao_api_key"`
	DoubaoModel     string           `mapstructure:"doubao_model"`
	MentionUsers    []string         `mapstructure:"mention_users"`
	OffWorkMessages []string         `mapstructure:"off_work_messages"`
	Holidays        []holiday.Holiday `mapstructure:"holidays"` // 用户自定义假期
}

func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("读取配置文件失败: ", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		logrus.Fatal("解析配置失败: ", err)
	}
}
