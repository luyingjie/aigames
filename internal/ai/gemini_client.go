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

// GeminiClient Gemini API客户端
type GeminiClient struct {
	config *config.AIConfig
	client *http.Client
}

// NewGeminiClient 创建Gemini API客户端
func NewGeminiClient(aiConfig *config.AIConfig) *GeminiClient {
	return &GeminiClient{
		config: aiConfig,
		client: &http.Client{
			Timeout: time.Duration(aiConfig.Timeout) * time.Second,
		},
	}
}

// GeminiChatRequest Gemini聊天请求
type GeminiChatRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
}

// Content 内容
type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part 部分内容
type Part struct {
	Text string `json:"text"`
}

// GenerationConfig 生成配置
type GenerationConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GeminiChatResponse Gemini聊天响应
type GeminiChatResponse struct {
	Candidates    []Candidate   `json:"candidates"`
	UsageMetadata UsageMetadata `json:"usageMetadata"`
}

// Candidate 候选回复
type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
}

// UsageMetadata 使用元数据
type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// Chat 发送聊天请求到Gemini API
func (c *GeminiClient) Chat(messages []Message) (string, error) {
	// 将标准消息格式转换为Gemini格式
	contents := make([]Content, len(messages))
	for i, msg := range messages {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		contents[i] = Content{
			Role: role,
			Parts: []Part{
				{Text: msg.Content},
			},
		}
	}

	request := GeminiChatRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature:     c.config.DefaultTemperature,
			MaxOutputTokens: c.config.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建Gemini API URL
	apiURL := fmt.Sprintf("%s/%s:generateContent?key=%s", c.config.APIURL, c.config.DefaultModel, c.config.APIKey)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	var geminiResponse GeminiChatResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(geminiResponse.Candidates) == 0 {
		logger.Warn("Gemini API返回空的candidates数组, 响应体: %s", string(body))
		return "", fmt.Errorf("API返回空结果")
	}

	candidate := geminiResponse.Candidates[0]
	logger.Info("Gemini API响应, finishReason: %s", candidate.FinishReason)

	// 提取回复文本
	var responseText string
	if candidate.FinishReason == "STOP" && len(candidate.Content.Parts) > 0 {
		responseText = candidate.Content.Parts[0].Text
	} else {
		logger.Warn("Gemini API返回了非STOP的finishReason或空的parts, 响应体: %s", string(body))
	}

	logger.Info("Gemini API调用成功，使用token: %d", geminiResponse.UsageMetadata.TotalTokenCount)
	return responseText, nil
}
