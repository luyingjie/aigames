package ai

import (
	"aigames/internal/models"
	"aigames/pkg/logger"
	"strconv"
	"strings"
)

// PlayerInterface 定义玩家接口
type PlayerInterface interface {
	GetUserName() string
}

// PlayerWrapper 包装器结构
type PlayerWrapper struct {
	UserName string
}

// GetUserName 获取玩家用户名
func (p *PlayerWrapper) GetUserName() string {
	return p.UserName
}

// CallLandlord AI叫地主操作
func CallLandlord(player *PlayerWrapper, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string, call bool) error {
	logger.Info("AI玩家 %s 叫地主: %t", player.GetUserName(), call)
	return gameService.CallLandlord(roomID, player.GetUserName(), call)
}

// PassTurn AI过牌操作
func PassTurn(player *PlayerWrapper, gameService interface {
	PassTurn(roomID, username string) error
}, roomID string) error {
	logger.Info("AI玩家 %s 过牌", player.GetUserName())
	return gameService.PassTurn(roomID, player.GetUserName())
}

// PlayCards AI出牌操作
func PlayCards(player *PlayerWrapper, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
}, roomID string, cards []models.Card) error {
	logger.Info("AI玩家 %s 出牌: %v", player.GetUserName(), cards)
	return gameService.PlayCards(roomID, player.GetUserName(), cards)
}

// APICallLandlord 使用AI API进行叫地主决策
func APICallLandlord(apiClient interface {
	Chat(messages []Message) (string, error)
}, promptBuilder interface {
	BuildCallLandlordPrompt(player *models.GamePlayer, game *models.Game) string
}, player *models.GamePlayer, game *models.Game, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string) error {
	// 构建提示词
	prompt := promptBuilder.BuildCallLandlordPrompt(player, game)

	// 发送API请求
	messages := []Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	response, err := apiClient.Chat(messages)
	if err != nil {
		logger.Error("AI API调用失败: %v", err)
		// 如果API调用失败，使用默认策略（不叫地主）
		return CallLandlord(&PlayerWrapper{UserName: player.UserName}, gameService, roomID, false)
	}

	return APICallLandlordWithResponse(response, &PlayerWrapper{UserName: player.UserName}, gameService, roomID)
}

// APICallLandlordWithResponse 使用AI API响应进行叫地主决策
func APICallLandlordWithResponse(response string, player *PlayerWrapper, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string) error {
	// 解析响应
	call := strings.Contains(response, "叫地主")
	logger.Info("AI玩家 %s 叫地主决策: %t (API响应: %s)", player.GetUserName(), call, response)

	return CallLandlord(player, gameService, roomID, call)
}

// APIPlayCards 使用AI API进行出牌决策
func APIPlayCards(apiClient interface {
	Chat(messages []Message) (string, error)
}, promptBuilder interface {
	BuildPlayCardPrompt(player *models.GamePlayer, game *models.Game) string
}, player *models.GamePlayer, game *models.Game, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
	PassTurn(roomID, username string) error
}, roomID string) error {
	// 构建提示词
	prompt := promptBuilder.BuildPlayCardPrompt(player, game)

	// 发送API请求
	messages := []Message{
		{Role: "system", Content: "你是一个专业的斗地主游戏AI玩家。"},
		{Role: "user", Content: prompt},
	}

	response, err := apiClient.Chat(messages)
	if err != nil {
		logger.Error("AI API调用失败: %v", err)
		// 如果API调用失败，使用默认策略（过牌）
		return PassTurn(&PlayerWrapper{UserName: player.UserName}, gameService, roomID)
	}

	return APIPlayCardsWithResponse(response, &PlayerWrapper{UserName: player.UserName}, player.Cards, gameService, roomID)
}

// APIPlayCardsWithResponse 使用AI API响应进行出牌决策
func APIPlayCardsWithResponse(response string, player *PlayerWrapper, handCards []models.Card, gameService interface {
	PlayCards(roomID, username string, cards []models.Card) error
	PassTurn(roomID, username string) error
}, roomID string) error {
	logger.Info("AI玩家 %s 出牌决策 (API响应): %s", player.GetUserName(), response)

	// 解析响应
	if strings.Contains(response, "过牌") {
		return PassTurn(player, gameService, roomID)
	}

	// 解析出牌
	cards := parseCardsFromResponse(response, handCards)
	if len(cards) == 0 {
		// 如果无法解析出牌，选择过牌
		return PassTurn(player, gameService, roomID)
	}

	return PlayCards(player, gameService, roomID, cards)
}

// parseCardsFromResponse 从API响应中解析出牌
func parseCardsFromResponse(response string, handCards []models.Card) []models.Card {
	// 查找"出牌:"部分
	if strings.Contains(response, "出牌:") {
		// 提取牌型部分
		parts := strings.Split(response, "出牌:")
		if len(parts) > 1 {
			// 移除可能的标点符号
			cardStr := strings.TrimSpace(parts[1])
			cardStr = strings.Trim(cardStr, " .,;!?")

			// 分割牌
			cardParts := strings.Split(cardStr, ",")
			var result []models.Card

			// 遍历手牌，匹配API返回的牌
			for _, handCard := range handCards {
				cardName := handCard.String()
				for _, part := range cardParts {
					part = strings.TrimSpace(part)
					// 检查是否匹配
					if strings.Contains(cardName, part) || strings.Contains(part, cardName) {
						result = append(result, handCard)
						break
					}
				}
			}

			// 如果没有匹配到，尝试按牌值匹配
			if len(result) == 0 {
				for _, handCard := range handCards {
					valueStr := strconv.Itoa(int(handCard.Value))
					for _, part := range cardParts {
						part = strings.TrimSpace(part)
						if part == valueStr || part == models.ValueNames[handCard.Value] {
							result = append(result, handCard)
							break
						}
					}
				}
			}

			return result
		}
	}
	return []models.Card{}
}
