package weather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"wechatrobot/internal/config"
	"wechatrobot/internal/log"
)

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

	resp, err := http.Post(config.Cfg.WecomWebhook, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送消息失败，状态码: %d", resp.StatusCode)
	}

	log.Info("消息发送成功")
	return nil
}
