package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"wechatrobot/internal/log"

	"github.com/sirupsen/logrus"
)

// 豆包 API 请求结构
type DoubaoRequest struct {
	Model string               `json:"model"`
	Input []DoubaoInputMessage `json:"input"`
}

type DoubaoInputMessage struct {
	Role    string              `json:"role"`
	Content []DoubaoContentItem `json:"content"`
}

type DoubaoContentItem struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// 豆包 API 响应结构
type DoubaoResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Role    string `json:"role"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

// OpenAI API 请求结构
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI API 响应结构
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateOffWorkReminder 使用 AI 生成下班提醒文案
// 会优先调用豆包（Doubao），失败则回退到 OpenAI
// 传入参数：doubaoURL, doubaoKey, doubaoModel, openaiKey
func GenerateOffWorkReminder(doubaoURL, doubaoKey, doubaoModel, openaiKey string) (string, error) {
	// 如果参数未传入，则尝试从环境变量读取
	if doubaoURL == "" {
		doubaoURL = os.Getenv("DOUBAO_URL")
	}
	if doubaoKey == "" {
		doubaoKey = os.Getenv("DOUBAO_API_KEY")
	}
	if doubaoModel == "" {
		doubaoModel = os.Getenv("DOUBAO_MODEL")
		if doubaoModel == "" {
			doubaoModel = "doubao-seed-1-8-251228" // 默认模型
		}
	}
	if openaiKey == "" {
		openaiKey = os.Getenv("OPENAI_API_KEY")
	}

	prompt := `生成一条有趣、温暖、鼓励的下班提醒文案。要求：
1. 字数在30-80字之间
2. 包含对员工的关怀和鼓励
3. 古诗词或者当前网络热梗或者热词（任选其一）来增加趣味性
4. 轻松幽默，避免枯燥
5. 不要使用编号或列表格式
6. 直接返回文案内容，不要包含任何前缀或说明
7. 如果是节假日，文案中要体现节日氛围

请生成：`

	// 如果配置了豆包，优先调用豆包
	if doubaoURL != "" && doubaoKey != "" {
		result, err := callDoubao(doubaoURL, doubaoKey, doubaoModel, prompt)
		if err == nil && result != "" {
			log.Info("Doubao 生成的下班提醒: ", result)
			return result, nil
		}
		logrus.Warnf("Doubao 调用失败: %v，回退到 OpenAI", err)
	}

	// 回退到 OpenAI
	if openaiKey == "" {
		return "", fmt.Errorf("没有可用的 AI 提供者（Doubao/OpenAI）")
	}

	result, err := callOpenAI(openaiKey, prompt)
	if err != nil {
		logrus.Errorf("OpenAI 调用失败: %v", err)
		return "", err
	}
	log.Info("OpenAI 生成的下班提醒: ", result)
	return result, nil
}

// callDoubao 调用豆包 API
func callDoubao(url, apiKey, model, prompt string) (string, error) {
	requestBody := DoubaoRequest{
		Model: model,
		Input: []DoubaoInputMessage{
			{
				Role: "user",
				Content: []DoubaoContentItem{
					{
						Type: "input_text",
						Text: prompt,
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用 Doubao API 失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Doubao API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	var doubaoResp DoubaoResponse
	if err := json.Unmarshal(bodyBytes, &doubaoResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取 output 中的内容
	if len(doubaoResp.Output) > 0 {
		for _, item := range doubaoResp.Output {
			if item.Type == "message" && len(item.Content) > 0 {
				for _, content := range item.Content {
					if content.Type == "output_text" && content.Text != "" {
						return content.Text, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("豆包 API 返回空响应或无法解析内容")
}

// callOpenAI 调用 OpenAI API
func callOpenAI(apiKey, prompt string) (string, error) {
	requestBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 150,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用 OpenAI API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI 返回空响应")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

// AskDoubao 使用豆包 AI 回答用户问题
func AskDoubao(question, url, apiKey, model string) (string, error) {
	if url == "" || apiKey == "" {
		return "", fmt.Errorf("豆包配置不完整")
	}

	requestBody := DoubaoRequest{
		Model: model,
		Input: []DoubaoInputMessage{
			{
				Role: "user",
				Content: []DoubaoContentItem{
					{
						Type: "input_text",
						Text: question,
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用 Doubao API 失败: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Doubao API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	var doubaoResp DoubaoResponse
	if err := json.Unmarshal(bodyBytes, &doubaoResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取 output 中的内容
	if len(doubaoResp.Output) > 0 {
		for _, item := range doubaoResp.Output {
			if item.Type == "message" && len(item.Content) > 0 {
				for _, content := range item.Content {
					if content.Type == "output_text" && content.Text != "" {
						return content.Text, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("豆包 API 返回空响应或无法解析内容")
}

// AskOpenAI 使用 OpenAI 回答用户问题
func AskOpenAI(question, apiKey string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API Key 未配置")
	}

	requestBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens: 500,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用 OpenAI API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("OpenAI 返回空响应")
	}

	return openaiResp.Choices[0].Message.Content, nil
}
