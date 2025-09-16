package services

import (
	"aigames/internal/ai"
	"aigames/internal/config"
	"aigames/internal/models"
	"aigames/pkg/logger"
	"strings"
	"time"
)

// AIService AI服务接口
type AIService struct {
	provider      string
	openAIClient  *ai.APIClient
	geminiClient  *ai.GeminiClient
	promptBuilder *ai.PromptBuilder
}

// NewAIService 创建AI服务
func NewAIService() *AIService {
	aiConfig := config.GetConfig().AI

	var openAIClient *ai.APIClient
	var geminiClient *ai.GeminiClient

	// 根据提供商创建相应的客户端
	if aiConfig.Provider == "gemini" {
		geminiClient = ai.NewGeminiClient(&aiConfig)
	} else {
		openAIClient = ai.NewAPIClient(&aiConfig)
	}

	promptBuilder := ai.NewPromptBuilder()
	// 设置提示词构建器的提供商
	promptBuilder.SetProvider(aiConfig.Provider)

	return &AIService{
		provider:      aiConfig.Provider,
		openAIClient:  openAIClient,
		geminiClient:  geminiClient,
		promptBuilder: promptBuilder,
	}
}

// CallLandlord AI叫地主决策
func (s *AIService) CallLandlord(player *models.GamePlayer, game *models.Game, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string) error {
	// 构建提示词
	prompt := s.promptBuilder.BuildCallLandlordPrompt(player, game)
	logger.Info("----------------------------------")
	logger.Info("AI CallLandlord Prompt: %s", prompt)
	logger.Info("----------------------------------")

	// 构建消息
	messages := []ai.Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	response, err := s.chatWithRetry(messages)
	if err != nil {
		// 如果API调用失败，使用默认策略（不叫地主）
		return ai.CallLandlord(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID, false)
	}
	return ai.APICallLandlordWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
}

// chatWithRetry 带有重试逻辑的聊天请求
func (s *AIService) chatWithRetry(messages []ai.Message) (string, error) {
	var response string
	var err error
	const maxRetries = 3
	const retryDelay = 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		if s.provider == "gemini" {
			response, err = s.geminiClient.Chat(messages)
		} else {
			response, err = s.openAIClient.Chat(messages)
		}

		if err == nil {
			return response, nil // Success
		}

		// 检查是否为503错误
		if strings.Contains(err.Error(), "状态码: 503") {
			logger.Info("API call failed with 503, retrying in %v... (%d/%d)", retryDelay, i+1, maxRetries)
			time.Sleep(retryDelay)
			continue
		}

		// 对于其他错误，立即失败
		return "", err
	}

	return "", err // 所有重试失败后返回最后一个错误
}

// PlayCards AI出牌决策
func (s *AIService) PlayCards(player *models.GamePlayer, game *models.Game, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
	PassTurn(roomID, username string) error
}, roomID string) error {
	// 构建提示词
	prompt := s.promptBuilder.BuildPlayCardPrompt(player, game)
	logger.Info("----------------------------------")
	logger.Info("AI PlayCards Prompt: %s", prompt)
	logger.Info("----------------------------------")

	// 构建消息
	messages := []ai.Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	response, err := s.chatWithRetry(messages)
	if err != nil {
		// API调用失败，检查是否可以过牌
		if game.LastPlayer != player.Position {
			return ai.PassTurn(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
		} else {
			// 不能过牌，必须出牌。作为备用策略，出手中最小的一张牌。
			if len(player.Cards) > 0 {
				// 玩家手牌已经排序，第一张就是最小的
				smallestCard := []models.Card{player.Cards[0]}
				return ai.PlayCards(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID, smallestCard)
			}
		}
		// 如果没有手牌了，这本身是一个错误状态，但还是尝试过牌
		return ai.PassTurn(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
	}

	return ai.APIPlayCardsWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, player.Cards, gameService, roomID)
}
