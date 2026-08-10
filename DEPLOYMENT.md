# 天气预报提醒机器人：部署说明

## 架构

```text
Cron（北京时间每天 08:00）
  -> Open-Meteo 当前天气
  -> Open-Meteo 每日预报（展示前三日）
  -> 体感温度、湿度和降水概率
  -> 组合天气文本
  -> 企业微信群机器人 Webhook
```

项目不包含 AI、群聊问答、数据库和 Web 管理页面，也不需要开放入站端口。

## 生产构建

```bash
go test ./...
go build -trimpath -ldflags="-s -w" -o weatherrobot ./cmd/weatherrobot
```

## systemd 示例

```ini
[Unit]
Description=Weather Reminder Robot
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=weatherrobot
WorkingDirectory=/opt/weatherrobot
ExecStart=/opt/weatherrobot/weatherrobot
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now weatherrobot
sudo systemctl status weatherrobot
```

## 运行要求

- 工作目录下必须存在 `config/config.yaml`
- `weather_cron` 控制发送时间，使用标准 5 段 Cron 表达式
- `timezone` 控制调度时区
- `city_names` 维护城市代码与显示名称
- `forecast_days` 控制展示预报天数
- `reminder_title` 和 `reminder_intro` 配置群提醒文案
- - `footer` 可自定义消息结尾
- 服务器需要访问 `api.open-meteo.com` 和 `qyapi.weixin.qq.com`
- 系统时间应正确；调度器固定使用 UTC+8
- 无需开放 9001 或其他入站端口

## 可靠性行为

- 外部 HTTP 请求超时为 10 秒
- 校验 Open-Meteo HTTP 状态和错误信息
- 校验企业微信 HTTP 状态和业务 `errcode`
- 不在日志中打印 API Key 或完整天气响应
- 所有城市数据均获取失败时不会发送空提醒

## 更新

```bash
git pull
go test ./...
go build -trimpath -ldflags="-s -w" -o weatherrobot ./cmd/weatherrobot
sudo systemctl restart weatherrobot
```
