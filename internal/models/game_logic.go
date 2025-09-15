package models

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// HandType 牌型
type HandType int

const (
	HandTypeNone                 HandType = 0  // 无效牌型
	HandTypeSingle               HandType = 1  // 单牌
	HandTypePair                 HandType = 2  // 对子
	HandTypeTriple               HandType = 3  // 三张
	HandTypeTripleSingle         HandType = 4  // 三带一
	HandTypeTriplePair           HandType = 5  // 三带二
	HandTypeStraight             HandType = 6  // 顺子
	HandTypePairStraight         HandType = 7  // 连对
	HandTypeTripleStraight       HandType = 8  // 飞机
	HandTypeTripleStraightSingle HandType = 9  // 飞机带单牌
	HandTypeTripleStraightPair   HandType = 10 // 飞机带对子
	HandTypeBomb                 HandType = 11 // 炸弹
	HandTypeRocket               HandType = 12 // 火箭（双王）
)

var HandTypeNames = map[HandType]string{
	HandTypeNone:                 "无效牌型",
	HandTypeSingle:               "单牌",
	HandTypePair:                 "对子",
	HandTypeTriple:               "三张",
	HandTypeTripleSingle:         "三带一",
	HandTypeTriplePair:           "三带二",
	HandTypeStraight:             "顺子",
	HandTypePairStraight:         "连对",
	HandTypeTripleStraight:       "飞机",
	HandTypeTripleStraightSingle: "飞机带单牌",
	HandTypeTripleStraightPair:   "飞机带对子",
	HandTypeBomb:                 "炸弹",
	HandTypeRocket:               "火箭",
}

// HandPattern 牌型分析结果
type HandPattern struct {
	Type      HandType `json:"type"`       // 牌型
	MainCards []Card   `json:"main_cards"` // 主牌（决定大小的牌）
	SubCards  []Card   `json:"sub_cards"`  // 副牌（带的牌）
	Weight    int      `json:"weight"`     // 权重（用于比较大小）
	IsValid   bool     `json:"is_valid"`   // 是否有效
	Length    int      `json:"length"`     // 长度（连牌的长度）
}

// GameLogic 游戏逻辑
type GameLogic struct {
	game *Game
}

// NewGameLogic 创建游戏逻辑实例
func NewGameLogic(game *Game) *GameLogic {
	return &GameLogic{game: game}
}

// ShuffleDeck 洗牌
func ShuffleDeck(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// DealCards 发牌
func (gl *GameLogic) DealCards() error {
	if gl.game.Status != GameStatusReady {
		return fmt.Errorf("游戏状态不正确")
	}

	// 创建并洗牌
	deck := NewDeck()
	ShuffleDeck(deck)

	// 给每个玩家发17张牌
	for i := 0; i < 17; i++ {
		for j := 0; j < 3; j++ {
			if gl.game.Players[j] != nil {
				gl.game.Players[j].Cards = append(gl.game.Players[j].Cards, deck[i*3+j])
			}
		}
	}

	// 剩余3张作为地主牌
	gl.game.LandlordCards = deck[51:]

	// 为每个玩家排序手牌
	for _, player := range gl.game.Players {
		if player != nil {
			player.SortCards()
		}
	}

	gl.game.Status = GameStatusCalling
	gl.game.AddLog("deal", Position1, nil, "发牌完成")

	return nil
}

// CallLandlord 叫地主
func (gl *GameLogic) CallLandlord(position PlayerPosition, call bool) error {
	if gl.game.Status != GameStatusCalling {
		return fmt.Errorf("当前不是叫地主阶段")
	}

	player := gl.game.GetPlayer(position)
	if player == nil {
		return fmt.Errorf("玩家不存在")
	}

	if position != gl.game.CurrentTurn {
		return fmt.Errorf("不是该玩家的回合")
	}

	player.CallLandlord = true

	if call {
		// 成为地主
		player.Role = RoleLandlord
		for _, p := range gl.game.Players {
			if p != nil && p != player {
				p.Role = RoleFarmer
			}
		}

		// 地主获得底牌
		player.AddCards(gl.game.LandlordCards)

		gl.game.Status = GameStatusPlaying
		gl.game.CurrentTurn = position // 地主先出牌
		gl.game.AddLog("call_landlord", position, nil, fmt.Sprintf("%s 成为地主", player.UserName))

		return nil
	}

	// 不叫地主，轮到下一个玩家
	gl.game.NextTurn()

	// 检查是否所有人都不叫地主
	allCalled := true
	anyCalled := false
	for _, p := range gl.game.Players {
		if p != nil {
			if !p.CallLandlord {
				allCalled = false
			} else {
				// 这里简单处理：如果有人叫了地主，就让最后一个叫的人当地主
				if p.Role != RoleLandlord {
					anyCalled = true
				}
			}
		}
	}

	if allCalled && !anyCalled {
		// 所有人都不叫，游戏结束
		gl.game.Status = GameStatusAbandoned
		gl.game.AddLog("abandon", position, nil, "没有人叫地主，游戏结束")
	}

	return nil
}

// AnalyzeHand 分析手牌牌型
func AnalyzeHand(cards []Card) HandPattern {
	if len(cards) == 0 {
		return HandPattern{Type: HandTypeNone, IsValid: false}
	}

	// 排序
	sortedCards := make([]Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].GetWeight() < sortedCards[j].GetWeight()
	})

	// 统计每种牌值的数量
	countMap := make(map[CardValue]int)
	for _, card := range sortedCards {
		countMap[card.Value]++
	}

	// 按数量分组
	var singles, pairs, triples, bombs []CardValue
	for value, count := range countMap {
		switch count {
		case 1:
			singles = append(singles, value)
		case 2:
			pairs = append(pairs, value)
		case 3:
			triples = append(triples, value)
		case 4:
			bombs = append(bombs, value)
		}
	}

	// 排序各组
	sort.Slice(singles, func(i, j int) bool { return singles[i] < singles[j] })
	sort.Slice(pairs, func(i, j int) bool { return pairs[i] < pairs[j] })
	sort.Slice(triples, func(i, j int) bool { return triples[i] < triples[j] })
	sort.Slice(bombs, func(i, j int) bool { return bombs[i] < bombs[j] })

	// 分析具体牌型
	return analyzeHandType(sortedCards, singles, pairs, triples, bombs)
}

// analyzeHandType 具体分析牌型
func analyzeHandType(cards []Card, singles, pairs, triples, bombs []CardValue) HandPattern {
	cardCount := len(cards)

	// 火箭（双王）
	if cardCount == 2 &&
		len(singles) == 0 && len(pairs) == 0 && len(triples) == 0 && len(bombs) == 0 {
		hasSmallJoker := false
		hasBigJoker := false
		for _, card := range cards {
			if card.Value == ValueSmallJoker {
				hasSmallJoker = true
			} else if card.Value == ValueBigJoker {
				hasBigJoker = true
			}
		}
		if hasSmallJoker && hasBigJoker {
			return HandPattern{
				Type:      HandTypeRocket,
				MainCards: cards,
				Weight:    int(ValueBigJoker),
				IsValid:   true,
			}
		}
	}

	// 炸弹
	if len(bombs) == 1 && cardCount == 4 {
		return HandPattern{
			Type:      HandTypeBomb,
			MainCards: cards,
			Weight:    int(bombs[0]),
			IsValid:   true,
		}
	}

	// 单牌
	if cardCount == 1 {
		return HandPattern{
			Type:      HandTypeSingle,
			MainCards: cards,
			Weight:    int(cards[0].Value),
			IsValid:   true,
		}
	}

	// 对子
	if cardCount == 2 && len(pairs) == 1 {
		return HandPattern{
			Type:      HandTypePair,
			MainCards: cards,
			Weight:    int(pairs[0]),
			IsValid:   true,
		}
	}

	// 三张
	if cardCount == 3 && len(triples) == 1 {
		return HandPattern{
			Type:      HandTypeTriple,
			MainCards: cards,
			Weight:    int(triples[0]),
			IsValid:   true,
		}
	}

	// 三带一
	if cardCount == 4 && len(triples) == 1 && len(singles) == 1 {
		mainCards := make([]Card, 0)
		subCards := make([]Card, 0)
		for _, card := range cards {
			if card.Value == triples[0] {
				mainCards = append(mainCards, card)
			} else {
				subCards = append(subCards, card)
			}
		}
		return HandPattern{
			Type:      HandTypeTripleSingle,
			MainCards: mainCards,
			SubCards:  subCards,
			Weight:    int(triples[0]),
			IsValid:   true,
		}
	}

	// 三带二
	if cardCount == 5 && len(triples) == 1 && len(pairs) == 1 {
		mainCards := make([]Card, 0)
		subCards := make([]Card, 0)
		for _, card := range cards {
			if card.Value == triples[0] {
				mainCards = append(mainCards, card)
			} else {
				subCards = append(subCards, card)
			}
		}
		return HandPattern{
			Type:      HandTypeTriplePair,
			MainCards: mainCards,
			SubCards:  subCards,
			Weight:    int(triples[0]),
			IsValid:   true,
		}
	}

	// 顺子（至少5张连续的单牌，不能包含2和王）
	if cardCount >= 5 && len(singles) == cardCount {
		if isStraight(singles, false) {
			return HandPattern{
				Type:      HandTypeStraight,
				MainCards: cards,
				Weight:    int(singles[0]),
				IsValid:   true,
				Length:    len(singles),
			}
		}
	}

	// 连对（至少3对连续的对子，不能包含2和王）
	if cardCount >= 6 && cardCount%2 == 0 && len(pairs) == cardCount/2 {
		if isStraight(pairs, false) {
			return HandPattern{
				Type:      HandTypePairStraight,
				MainCards: cards,
				Weight:    int(pairs[0]),
				IsValid:   true,
				Length:    len(pairs),
			}
		}
	}

	// 飞机（连续的三张，至少2组）
	if len(triples) >= 2 && isStraight(triples, false) {
		tripleCount := len(triples)
		expectedCards := tripleCount * 3

		if cardCount == expectedCards {
			// 纯飞机
			return HandPattern{
				Type:      HandTypeTripleStraight,
				MainCards: cards,
				Weight:    int(triples[0]),
				IsValid:   true,
				Length:    tripleCount,
			}
		} else if cardCount == expectedCards+tripleCount && len(singles) == tripleCount {
			// 飞机带单牌
			mainCards := make([]Card, 0)
			subCards := make([]Card, 0)
			for _, card := range cards {
				isTriple := false
				for _, triple := range triples {
					if card.Value == triple {
						isTriple = true
						break
					}
				}
				if isTriple {
					mainCards = append(mainCards, card)
				} else {
					subCards = append(subCards, card)
				}
			}
			return HandPattern{
				Type:      HandTypeTripleStraightSingle,
				MainCards: mainCards,
				SubCards:  subCards,
				Weight:    int(triples[0]),
				IsValid:   true,
				Length:    tripleCount,
			}
		} else if cardCount == expectedCards+tripleCount*2 && len(pairs) == tripleCount {
			// 飞机带对子
			mainCards := make([]Card, 0)
			subCards := make([]Card, 0)
			for _, card := range cards {
				isTriple := false
				for _, triple := range triples {
					if card.Value == triple {
						isTriple = true
						break
					}
				}
				if isTriple {
					mainCards = append(mainCards, card)
				} else {
					subCards = append(subCards, card)
				}
			}
			return HandPattern{
				Type:      HandTypeTripleStraightPair,
				MainCards: mainCards,
				SubCards:  subCards,
				Weight:    int(triples[0]),
				IsValid:   true,
				Length:    tripleCount,
			}
		}
	}

	// 无效牌型
	return HandPattern{Type: HandTypeNone, IsValid: false}
}

// isStraight 检查是否为连续牌（顺子检查）
func isStraight(values []CardValue, allowJokers bool) bool {
	if len(values) < 2 {
		return false
	}

	for i := 1; i < len(values); i++ {
		// 不允许包含2和王
		if !allowJokers {
			if values[i-1] >= Value2 || values[i] >= Value2 {
				return false
			}
		}

		// 检查是否连续
		if int(values[i])-int(values[i-1]) != 1 {
			return false
		}
	}

	return true
}

// CanBeat 判断牌型A是否能打过牌型B
func CanBeat(a, b HandPattern) bool {
	// 无效牌型不能打过任何牌
	if !a.IsValid {
		return false
	}

	// 火箭能打过所有牌
	if a.Type == HandTypeRocket {
		return true
	}

	// 如果B是火箭，只有火箭能打过
	if b.Type == HandTypeRocket && a.Type != HandTypeRocket {
		return false
	}

	// 炸弹能打过非炸弹和非火箭的所有牌
	if a.Type == HandTypeBomb && b.Type != HandTypeBomb && b.Type != HandTypeRocket {
		return true
	}

	// 如果B是炸弹，只有更大的炸弹或火箭能打过
	if b.Type == HandTypeBomb && a.Type != HandTypeBomb && a.Type != HandTypeRocket {
		return false
	}

	// 同类型牌比较
	if a.Type != b.Type {
		return false
	}

	// 对于有长度要求的牌型，长度必须相同
	if (a.Type == HandTypeStraight || a.Type == HandTypePairStraight ||
		a.Type == HandTypeTripleStraight || a.Type == HandTypeTripleStraightSingle ||
		a.Type == HandTypeTripleStraightPair) && a.Length != b.Length {
		return false
	}

	// 比较权重
	return a.Weight > b.Weight
}
