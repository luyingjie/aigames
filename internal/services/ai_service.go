package services

import (
	"aigames/internal/ai"
	"aigames/internal/config"
	"aigames/internal/models"
)

// AIService AI服务接口
type AIService struct {
	apiClient     *ai.APIClient
	promptBuilder *ai.PromptBuilder
}

// NewAIService 创建AI服务
func NewAIService() *AIService {
	// 创建AI API客户端和提示词构建器
	apiClient := ai.NewAPIClient(&config.GetConfig().AI)
	promptBuilder := ai.NewPromptBuilder()

	return &AIService{
		apiClient:     apiClient,
		promptBuilder: promptBuilder,
	}
}

// CallLandlord AI叫地主决策
func (s *AIService) CallLandlord(player *models.GamePlayer, game *models.Game, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string) error {
	return ai.APICallLandlord(s.apiClient, s.promptBuilder, player, game, gameService, roomID)
}

// PlayCards AI出牌决策
func (s *AIService) PlayCards(player *models.GamePlayer, game *models.Game, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
	PassTurn(roomID, username string) error
}, roomID string) error {
	return ai.APIPlayCards(s.apiClient, s.promptBuilder, player, game, gameService, roomID)
}
