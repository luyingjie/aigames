package ai

import (
	"fmt"
	"strings"

	"aigames/internal/models"
)

// PromptBuilder 提示词构建器
type PromptBuilder struct{}

// NewPromptBuilder 创建提示词构建器
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// BuildGameRulesPrompt 构建游戏规则提示词
func (pb *PromptBuilder) BuildGameRulesPrompt() string {
	return `你是一个专业的斗地主游戏AI玩家。请遵循以下规则：

1. 游戏基本规则：
   - 使用一副54张的扑克牌（包括大小王）
   - 三名玩家，其中一名地主，两名农民
   - 地主先出牌，按逆时针顺序出牌
   - 目标是尽快出完手中的牌

2. 角色关系：
   - 地主：独立对抗两名农民，需要独自出完所有牌
   - 农民：合作对抗地主，需要阻止地主先出完牌

3. 牌型规则：
   - 单牌：任意一张牌
   - 对子：两张相同点数的牌
   - 三张：三张相同点数的牌
   - 三带一：三张相同点数的牌带一张单牌
   - 三带二：三张相同点数的牌带一对牌
   - 顺子：五张或更多连续的单牌（不能包含2和王）
   - 连对：三对或更多连续的对子（不能包含2和王）
   - 飞机：两个或更多连续的三张（不能包含2和王）
   - 飞机带单牌：飞机带同等数量的单牌
   - 飞机带对子：飞机带同等数量的对子
   - 炸弹：四张相同点数的牌
   - 火箭：双王（大王和小王），最大牌型

4. 牌的大小顺序（从小到大）：
   3 < 4 < 5 < 6 < 7 < 8 < 9 < 10 < J < Q < K < A < 2 < 小王 < 大王

5. 出牌规则：
   - 第一手牌可以出任意合法牌型
   - 后续出牌必须能压过上一手牌（牌型相同且更大，或使用炸弹/火箭）
   - 如果无法或不想压过上一手牌，可以选择"过牌"
   - 一轮中所有其他玩家都过牌后，最后出牌的玩家可以出任意合法牌型

6. 决策策略：
   - 地主目标：尽快出完所有牌
   - 农民目标：合作阻止地主，农民之间需要默契配合
   - 优先保留关键牌型（如炸弹）在关键时刻使用
   - 根据当前手牌和已出牌情况，计算最优出牌策略`
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

	sb.WriteString("\n请分析你的手牌，判断是否叫地主。")
	sb.WriteString("\n回答格式：只回答'叫地主'或'不叫'，不要包含其他内容。")

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
	sb.WriteString(fmt.Sprintf("- 当前回合：%s\n", getPlayerNameByPosition(game, game.CurrentTurn)))
	sb.WriteString(fmt.Sprintf("- 上一手牌（%s）：%s\n", getPlayerNameByPosition(game, game.LastPlayer), pb.formatCards(game.LastPlayCards)))

	sb.WriteString("\n其他玩家手牌数量：\n")
	for _, p := range game.Players {
		if p != nil && p.UserName != player.UserName {
			sb.WriteString(fmt.Sprintf("- %s：%d张\n", p.UserName, p.GetCardCount()))
		}
	}

	sb.WriteString("\n请分析当前情况，决定出什么牌。")
	sb.WriteString("\n回答格式：只回答'过牌'或'出牌:牌型'，例如'出牌:3,4,5'，不要包含其他内容。")

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
