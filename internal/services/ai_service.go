package services

import (
	"aigames/internal/ai"
	"aigames/internal/config"
	"aigames/internal/models"
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

	// 构建消息
	messages := []ai.Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	// 根据提供商调用相应的API
	if s.provider == "gemini" {
		response, err := s.geminiClient.Chat(messages)
		if err != nil {
			// 如果API调用失败，使用默认策略（不叫地主）
			return ai.CallLandlord(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID, false)
		}
		return ai.APICallLandlordWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
	} else {
		response, err := s.openAIClient.Chat(messages)
		if err != nil {
			// 如果API调用失败，使用默认策略（不叫地主）
			return ai.CallLandlord(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID, false)
		}
		return ai.APICallLandlordWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
	}
}

// PlayCards AI出牌决策
func (s *AIService) PlayCards(player *models.GamePlayer, game *models.Game, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
	PassTurn(roomID, username string) error
}, roomID string) error {
	// 构建提示词
	prompt := s.promptBuilder.BuildPlayCardPrompt(player, game)

	// 构建消息
	messages := []ai.Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	// 根据提供商调用相应的API
	if s.provider == "gemini" {
		response, err := s.geminiClient.Chat(messages)
		if err != nil {
			// 如果API调用失败，使用默认策略（过牌）
			return ai.PassTurn(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
		}
		return ai.APIPlayCardsWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, player.Cards, gameService, roomID)
	} else {
		response, err := s.openAIClient.Chat(messages)
		if err != nil {
			// 如果API调用失败，使用默认策略（过牌）
			return ai.PassTurn(&ai.PlayerWrapper{UserName: player.UserName}, gameService, roomID)
		}
		return ai.APIPlayCardsWithResponse(response, &ai.PlayerWrapper{UserName: player.UserName}, player.Cards, gameService, roomID)
	}
}
