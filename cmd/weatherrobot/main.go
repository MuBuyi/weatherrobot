package main

import (
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/cronn"
	"wechatrobot/internal/log"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	log.Init()

	// 加载配置
	config.Load()

	// 按配置创建定时任务
	location, err := time.LoadLocation(config.Cfg.Timezone)
	if err != nil {
		logrus.Fatal("无效时区: ", err)
	}
	c := cron.New(cron.WithLocation(location))

	// 注册天气日报任务
	_, err = c.AddFunc(config.Cfg.WeatherCron, cronn.SendDailyReport)
	if err != nil {
		logrus.Fatal("创建定时任务失败: ", err)
	}
	logrus.Infof("天气报告定时任务已添加：%s（%s）", config.Cfg.WeatherCron, config.Cfg.Timezone)

	// 注册午餐点餐提醒任务
	_, err = c.AddFunc(config.Cfg.LunchCron, cronn.SendLunchReminder)
	if err != nil {
		logrus.Fatal("创建午餐提醒定时任务失败: ", err)
	}
	logrus.Infof("午餐点餐提醒任务已添加：%s（%s）", config.Cfg.LunchCron, config.Cfg.Timezone)

	// 注册每天 18:00 的固定下班打卡提醒
	_, err = c.AddFunc(config.Cfg.OffWorkCron, cronn.SendOffWorkReminder)
	if err != nil {
		logrus.Fatal("创建下班打卡提醒定时任务失败: ", err)
	}
	logrus.Infof("下班打卡提醒任务已添加：%s（%s）", config.Cfg.OffWorkCron, config.Cfg.Timezone)

	// 启动定时任务
	c.Start()
	logrus.Info("天气预报提醒机器人已启动")

	// 保持程序运行
	select {}
}
