# 天气机器人 - 生产部署指南

## 功能概述

这是一个企业微信集成的天气机器人，具有以下功能：

1. **定时天气报告** - 每天定时发送天气预报和生活指数
2. **下班提醒** - AI 生成温暖的下班提醒文案
3. **群聊交互** - 在群中 @机器人 可以实时获得 AI 回答

## 系统架构

```
┌─────────────────────┐
│   定时任务调度      │
├─────────────────────┤
│  • 天气报告 (1分钟) │
│  • 下班提醒 (1分钟) │
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│  天气机器人核心    │
├─────────────────────┤
│  • 配置管理         │
│  • AI 集成          │  ◄─── 豆包 AI (主) / OpenAI (备)
│  • 消息处理         │
│  • 日志系统         │
└──────────┬──────────┘
           │
┌──────────▼──────────┐
│   企业微信 Webhook  │
├─────────────────────┤
│  • 消息接收 (9001)  │
│  • 消息发送         │
└─────────────────────┘
```

## 部署步骤

### 1. 编译程序

```bash
cd /path/to/weatherrobot
go build -o weatherrobot ./cmd/weatherrobot
```

### 2. 配置文件

编辑 `config/config.yaml`：

```yaml
# 企业微信 Webhook 地址
wecom_webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY"

# 天气 API
weather_api_key: "YOUR_WEATHER_API_KEY"

# 豆包 AI 配置（推荐）
doubao_url: "https://ark.cn-beijing.volces.com/api/v3/responses"
doubao_api_key: "YOUR_DOUBAO_API_KEY"
doubao_model: "doubao-seed-1-8-251228"

# OpenAI 配置（备用）
openai_api_key: "sk-YOUR_OPENAI_KEY"

# 城市代码（QWeather）
locations:
  - "101200805"  # 监利

# 提及人员
mention_users:
  - "@all"

# 下班文案
use_ai_reminder: true
off_work_messages: [...]
```

### 3. 启动程序

```bash
# 前台运行（调试）
./weatherrobot

# 后台运行（生产）
nohup ./weatherrobot > weatherrobot.log 2>&1 &

# 使用 systemd 开机自启
# 见下面的 Systemd 配置部分
```

### 4. 企业微信配置

在企业微信应用中配置 Webhook 回调：

**接收消息配置：**
- 回调地址: `http://your-server-ip:9001/wecom/message`
- 需要在消息接收中启用

**发送消息配置：**
- 使用机器人 Webhook 地址（在 config.yaml 中配置）

## 消息格式

### 群聊交互

在企业微信群中 @机器人，格式如下：

```
@皮皮虾 今天天气怎么样？
@机器人 下班后该做什么？
@皮皮虾 如何更好地学习？
```

机器人会自动识别 `@机器人` 格式的消息，提取问题部分，并通过豆包 AI 生成回答。

**不支持的格式（会被忽略）：**
```
这是一条普通消息
```

### 定时任务

```
每 1 分钟执行:
✓ 发送天气报告（可改为每天 8 点）
✓ 发送下班提醒（可改为每天 18 点）
```

## Systemd 配置（可选）

创建 `/etc/systemd/system/weatherrobot.service`：

```ini
[Unit]
Description=Weather Robot Service
After=network.target

[Service]
Type=simple
User=your_user
WorkingDirectory=/path/to/weatherrobot
ExecStart=/path/to/weatherrobot/weatherrobot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启用和启动：

```bash
sudo systemctl daemon-reload
sudo systemctl enable weatherrobot
sudo systemctl start weatherrobot
sudo systemctl status weatherrobot
```

## 日志管理

日志文件位置：
- 标准输出：直接打印到控制台
- 文件保存：可用 `nohup` 重定向到文件

查看日志：

```bash
# 实时日志
tail -f weatherrobot.log

# 关键日志
grep "ERROR\|WARN" weatherrobot.log

# 消息处理日志
grep "收到来自\|发送成功" weatherrobot.log
```

## 故障排查

### 问题 1: 消息收不到

**检查清单：**
- [ ] 确认企业微信应用已配置 webhook 回调地址
- [ ] 确认服务器防火墙允许 9001 端口入站
- [ ] 检查程序是否成功启动：`ps aux | grep weatherrobot`
- [ ] 查看日志：`grep "企业微信消息服务启动" weatherrobot.log`

### 问题 2: @机器人 没有回复

**检查清单：**
- [ ] 确认消息格式是 `@机器人名称 问题`
- [ ] 检查豆包 API 密钥是否正确
- [ ] 查看日志中的 AI 调用错误：`grep "调用豆包\|调用 OpenAI" weatherrobot.log`

### 问题 3: 定时任务不执行

**检查清单：**
- [ ] 程序日志中是否有任务启动信息
- [ ] 检查配置文件中的 Webhook 地址是否正确
- [ ] 查看是否有网络连接问题

### 问题 4: 内存占用过高

**解决方案：**
- 检查是否有大量日志输出，可以降低日志级别
- 检查是否有消息处理阻塞，查看日志延迟

## 性能优化

### 调整定时任务频率

编辑 `cmd/weatherrobot/main.go`：

```go
// 改为每天 8 点执行天气报告
c.AddFunc("0 8 * * *", cronn.SendDailyReport)

// 改为每天 18 点执行下班提醒
c.AddFunc("0 18 * * *", cronn.SendOffWorkReminder)
```

### 并发控制

当前采用异步处理，单个消息不会阻塞主程序。

### 缓存优化

天气数据由 QWeather API 提供缓存，无需额外配置。

## 监控和告警

建议监控以下指标：

1. **进程状态**：程序是否还在运行
2. **端口监听**：9001 端口是否开放
3. **磁盘空间**：日志文件大小
4. **API 调用**：豆包/OpenAI 的成功率

## 备份和恢复

### 备份配置

```bash
cp config/config.yaml config/config.yaml.bak
```

### 更新程序

```bash
git pull origin main
go build -o weatherrobot ./cmd/weatherrobot
systemctl restart weatherrobot
```

## 常见问题

**Q: 支持多个城市吗？**
A: 支持。在 `config.yaml` 的 `locations` 中添加多个城市代码。

**Q: 可以自定义下班文案吗？**
A: 可以。编辑 `config.yaml` 的 `off_work_messages` 或设置 `use_ai_reminder: true` 使用 AI 生成。

**Q: 支持中英文混合吗？**
A: 支持。豆包和 OpenAI 都支持多语言。

**Q: 可以取消某个定时任务吗？**
A: 可以。在 `main.go` 中注释掉对应的 `c.AddFunc()` 调用。

## 获取帮助

- 查看日志文件了解错误详情
- 检查 GitHub 仓库的 Issue
- 参考 [WECOM_INTERACTION.md](WECOM_INTERACTION.md) 了解交互功能

## 版本信息

- **Go 版本**: 1.18+
- **依赖**: robfig/cron/v3, sirupsen/logrus, spf13/viper
- **最后更新**: 2026-02-05
