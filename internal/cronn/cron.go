package cronn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/log"
	"wechatrobot/internal/weather"

	"github.com/sirupsen/logrus"
)

func SendDailyReport() {
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

	log.Info("每日天气报告发送成功")
}

// 发送下班提醒
func SendOffWorkReminder() {
	if len(config.Cfg.OffWorkMessages) == 0 {
		logrus.Error("下班结束语配置为空")
		return
	}

	rand.Seed(time.Now().Unix())
	message := config.Cfg.OffWorkMessages[rand.Intn(len(config.Cfg.OffWorkMessages))]
	content := fmt.Sprintf("提醒：%s", message)

	if err := sendWecomMessage(content, config.Cfg.MentionUsers); err != nil {
		logrus.Error("发送下班提醒失败: ", err)
		return
	}

	logrus.Info("下班提醒发送成功")
}

// 发送企业微信消息
func sendWecomMessage(content string, mentionUsers []string) error {
	message := struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string   `json:"content"`
			Mention []string `json:"mentioned_list"`
		} `json:"text"`
	}{
		MsgType: "text",
	}
	message.Text.Content = content
	message.Text.Mention = mentionUsers

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.Cfg.WecomWebhook, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送消息失败，状态码: %d", resp.StatusCode)
	}

	logrus.Info("消息发送成功")
	return nil
}