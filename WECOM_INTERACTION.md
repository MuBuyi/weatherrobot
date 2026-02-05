# 微信群聊交互功能说明

## 功能概述

已添加微信群聊交互功能，用户可以在企业微信群中与机器人互动。机器人通过豆包 AI（ByteDance Doubao）来理解和回答用户的问题。

## 架构设计

```
用户发送消息 → 企业微信 → webhook 回调 → 机器人 HTTP 服务
                                           ↓
                                    消息处理线程
                                           ↓
                                    豆包 AI 或 OpenAI
                                           ↓
                                    生成回复文案
                                           ↓
                                    发送回复到企业微信
```

## 工作流程

1. **消息接收**: 机器人在端口 9001 监听 `/wecom/message` 路由
2. **快速响应**: 立即返回 200 OK 给企业微信，避免超时
3. **异步处理**: 在后台线程中处理用户消息
4. **AI 生成**: 调用豆包 API 生成回复（主）或 OpenAI（备用）
5. **发送回复**: 通过企业微信 Webhook 发送回复消息

## 配置步骤

### 1. 企业微信应用配置

需要在企业微信应用中配置 webhook 回调地址：

```
回调地址: http://<your-server-ip>:9001/wecom/message
```

其中 `<your-server-ip>` 是运行机器人的服务器 IP 地址。

### 2. 配置文件（config/config.yaml）

已包含所有必需的配置：

```yaml
wecom_webhook: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=..."
doubao_url: "https://ark.cn-beijing.volces.com/api/v3/responses"
doubao_api_key: "your-doubao-api-key"
doubao_model: "doubao-seed-1-8-251228"
openai_api_key: "sk-your-openai-api-key"  # 备用
```

## 代码结构

### 新增文件

- `internal/wecom/handler.go` - 企业微信消息处理逻辑
  - `StartWecomServer()` - 启动 HTTP 服务
  - `HandleWecomMessage()` - Webhook 处理函数
  - `ProcessUserMessage()` - 异步处理消息

### 修改文件

- `internal/ai/ai.go`
  - `AskDoubao()` - 调用豆包 AI 回答问题
  - `AskOpenAI()` - 调用 OpenAI 回答问题（备用）

- `cmd/weatherrobot/main.go`
  - 添加后台启动微信消息服务

## 使用示例

### 用户在群聊中提问

```
用户: "今天天气怎么样？"
机器人: "@user_id

今天是一个多云的日子，温度约在2-8°C之间...
```

### 自动回复流程

```
1. 企业微信收到群消息
2. 触发 webhook 回调到机器人
3. 机器人快速返回 200 OK
4. 后台调用豆包 AI 处理
5. AI 生成回复
6. 机器人通过 webhook 发送回复
```

## 功能特点

✅ **快速响应** - 异步处理，不阻塞 webhook 回调  
✅ **AI 驱动** - 支持豆包和 OpenAI 两种 AI 引擎  
✅ **容错机制** - 豆包失败自动回退到 OpenAI  
✅ **灵活部署** - 可配置监听端口和 webhook 地址  
✅ **日志完整** - 详细的执行日志用于调试  

## 故障排查

### 问题：机器人收不到消息

1. 检查企业微信应用是否配置了 webhook 回调地址
2. 确认服务器防火墙允许端口 9001 的入站流量
3. 查看日志 `INFO[xxxx] 企业微信消息服务启动，监听端口 9001`

### 问题：回复消息延迟

- 这是正常的，因为需要调用 AI API，通常需要 2-10 秒

### 问题：收不到 AI 的回复

1. 检查豆包 API key 是否有效
2. 查看日志中的错误信息
3. 确认 OpenAI API key 已配置（作为备用）

## 性能考虑

- 当前采用异步处理，单个消息处理不影响其他消息
- 建议单个群组每秒不超过 10 条消息（避免 API 配额限制）
- 可根据需要调整日志级别来优化性能

## 安全建议

⚠️ **重要**：不要在代码中硬编码 API 密钥，使用配置文件或环境变量

✅ 已实现：使用 config.yaml 管理所有敏感信息

## 后续扩展方向

- [ ] 添加消息内容过滤/审核
- [ ] 支持多种消息类型（图片、文件等）
- [ ] 添加用户使用统计
- [ ] 支持自定义回复模板
- [ ] 实现消息持久化存储
