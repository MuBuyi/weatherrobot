package main

import (
    "wechatrobot/internal/config"
    "wechatrobot/internal/cronn"
    "wechatrobot/internal/log"

    "github.com/sirupsen/logrus"
)

func main() {
    log.Init()
    config.Load()

    logrus.Info("触发测试：开始调用 SendOffWorkReminder")
    cronn.SendOffWorkReminder()
    logrus.Info("触发测试：SendOffWorkReminder 调用完成")
}
