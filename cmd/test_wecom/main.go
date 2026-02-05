package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 模拟企业微信的消息结构
type TestMessage struct {
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

func main() {
	// 测试消息列表
	testMessages := []string{
		"今天天气怎么样？",
		"请告诉我明天的天气",
		"下班后做什么比较好？",
		"如何更好地工作？",
		"豆包 AI 的优势是什么？",
	}

	fmt.Println("=== 微信交互功能测试 ===\n")
	fmt.Println("机器人应该监听 http://localhost:9001/wecom/message")
	fmt.Println("确保程序已启动...\n")

	for i, question := range testMessages {
		fmt.Printf("\n[测试 %d] 发送问题: %s\n", i+1, question)

		msg := TestMessage{
			ToUserID:   "robot",
			FromUserID: "test_user",
			CreateTime: time.Now().Unix(),
			MsgType:    "text",
			Content:    question,
			MsgID:      fmt.Sprintf("msg_%d", time.Now().UnixNano()),
			AgentID:    1000001,
		}
		msg.Text.Content = question

		jsonBody, _ := json.Marshal(msg)

		resp, err := http.Post(
			"http://localhost:9001/wecom/message",
			"application/json",
			bytes.NewBuffer(jsonBody),
		)

		if err != nil {
			fmt.Printf("❌ 发送失败: %v\n", err)
			continue
		}

		fmt.Printf("✓ 发送成功，状态码: %d\n", resp.StatusCode)
		resp.Body.Close()

		// 等待 AI 处理（通常需要 2-10 秒）
		fmt.Println("⏳ 等待 AI 处理中（3 秒）...")
		time.Sleep(3 * time.Second)
	}

	fmt.Println("\n✓ 所有测试完成！检查机器人日志查看处理结果。")
}
