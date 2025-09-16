package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// GameStatus 游戏状态
type GameStatus int

const (
	GameStatusWaiting    GameStatus = 0 // 等待玩家
	GameStatusReady      GameStatus = 1 // 准备开始
	GameStatusDealing    GameStatus = 2 // 发牌中
	GameStatusCalling    GameStatus = 3 // 叫地主
	GameStatusPlaying    GameStatus = 4 // 游戏进行中
	GameStatusFinished   GameStatus = 5 // 游戏结束
	GameStatusAbandoned  GameStatus = 6 // 游戏中止
)

var GameStatusNames = map[GameStatus]string{
	GameStatusWaiting:    "等待玩家",
	GameStatusReady:      "准备开始",
	GameStatusDealing:    "发牌中",
	GameStatusCalling:    "叫地主",
	GameStatusPlaying:    "游戏进行中",
	GameStatusFinished:   "游戏结束",
	GameStatusAbandoned:  "游戏中止",
}

// PlayerRole 玩家角色
type PlayerRole int

const (
	RoleNone     PlayerRole = 0 // 无角色
	RoleLandlord PlayerRole = 1 // 地主
	RoleFarmer   PlayerRole = 2 // 农民
)

var RoleNames = map[PlayerRole]string{
	RoleNone:     "观众",
	RoleLandlord: "地主",
	RoleFarmer:   "农民",
}

// PlayerPosition 玩家位置
type PlayerPosition int

const (
	Position1 PlayerPosition = 0 // 位置1
	Position2 PlayerPosition = 1 // 位置2
	Position3 PlayerPosition = 2 // 位置3
)

// GamePlayer 游戏中的玩家
type GamePlayer struct {
	UserName     string         `json:"username"`      // 用户名
	Position     PlayerPosition `json:"position"`      // 位置
	Role         PlayerRole     `json:"role"`          // 角色
	Cards        []Card         `json:"cards"`         // 手牌
	IsReady      bool           `json:"is_ready"`      // 是否准备
	IsOnline     bool           `json:"is_online"`     // 是否在线
	IsAI         bool           `json:"is_ai"`         // 是否为AI玩家
	Score        int            `json:"score"`         // 得分
	CallLandlord bool           `json:"call_landlord"` // 是否叫过地主
}

// GetCardCount 获取手牌数量
func (p *GamePlayer) GetCardCount() int {
	return len(p.Cards)
}

// AddCards 添加手牌
func (p *GamePlayer) AddCards(cards []Card) {
	p.Cards = append(p.Cards, cards...)
	p.SortCards()
}

// RemoveCards 移除指定的牌
func (p *GamePlayer) RemoveCards(cards []Card) bool {
	// 创建一个副本用于操作
	remainingCards := make([]Card, 0, len(p.Cards))
	cardsToRemove := make([]Card, len(cards))
	copy(cardsToRemove, cards)

	// 遍历手牌，移除指定的牌
	for _, card := range p.Cards {
		removed := false
		for i, toRemove := range cardsToRemove {
			if card.Suit == toRemove.Suit && card.Value == toRemove.Value {
				// 移除这张牌（通过将其标记为已处理）
				cardsToRemove = append(cardsToRemove[:i], cardsToRemove[i+1:]...)
				removed = true
				break
			}
		}
		if !removed {
			remainingCards = append(remainingCards, card)
		}
	}

	// 检查是否所有要移除的牌都找到了
	if len(cardsToRemove) > 0 {
		return false // 没有找到所有要移除的牌
	}

	p.Cards = remainingCards
	return true
}

// SortCards 对手牌排序
func (p *GamePlayer) SortCards() {
	sort.Slice(p.Cards, func(i, j int) bool {
		return p.Cards[i].GetWeight() < p.Cards[j].GetWeight()
	})
}

// HasCards 检查是否有指定的牌
func (p *GamePlayer) HasCards(cards []Card) bool {
	cardMap := make(map[string]int)

	// 统计手牌
	for _, card := range p.Cards {
		key := fmt.Sprintf("%d-%d", card.Suit, card.Value)
		cardMap[key]++
	}

	// 检查是否有足够的牌
	for _, card := range cards {
		key := fmt.Sprintf("%d-%d", card.Suit, card.Value)
		if cardMap[key] <= 0 {
			return false
		}
		cardMap[key]--
	}

	return true
}

// Game 游戏对象
type Game struct {
	ID             string            `json:"id"`              // 游戏ID
	RoomID         string            `json:"room_id"`         // 房间ID
	Status         GameStatus        `json:"status"`          // 游戏状态
	Players        [3]*GamePlayer    `json:"players"`         // 玩家列表（固定3个位置）
	LandlordCards  []Card            `json:"landlord_cards"`  // 地主牌（底牌）
	CurrentTurn    PlayerPosition    `json:"current_turn"`    // 当前回合
	LastPlayCards  []Card            `json:"last_play_cards"` // 上一次出的牌
	LastPlayer     PlayerPosition    `json:"last_player"`     // 上次出牌的玩家
	CreatedAt      time.Time         `json:"created_at"`      // 创建时间
	StartedAt      *time.Time        `json:"started_at"`      // 开始时间
	FinishedAt     *time.Time        `json:"finished_at"`     // 结束时间
	Winner         PlayerPosition    `json:"winner"`          // 获胜者
	GameLog        []GameLogEntry    `json:"game_log"`        // 游戏日志
}

// GameLogEntry 游戏日志条目
type GameLogEntry struct {
	Type      string         `json:"type"`      // 日志类型
	Player    PlayerPosition `json:"player"`    // 玩家位置
	Cards     []Card         `json:"cards"`     // 相关卡牌
	Message   string         `json:"message"`   // 消息内容
	Timestamp time.Time      `json:"timestamp"` // 时间戳
}

// NewGame 创建新游戏
func NewGame(gameID, roomID string) *Game {
	return &Game{
		ID:            gameID,
		RoomID:        roomID,
		Status:        GameStatusWaiting,
		Players:       [3]*GamePlayer{nil, nil, nil},
		LandlordCards: make([]Card, 0, 3),
		CurrentTurn:   Position1,
		CreatedAt:     time.Now(),
		GameLog:       make([]GameLogEntry, 0),
	}
}

// AddPlayer 添加玩家到指定位置
func (g *Game) AddPlayer(username string, position PlayerPosition) bool {
	if position < Position1 || position > Position3 {
		return false
	}
	if g.Players[position] != nil {
		return false // 位置已被占用
	}

	g.Players[position] = &GamePlayer{
		UserName:     username,
		Position:     position,
		Role:         RoleNone,
		Cards:        make([]Card, 0, 20),
		IsReady:      false,
		IsOnline:     true,
		Score:        0,
		CallLandlord: false,
	}

	g.AddLog("join", position, nil, fmt.Sprintf("玩家 %s 加入游戏", username))
	return true
}

// RemovePlayer 移除指定位置的玩家
func (g *Game) RemovePlayer(position PlayerPosition) {
	if position < Position1 || position > Position3 {
		return
	}
	if g.Players[position] != nil {
		username := g.Players[position].UserName
		g.Players[position] = nil
		g.AddLog("leave", position, nil, fmt.Sprintf("玩家 %s 离开游戏", username))
	}
}

// GetPlayer 获取指定位置的玩家
func (g *Game) GetPlayer(position PlayerPosition) *GamePlayer {
	if position < Position1 || position > Position3 {
		return nil
	}
	return g.Players[position]
}

// GetPlayerByName 根据用户名获取玩家
func (g *Game) GetPlayerByName(username string) *GamePlayer {
	for _, player := range g.Players {
		if player != nil && player.UserName == username {
			return player
		}
	}
	return nil
}

// GetPlayerPosition 获取玩家位置
func (g *Game) GetPlayerPosition(username string) (PlayerPosition, bool) {
	for i, player := range g.Players {
		if player != nil && player.UserName == username {
			return PlayerPosition(i), true
		}
	}
	return Position1, false
}

// IsPlayerFull 检查玩家是否已满
func (g *Game) IsPlayerFull() bool {
	count := 0
	for _, player := range g.Players {
		if player != nil {
			count++
		}
	}
	return count >= 3
}

// IsAllReady 检查是否所有玩家都准备
func (g *Game) IsAllReady() bool {
	if !g.IsPlayerFull() {
		return false
	}
	for _, player := range g.Players {
		if player != nil && !player.IsReady {
			return false
		}
	}
	return true
}

// GetLandlord 获取地主玩家
func (g *Game) GetLandlord() *GamePlayer {
	for _, player := range g.Players {
		if player != nil && player.Role == RoleLandlord {
			return player
		}
	}
	return nil
}

// NextTurn 下一个回合
func (g *Game) NextTurn() {
	g.CurrentTurn = (g.CurrentTurn + 1) % 3
}

// AddLog 添加游戏日志
func (g *Game) AddLog(logType string, player PlayerPosition, cards []Card, message string) {
	entry := GameLogEntry{
		Type:      logType,
		Player:    player,
		Cards:     cards,
		Message:   message,
		Timestamp: time.Now(),
	}
	g.GameLog = append(g.GameLog, entry)
}

// ToJSON 转换为JSON字符串
func (g *Game) ToJSON() (string, error) {
	data, err := json.Marshal(g)
	return string(data), err
}

// FromJSON 从JSON字符串创建游戏对象
func (g *Game) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), g)
}