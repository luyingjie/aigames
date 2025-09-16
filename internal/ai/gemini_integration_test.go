package ai

import (
	"aigames/internal/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGeminiIntegration 测试Gemini API集成
func TestGeminiIntegration(t *testing.T) {
	// 创建一个模拟的Gemini API服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求路径
		expectedPath := "/v1beta/models/gemini-2.5-pro:generateContent"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 检查请求方法
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 检查Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 返回模拟响应
		response := GeminiChatResponse{
			Candidates: []Candidate{
				{
					Content: Content{
						Role: "model",
						Parts: []Part{
							{Text: "叫地主"},
						},
					},
					FinishReason: "STOP",
				},
			},
			UsageMetadata: UsageMetadata{
				PromptTokenCount:     50,
				CandidatesTokenCount: 10,
				TotalTokenCount:      60,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建AI配置
	aiConfig := config.AIConfig{
		DefaultModel:       "gemini-2.5-pro",
		Timeout:            30,
		DefaultTemperature: 0.7,
		MaxTokens:          1000,
		APIURL:             server.URL + "/v1beta/models",
		APIKey:             "test-api-key",
		Provider:           "gemini",
	}

	// 创建Gemini客户端
	client := NewGeminiClient(&aiConfig)

	// 发送测试请求
	messages := []Message{
		{Role: "user", Content: "你是否要叫地主？"},
	}

	response, err := client.Chat(messages)
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}

	// 验证响应
	expectedResponse := "叫地主"
	if response != expectedResponse {
		t.Errorf("Expected response '%s', got '%s'", expectedResponse, response)
	}
}

// TestGeminiClientWithDifferentModels 测试不同Gemini模型
func TestGeminiClientWithDifferentModels(t *testing.T) {
	models := []string{"gemini-2.5-pro", "gemini-2.5-flash"}

	for _, model := range models {
		aiConfig := config.AIConfig{
			DefaultModel:       model,
			Timeout:            30,
			DefaultTemperature: 0.7,
			MaxTokens:          1000,
			APIURL:             "https://generativelanguage.googleapis.com/v1beta/models",
			APIKey:             "test-api-key",
			Provider:           "gemini",
		}

		client := NewGeminiClient(&aiConfig)

		if client.config.DefaultModel != model {
			t.Errorf("Expected model %s, got %s", model, client.config.DefaultModel)
		}
	}
}
