package cronn

import (
	"fmt"
	"time"

	"wechatrobot/internal/config"
	"wechatrobot/internal/holiday"
	"wechatrobot/internal/log"
	"wechatrobot/internal/weather"

	"github.com/sirupsen/logrus"
)

func SendLunchReminder() {
	content := fmt.Sprintf("%s\n%s", config.Cfg.LunchReminderTitle, config.Cfg.LunchReminderText)
	if err := weather.SendWecomMessage(content, config.Cfg.MentionUsers); err != nil {
		log.Error("发送午餐点餐提醒失败: ", err)
		return
	}
	logrus.Info("午餐点餐提醒发送成功")
}

func SendOffWorkReminder() {
	content := fmt.Sprintf("%s\n%s", config.Cfg.OffWorkReminderTitle, config.Cfg.OffWorkReminderText)
	if err := weather.SendWecomMessage(content, config.Cfg.MentionUsers); err != nil {
		log.Error("发送下班打卡提醒失败: ", err)
		return
	}
	logrus.Info("下班打卡提醒发送成功")
}

func SendDailyReport() {
	logrus.Info("天气报告定时任务触发")
	forecast, err := weather.GetForecast()
	if err != nil {
		log.Error("获取天气预报失败: ", err)
		return
	}

	var report string
	report = fmt.Sprintf("%s\n%s\n\n", config.Cfg.ReminderTitle, config.Cfg.ReminderIntro)
	if isFestival, festival := holiday.IsFestival(time.Now()); isFestival && festival != nil {
		report += festival.Greeting + "\n\n"
	}
	today := forecast.Daily.Time[0]
	report += fmt.Sprintf("🏙 城市：%s\n🌤 【%s】当前天气：%s，温度：%.1f℃，体感：%.1f℃\n💧 相对湿度：%d%%\n🎐 风速：%.1f km/h，风向：%s\n\n📅 未来%d日预报：\n",
		config.Cfg.CityName, today, weather.WeatherDescription(forecast.Current.WeatherCode), forecast.Current.Temperature,
		forecast.Current.ApparentTemperature, forecast.Current.Humidity, forecast.Current.WindSpeed,
		weather.WindDirection(forecast.Current.WindDirection), config.Cfg.ForecastDays)

	// Index 0 is today; the future forecast starts from tomorrow.
	count := len(forecast.Daily.Time)
	for i, shown := 1, 0; i < count && shown < config.Cfg.ForecastDays; i, shown = i+1, shown+1 {
		if i >= len(forecast.Daily.WeatherCode) || i >= len(forecast.Daily.TempMin) || i >= len(forecast.Daily.TempMax) {
			break
		}
		rain := 0
		if i < len(forecast.Daily.PrecipitationProbability) {
			rain = forecast.Daily.PrecipitationProbability[i]
		}
		report += fmt.Sprintf("【%s】%s %.1f℃～%.1f℃，降水概率 %d%%\n", forecast.Daily.Time[i],
			weather.WeatherDescription(forecast.Daily.WeatherCode[i]), forecast.Daily.TempMin[i], forecast.Daily.TempMax[i], rain)
	}
	report += "\n" + config.Cfg.Footer
	if err := weather.SendWecomMessage(report, config.Cfg.MentionUsers); err != nil {
		log.Error("发送消息失败: ", err)
		return
	}
	logrus.Info("每日天气报告发送成功")
}
