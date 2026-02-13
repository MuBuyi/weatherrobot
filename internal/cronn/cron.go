package cronn

import (
	"fmt"
	"math/rand"
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/holiday"
	"wechatrobot/internal/ai"
	"wechatrobot/internal/log"
	"wechatrobot/internal/weather"

	"github.com/sirupsen/logrus"
)

func SendDailyReport() {
	logrus.Info("天气报告定时任务触发")
	
	// 检查是否应该发送提醒
	shouldSend, isFestival, festival := holiday.ShouldSendReminder(config.Cfg.Holidays)
	
	if isFestival && festival != nil {
		// 节假日：发送特色问候
		festivalMessage := fmt.Sprintf("%s\n\n除夕至初六期间天气报告已暂停，假期结束后继续为您服务。", festival.Greeting)
		if err := weather.SendWecomMessage(festivalMessage, config.Cfg.MentionUsers); err != nil {
			log.Error("发送节假日问候失败: ", err)
			weather.SendErrorAlert(err)
			return
		}
		logrus.Info("节假日问候发送成功")
		return
	}
	
	if !shouldSend {
		logrus.Info("当前为假期或非工作日，跳过天气报告")
		return
	}
	
	var fullReport string

	// 遍历所有配置的城市
	for _, location := range config.Cfg.Locations {
		// 获取实时天气
		currentWeather, err := weather.GetWeather(location, "current")
		if err != nil {
			log.Error("获取实时天气失败: ", err)
			weather.SendErrorAlert(err)
			continue
		}

		// 获取天气预报
		forecast, err := weather.GetWeather(location, "7d")
		if err != nil {
			log.Error("获取天气预报失败: ", err)
			weather.SendErrorAlert(err)
			continue
		}

		// 获取生活指数
		indices, err := weather.GetLivingIndices(location)
		if err != nil {
			log.Error("获取生活指数失败: ", err)
			weather.SendErrorAlert(err)
			continue
		}

		// 构造城市报告
		cityReport := fmt.Sprintf("🏙 城市：%s\n"+
			"🌤 当前天气：%s，温度：%s℃\n"+
			"🎐 风力：%s级，风向：%s\n\n"+
			"📅 三日预报：\n",
			weather.GetCityName(location),
			currentWeather.Now.Text,
			currentWeather.Now.Temp,
			currentWeather.Now.WindScale,
			currentWeather.Now.WindDir)

		// 添加天气预报
		for i := 0; i < 3 && i < len(forecast.Daily); i++ {
			cityReport += fmt.Sprintf("【%s】%s %s℃～%s℃\n",
				forecast.Daily[i].FxDate,
				forecast.Daily[i].TextDay,
				forecast.Daily[i].TempMin,
				forecast.Daily[i].TempMax)
		}

		// 添加生活指数
		cityReport += "\n📊 生活指数：\n"
		for _, index := range indices.Daily {
			cityReport += fmt.Sprintf("• %s：%s\n", index.Name, index.Text)
		}

		fullReport += cityReport + "\n------------------------\n"
	}

	// 添加结尾
	fullReport += "💡 温馨提示：记得关注天气变化哦！"

	// 发送消息
	if err := weather.SendWecomMessage(fullReport, config.Cfg.MentionUsers); err != nil {
		log.Error("发送消息失败: ", err)
		weather.SendErrorAlert(err)
		return
	}

	logrus.Info("每日天气报告发送成功")
}

// 发送下班提醒
func SendOffWorkReminder() {
	// 检查是否应该发送下班提醒
	shouldSend, _, _ := holiday.ShouldSendOffWorkReminder(config.Cfg.Holidays)
	
	if !shouldSend {
		logrus.Info("当前为假期、节假日或非工作日，跳过下班提醒")
		return
	}
	
	var content string

	// 如果启用了 AI 模式，使用 AI 生成提醒
	if config.Cfg.UseAIReminder {
		generatedMessage, err := ai.GenerateOffWorkReminder(config.Cfg.DoubaoURL, config.Cfg.DoubaoAPIKey, config.Cfg.DoubaoModel, config.Cfg.OpenAIAPIKey)
		if err != nil {
			logrus.Error("使用 AI 生成提醒失败: ", err)
			// 降级到静态文案
			if len(config.Cfg.OffWorkMessages) > 0 {
				content = fmt.Sprintf("提醒：%s", config.Cfg.OffWorkMessages[0])
			} else {
				logrus.Error("生成提醒失败，且没有备用文案")
				return
			}
		} else {
			content = fmt.Sprintf("提醒：%s", generatedMessage)
		}
	} else {
		// 使用静态文案模式
		if len(config.Cfg.OffWorkMessages) == 0 {
			logrus.Error("下班结束语配置为空且未启用 AI 模式")
			return
		}
		rand.Seed(time.Now().Unix())
		message := config.Cfg.OffWorkMessages[rand.Intn(len(config.Cfg.OffWorkMessages))]
		content = fmt.Sprintf("提醒：%s", message)
	}

	if err := weather.SendWecomMessage(content, config.Cfg.MentionUsers); err != nil {
		logrus.Error("发送下班提醒失败: ", err)
		return
	}

	logrus.Info("下班提醒发送成功")
}
