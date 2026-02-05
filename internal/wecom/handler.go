package wecom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"wechatrobot/internal/ai"
	"wechatrobot/internal/config"
	"wechatrobot/internal/weather"

	"github.com/sirupsen/logrus"
)

// 企业微信消息接收结构
type WecomIncomingMessage struct {
	ToUserID   string `json:"ToUserID"`
	FromUserID string `json:"FromUserID"`
	CreateTime int64  `json:"CreateTime"`
	MsgType    string `json:"MsgType"`
	Content    string `json:"Content"`
	MsgID      string `json:"MsgId"`
	AgentID    int    `json:"AgentID"`
	Text       struct {
		Content string `json:"content"`
	} `json:"Text"`
}

// StartWecomServer 启动企业微信消息接收服务
func StartWecomServer(port string) {
	http.HandleFunc("/wecom/message", HandleWecomMessage)
	logrus.Infof("企业微信消息服务启动，监听端口 %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Errorf("企业微信消息服务启动失败: %v", err)
	}
}

// HandleWecomMessage 处理企业微信的消息 webhook
func HandleWecomMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logrus.Debugf("收到企业微信消息: %s", string(body))

	var msg WecomIncomingMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		logrus.Errorf("解析消息失败: %v", err)
		http.Error(w, "Failed to parse message", http.StatusBadRequest)
		return
	}

	// 只处理文本消息
	if msg.MsgType != "text" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	// 获取消息内容
	var messageContent string
	if msg.Text.Content != "" {
		messageContent = msg.Text.Content
	} else if msg.Content != "" {
		messageContent = msg.Content
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	logrus.Infof("收到来自 %s 的消息: %s", msg.FromUserID, messageContent)

	// 检查是否是 @机器人 的消息
	question := extractQuestion(messageContent)
	if question == "" {
		// 不是 @机器人 的消息，忽略
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	// 异步处理消息，快速响应
	go ProcessUserMessage(question, msg.FromUserID)

	// 立即响应 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// ProcessUserMessage 处理用户消息并通过 AI 生成回复
func ProcessUserMessage(userMessage, userID string) {
	// 使用豆包 AI 来回答问题
	reply, err := ai.AskDoubao(userMessage, config.Cfg.DoubaoURL, config.Cfg.DoubaoAPIKey, config.Cfg.DoubaoModel)
	if err != nil {
		logrus.Warnf("调用豆包 AI 失败: %v，尝试回退到 OpenAI", err)
		reply, err = ai.AskOpenAI(userMessage, config.Cfg.OpenAIAPIKey)
		if err != nil {
			logrus.Errorf("调用 OpenAI 也失败: %v", err)
			reply = "抱歉，我现在无法处理您的问题，请稍后再试。"
		}
	}

	// 构造回复消息
	responseContent := fmt.Sprintf("@%s\n\n%s", userID, reply)

	// 发送回复
	if err := weather.SendWecomMessage(responseContent, nil); err != nil {
		logrus.Errorf("发送回复消息失败: %v", err)
	} else {
		logrus.Infof("已向 %s 发送回复", userID)
	}
}

// extractQuestion 从企业微信消息中提取问题
// 支持的格式: "@机器人名称 问题内容" 或 "@机器人名称  问题内容"
// 如果不是 @机器人 的消息，返回空字符串
func extractQuestion(message string) string {
	// 移除消息前的空格
	trimmed := strings.TrimSpace(message)
	
	// 检查是否以 @ 开头
	if !strings.HasPrefix(trimmed, "@") {
		return ""
	}
	
	// 找到空格位置，分离 @名称 和 问题
	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) < 2 {
		// 没有问题内容，只有 @
		return ""
	}
	
	// 第二部分是问题，移除前后空格
	question := strings.TrimSpace(parts[1])
	
	// 如果问题为空，忽略
	if question == "" {
		return ""
	}
	
	return question
}
