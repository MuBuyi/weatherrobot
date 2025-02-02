package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg Config

type Config struct {
	WecomWebhook    string   `mapstructure:"wecom_webhook"`
	WeatherAPIKey   string   `mapstructure:"weather_api_key"`
	Locations       []string `mapstructure:"locations"`
	MentionUsers    []string `mapstructure:"mention_users"`
	OffWorkMessages []string `mapstructure:"off_work_messages"` // 新增下班结束语配置项
}

func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/root/yangjian/wechatrobot/config") // 修改为配置文件所在的目录
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("读取配置文件失败: ", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		logrus.Fatal("解析配置失败: ", err)
	}
}
