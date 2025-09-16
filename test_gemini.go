package main

import (
	"aigames/internal/ai"
	"aigames/internal/config"
	"fmt"
	"log"
)

func main() {
	// 初始化配置
	err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 获取AI配置
	aiConfig := config.GetConfig().AI

	// 检查是否配置了Gemini
	if aiConfig.Provider != "gemini" {
		fmt.Println("Gemini provider not configured. Please check your config.yaml")
		return
	}

	fmt.Printf("Testing Gemini integration with model: %s\n", aiConfig.DefaultModel)

	// 创建Gemini客户端
	geminiClient := ai.NewGeminiClient(&aiConfig)

	// 构建测试消息
	messages := []ai.Message{
		{Role: "user", Content: "你好，Gemini！请简单介绍一下你自己。"},
	}

	// 发送请求
	response, err := geminiClient.Chat(messages)
	if err != nil {
		fmt.Printf("Error calling Gemini API: %v\n", err)
		return
	}

	fmt.Printf("Gemini response: %s\n", response)
	fmt.Println("Gemini integration test completed successfully!")
}
