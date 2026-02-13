package main

import (
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/cronn"
	"wechatrobot/internal/log"
	"wechatrobot/internal/wecom"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	log.Init()

	// 加载配置
	config.Load()

	// 启动企业微信消息接收服务（在后台运行）
	go wecom.StartWecomServer("9001")

	// 创建cron实例（中国时区）
	// c := cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600)))
	c := cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600)))

	// 添加定时任务，每分钟执行一次
	_, err := c.AddFunc("0 8 * * *", cronn.SendDailyReport)
	// _, err := c.AddFunc("*/1 * * * *", cronn.SendDailyReport) // 1分钟执行一次
	if err != nil {
		logrus.Fatal("创建定时任务失败: ", err)
	}
	logrus.Info("天气报告定时任务已添加（每1分钟执行一次）")

	// 添加定时任务，每天下午6点提醒下班
	_, err = c.AddFunc("0 18 * * *", cronn.SendOffWorkReminder)
	// _, err = c.AddFunc("*/1 * * * *", cronn.SendOffWorkReminder) // 1分钟执行一次

	if err != nil {
		logrus.Fatal("创建定时任务失败: ", err)
	}
	logrus.Info("下班提醒定时任务已添加（每1分钟执行一次）")

	// 启动定时任务
	c.Start()
	logrus.Info("天气机器人已启动（包含定时任务和微信交互服务）")

	// 保持程序运行
	select {}
}
