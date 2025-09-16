package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"aigames/internal/config"
	"aigames/pkg/logger"
)

// APIClient AI API客户端
type APIClient struct {
	config *config.AIConfig
	client *http.Client
}

// NewAPIClient 创建AI API客户端
func NewAPIClient(aiConfig *config.AIConfig) *APIClient {
	return &APIClient{
		config: aiConfig,
		client: &http.Client{
			Timeout: time.Duration(aiConfig.Timeout) * time.Second,
		},
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat 发送聊天请求
func (c *APIClient) Chat(messages []Message) (string, error) {
	request := ChatRequest{
		Model:       c.config.DefaultModel,
		Messages:    messages,
		Temperature: c.config.DefaultTemperature,
		MaxTokens:   c.config.MaxTokens,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("API返回空结果")
	}

	logger.Info("AI API调用成功，使用token: %d", chatResponse.Usage.TotalTokens)
	return chatResponse.Choices[0].Message.Content, nil
}
