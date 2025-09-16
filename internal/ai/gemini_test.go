package ai

import (
	"aigames/internal/config"
	"testing"
)

// TestGeminiClientInitialization 测试Gemini客户端初始化
func TestGeminiClientInitialization(t *testing.T) {
	// 创建一个模拟的AI配置
	aiConfig := config.AIConfig{
		DefaultModel:       "gemini-2.5-pro",
		Timeout:            30,
		DefaultTemperature: 0.7,
		MaxTokens:          1000,
		APIURL:             "https://generativelanguage.googleapis.com/v1beta/models",
		APIKey:             "test-api-key",
		Provider:           "gemini",
	}

	// 创建Gemini客户端
	client := NewGeminiClient(&aiConfig)

	if client == nil {
		t.Error("Failed to create Gemini client")
	}

	if client.config.DefaultModel != "gemini-2.5-pro" {
		t.Errorf("Expected default model to be 'gemini-2.5-pro', got '%s'", client.config.DefaultModel)
	}
}

// TestGeminiChatRequest 测试Gemini聊天请求结构
func TestGeminiChatRequest(t *testing.T) {
	request := GeminiChatRequest{
		Contents: []Content{
			{
				Role: "user",
				Parts: []Part{
					{Text: "Hello, Gemini!"},
				},
			},
		},
		GenerationConfig: GenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 1000,
		},
	}

	if len(request.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(request.Contents))
	}

	if request.Contents[0].Role != "user" {
		t.Errorf("Expected role to be 'user', got '%s'", request.Contents[0].Role)
	}

	if request.Contents[0].Parts[0].Text != "Hello, Gemini!" {
		t.Errorf("Expected text to be 'Hello, Gemini!', got '%s'", request.Contents[0].Parts[0].Text)
	}
}

// TestGeminiChatResponse 测试Gemini聊天响应结构
func TestGeminiChatResponse(t *testing.T) {
	response := GeminiChatResponse{
		Candidates: []Candidate{
			{
				Content: Content{
					Role: "model",
					Parts: []Part{
						{Text: "Hello! How can I help you today?"},
					},
				},
				FinishReason: "STOP",
			},
		},
		UsageMetadata: UsageMetadata{
			PromptTokenCount:     10,
			CandidatesTokenCount: 20,
			TotalTokenCount:      30,
		},
	}

	if len(response.Candidates) != 1 {
		t.Errorf("Expected 1 candidate, got %d", len(response.Candidates))
	}

	if response.Candidates[0].Content.Parts[0].Text != "Hello! How can I help you today?" {
		t.Errorf("Expected response text to be 'Hello! How can I help you today?', got '%s'", response.Candidates[0].Content.Parts[0].Text)
	}

	if response.UsageMetadata.TotalTokenCount != 30 {
		t.Errorf("Expected total token count to be 30, got %d", response.UsageMetadata.TotalTokenCount)
	}
}
