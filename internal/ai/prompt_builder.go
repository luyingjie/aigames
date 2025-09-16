package ai

import (
	"fmt"
	"strings"

	"aigames/internal/models"
)

// PromptBuilder 提示词构建器
type PromptBuilder struct {
	provider string
}

// NewPromptBuilder 创建提示词构建器
func NewPromptBuilder() *PromptBuilder {
	// 这里可以获取实际的提供商配置
	// 为了简化，我们使用默认值
	return &PromptBuilder{
		provider: "gemini", // 默认使用Gemini
	}
}

// BuildGameRulesPrompt 构建游戏规则提示词
func (pb *PromptBuilder) BuildGameRulesPrompt() string {
	basePrompt := `你是一个专业的斗地主AI。
规则:
- 牌序: 3,4,5,6,7,8,9,10,J,Q,K,A,2,小王,大王
- 牌型: 单,对,三,三带一,三带二,顺子,连对,飞机,炸弹,火箭
- 目标: 地主先出完牌获胜，农民合作阻止地主。
- 出牌: 你必须出比上家大的牌，或者"过牌"。
- 格式: 严格按照'叫地主'/'不叫'或'出牌:牌'/'过牌'格式回答。例如: '出牌:3,4,5,6,7'
`

	return basePrompt
}

// BuildCallLandlordPrompt 构建叫地主提示词
func (pb *PromptBuilder) BuildCallLandlordPrompt(player *models.GamePlayer, game *models.Game) string {
	var sb strings.Builder

	sb.WriteString(pb.BuildGameRulesPrompt())
	sb.WriteString("\n\n现在是叫地主阶段，请根据以下信息决定是否叫地主：\n\n")

	sb.WriteString(fmt.Sprintf("你的玩家信息：\n"))
	sb.WriteString(fmt.Sprintf("- 用户名：%s\n", player.UserName))
	sb.WriteString(fmt.Sprintf("- 手牌数量：%d\n", player.GetCardCount()))
	sb.WriteString(fmt.Sprintf("- 手牌：%s\n", pb.formatCards(player.Cards)))

	sb.WriteString("\n其他玩家信息：\n")
	for _, p := range game.Players {
		if p != nil && p.UserName != player.UserName {
			sb.WriteString(fmt.Sprintf("- %s，手牌数量：%d\n", p.UserName, p.GetCardCount()))
		}
	}

	if pb.provider == "gemini" {
		sb.WriteString("\n请分析你的手牌，判断是否叫地主。")
		sb.WriteString("\n你的回答必须严格遵循以下格式之一：'叫地主' 或 '不叫'。不要包含任何解释或额外的文字。")
	} else {
		sb.WriteString("\n请分析你的手牌，判断是否叫地主。")
		sb.WriteString("\n你的回答必须严格遵循以下格式之一：'叫地主' 或 '不叫'。不要包含任何解释或额外的文字。")
	}

	return sb.String()
}

// BuildPlayCardPrompt 构建出牌提示词
func (pb *PromptBuilder) BuildPlayCardPrompt(player *models.GamePlayer, game *models.Game) string {
	var sb strings.Builder

	sb.WriteString(pb.BuildGameRulesPrompt())
	sb.WriteString("\n\n现在是出牌阶段，请根据以下信息决定出什么牌：\n\n")

	sb.WriteString(fmt.Sprintf("你的玩家信息：\n"))
	sb.WriteString(fmt.Sprintf("- 用户名：%s\n", player.UserName))
	sb.WriteString(fmt.Sprintf("- 角色：%s\n", models.RoleNames[player.Role]))
	sb.WriteString(fmt.Sprintf("- 手牌：%s\n", pb.formatCards(player.Cards)))

	sb.WriteString("\n游戏状态：\n")
	landlord := game.GetLandlord()
	if landlord != nil {
		sb.WriteString(fmt.Sprintf("- 地主是: %s\n", landlord.UserName))
	}
	sb.WriteString(fmt.Sprintf("- 当前回合：%s\n", getPlayerNameByPosition(game, game.CurrentTurn)))
	sb.WriteString(fmt.Sprintf("- 上一手牌（%s）：%s\n", getPlayerNameByPosition(game, game.LastPlayer), pb.formatCards(game.LastPlayCards)))

	sb.WriteString("\n其他玩家信息：\n")
	for _, p := range game.Players {
		if p != nil && p.UserName != player.UserName {
			sb.WriteString(fmt.Sprintf("- %s (%s)：%d张\n", p.UserName, models.RoleNames[p.Role], p.GetCardCount()))
		}
	}

	if pb.provider == "gemini" {
		sb.WriteString("\n请分析当前情况，决定出什么牌。")
		sb.WriteString("\n你的回答必须严格遵循以下格式之一：'过牌' 或 '出牌:牌1,牌2,牌3'。不要包含任何解释或额外的文字。")
		sb.WriteString("\n注意：请确保选择的牌在你的手牌中，并且符合斗地主的出牌规则。")
	} else {
		sb.WriteString("\n请分析当前情况，决定出什么牌。")
		sb.WriteString("\n你的回答必须严格遵循以下格式之一：'过牌' 或 '出牌:牌1,牌2,牌3'。不要包含任何解释或额外的文字。")
	}

	return sb.String()
}

// formatCards 格式化牌列表
func (pb *PromptBuilder) formatCards(cards []models.Card) string {
	if len(cards) == 0 {
		return "无"
	}

	cardStrings := make([]string, len(cards))
	for i, card := range cards {
		cardStrings[i] = card.String()
	}

	return strings.Join(cardStrings, ", ")
}

// getPlayerNameByPosition 根据位置获取玩家名
func getPlayerNameByPosition(game *models.Game, position models.PlayerPosition) string {
	player := game.GetPlayer(position)
	if player == nil {
		return "未知玩家"
	}
	return player.UserName
}

// SetProvider 设置AI提供商
func (pb *PromptBuilder) SetProvider(provider string) {
	pb.provider = provider
}
