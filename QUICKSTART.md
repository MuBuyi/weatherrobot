# 天气预报提醒机器人：快速开始

这是一个纯 Go 实现的企业微信天气提醒机器人。程序每天北京时间 08:00 获取配置城市的实时天气、三日预报和生活指数，并通过企业微信群机器人 Webhook 推送。

## 前置条件

- Go 1.18+
- 无需天气 API Key
- 企业微信群机器人 Webhook

## 配置

创建 `config/config.yaml`：

```yaml
wecom_webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY"
city_name: "武汉市"
latitude: 30.5928
longitude: 114.3055
weather_cron: "0 8 * * *"
timezone: "Asia/Shanghai"
weather_base_url: "https://api.open-meteo.com"
http_timeout_seconds: 10
forecast_days: 3
reminder_title: "🌦️ 武汉市今日天气提醒"
reminder_intro: "早上好！新的一天开始了，请大家及时关注天气变化，合理安排出行与工作。"
footer: "💡 温馨提示：记得关注天气变化哦！"
mention_users:
  - "@all"
holidays: []
```

配置文件已被 Git 忽略，请勿提交真实密钥。

## 构建与启动

```bash
go test ./...
go build -o weatherrobot ./cmd/weatherrobot
./weatherrobot
```

启动后不监听 HTTP 端口，只运行每天 08:00 的定时天气任务。

## 手动测试

```bash
go run ./cmd/test-weather
```

该命令会立即请求天气数据并发送到配置的企业微信群，请谨慎在生产群执行。

## 当前消息内容

- 城市名称
- 当前天气和温度
- 风向、风力
- 未来三日天气及温度范围
- 生活指数
- 法定节假日问候

## 常见问题

- 启动失败：检查配置文件是否存在，以及 Webhook、天气 Key、城市列表是否完整。
- 不发送消息：检查日志中的Open-Meteo HTTP 错误信息或企业微信 `errcode`。
- 城市显示未知：在 `internal/weather/weather.go` 的 `GetCityName` 中补充城市代码映射。
