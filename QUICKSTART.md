# 快速开始指南

## 5分钟快速部署

### 前置条件
- Linux 服务器（或 macOS）
- Go 1.18+ 环境
- 企业微信应用和 Webhook 密钥
- 豆包 AI 或 OpenAI API 密钥

### 第一步：获取代码

```bash
git clone https://github.com/MuBuyi/weatherrobot.git
cd weatherrobot
```

### 第二步：配置密钥

编辑 `config/config.yaml`，填入你的密钥：

```yaml
# ⚠️ 重要：这三个必须配置
wecom_webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY_HERE"
doubao_api_key: "YOUR_DOUBAO_API_KEY"
weather_api_key: "YOUR_WEATHER_API_KEY"

# 可选：OpenAI 备用方案
openai_api_key: "sk-your-openai-key-here"
```

### 第三步：编译

```bash
go build -o weatherrobot ./cmd/weatherrobot
```

### 第四步：启动

```bash
# 前台运行（调试）
./weatherrobot

# 后台运行（生产）
nohup ./weatherrobot > weatherrobot.log 2>&1 &
```

### 第五步：配置企业微信

在企业微信应用的**消息接收**中，配置回调地址：

```
http://your-server-ip:9001/wecom/message
```

## 测试功能

### 在企业微信群中测试

在群中输入：
```
@皮皮虾 今天天气怎么样？
```

机器人应该在几秒内回复。

### 查看日志

```bash
# 查看实时日志
tail -f weatherrobot.log

# 查看消息处理
grep "收到来自\|发送成功" weatherrobot.log
```

## 常用命令

```bash
# 编译
go build -o weatherrobot ./cmd/weatherrobot

# 运行
./weatherrobot

# 后台运行
nohup ./weatherrobot > weatherrobot.log 2>&1 &

# 停止
pkill -f weatherrobot

# 查看状态
ps aux | grep weatherrobot

# 查看日志
tail -100 weatherrobot.log

# 查看错误
grep -i error weatherrobot.log
```

## 支持的功能

✅ **定时天气报告** - 自动发送天气预报  
✅ **AI 下班提醒** - 豆包生成个性化文案  
✅ **群聊交互** - @机器人 实时回答问题  
✅ **自动容错** - 豆包失败自动用 OpenAI  
✅ **异步处理** - 不阻塞主程序流程  

## 消息格式

### 支持的格式

```
✓ @皮皮虾 今天天气怎么样？
✓ @机器人 请帮我回答
✓ @AI 下班后该做什么？
```

### 不支持的格式（会被忽略）

```
✗ 这是一条普通消息
✗ @ 只有@没有问题
✗ 没有@的问题
```

## 获取 API 密钥

### 豆包 API（推荐）
1. 访问 [火山引擎](https://console.volcengine.com/)
2. 申请豆包 API
3. 获取 API 密钥

### 天气 API
1. 访问 [QWeather](https://www.qweather.com/)
2. 注册并获取 API 密钥
3. 查询城市代码

### OpenAI（备用）
1. 访问 [OpenAI](https://platform.openai.com/)
2. 创建 API 密钥

## 故障排查

| 问题 | 解决方案 |
|------|--------|
| 程序无法启动 | 检查端口 9001 是否被占用：`lsof -i :9001` |
| 消息收不到 | 确认 webhook 回调地址配置正确 |
| AI 不回复 | 检查 API 密钥和网络连接 |
| 日志文件过大 | 定期备份和清理：`rm weatherrobot.log` |

## 性能建议

- **单服务器**：支持 100+ 用户并发
- **消息处理**：平均延迟 2-5 秒（包含 AI 调用）
- **内存占用**：约 30-50 MB
- **CPU 占用**：大部分时间 < 1%

## 升级更新

```bash
# 获取最新代码
git pull origin main

# 重新编译
go build -o weatherrobot ./cmd/weatherrobot

# 重启程序
pkill -f weatherrobot
nohup ./weatherrobot > weatherrobot.log 2>&1 &
```

## 需要帮助？

1. 查看 [DEPLOYMENT.md](DEPLOYMENT.md) 了解详细部署说明
2. 查看 [WECOM_INTERACTION.md](WECOM_INTERACTION.md) 了解交互功能
3. 查看日志文件诊断问题
4. 在 GitHub 提出 Issue

## 后续扩展

- [ ] 添加命令识别（如 `/天气`、`/帮助`）
- [ ] 支持消息内容审核
- [ ] 实现对话上下文记忆
- [ ] 支持更多城市和天气源
- [ ] 添加 Web 管理界面

---

祝您使用愉快！🎉
