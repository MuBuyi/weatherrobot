package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/log"
)

func wecomHTTPClient() *http.Client {
	return &http.Client{Timeout: time.Duration(config.Cfg.HTTPTimeoutSeconds) * time.Second}
}

type wecomResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WecomMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string   `json:"content"`
		Mention []string `json:"mentioned_list"`
	} `json:"text"`
}

func SendWecomMessage(content string, mentionUsers []string) error {
	message := WecomMessage{
		MsgType: "text",
	}
	message.Text.Content = content
	message.Text.Mention = mentionUsers

	messageBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := wecomHTTPClient().Post(config.Cfg.WecomWebhook, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取企业微信响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送消息失败，HTTP 状态码: %d", resp.StatusCode)
	}

	var result wecomResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析企业微信响应失败: %w", err)
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("企业微信发送失败: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	log.Info("消息发送成功")
	return nil
}
