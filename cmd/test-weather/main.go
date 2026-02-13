package main

import (
	"wechatrobot/internal/config"
	"wechatrobot/internal/cronn"
	"wechatrobot/internal/log"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	log.Init()

	// 加载配置
	config.Load()

	logrus.Info("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	logrus.Info("手动测试天气报告 - 立即发送")
	logrus.Info("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")

	// 直接调用天气报告函数
	cronn.SendDailyReport()

	logrus.Info("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	logrus.Info("测试完毕")
	logrus.Info("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
}
