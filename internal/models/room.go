package models

import (
	"encoding/json"
	"time"
)

// RoomStatus 房间状态
type RoomStatus int

const (
	RoomStatusIdle   RoomStatus = 0 // 空闲
	RoomStatusWaiting RoomStatus = 1 // 等待玩家
	RoomStatusPlaying RoomStatus = 2 // 游戏中
)

var RoomStatusNames = map[RoomStatus]string{
	RoomStatusIdle:    "空闲",
	RoomStatusWaiting: "等待玩家",
	RoomStatusPlaying: "游戏中",
}

// RoomType 房间类型
type RoomType int

const (
	RoomTypePublic  RoomType = 0 // 公开房间
	RoomTypePrivate RoomType = 1 // 私人房间
)

var RoomTypeNames = map[RoomType]string{
	RoomTypePublic:  "公开房间",
	RoomTypePrivate: "私人房间",
}

// Room 房间对象
type Room struct {
	ID          string     `json:"id"`          // 房间ID
	Name        string     `json:"name"`        // 房间名称
	Owner       string     `json:"owner"`       // 房主
	Type        RoomType   `json:"type"`        // 房间类型
	Status      RoomStatus `json:"status"`      // 房间状态
	MaxPlayers  int        `json:"max_players"` // 最大玩家数
	Password    string     `json:"password"`    // 房间密码（私人房间）
	CreatedAt   time.Time  `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`  // 更新时间
	CurrentGame *Game      `json:"current_game"` // 当前游戏
}

// NewRoom 创建新房间
func NewRoom(id, name, owner string, roomType RoomType, password string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Owner:      owner,
		Type:       roomType,
		Status:     RoomStatusIdle,
		MaxPlayers: 3, // 斗地主固定3人
		Password:   password,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// StartGame 开始新游戏
func (r *Room) StartGame() *Game {
	gameID := r.ID + "_" + time.Now().Format("20060102150405")
	r.CurrentGame = NewGame(gameID, r.ID)
	r.Status = RoomStatusPlaying
	r.UpdatedAt = time.Now()
	return r.CurrentGame
}

// EndGame 结束当前游戏
func (r *Room) EndGame() {
	if r.CurrentGame != nil {
		now := time.Now()
		r.CurrentGame.FinishedAt = &now
		r.CurrentGame.Status = GameStatusFinished
	}
	r.Status = RoomStatusIdle
	r.UpdatedAt = time.Now()
}

// IsGameActive 检查是否有活跃的游戏
func (r *Room) IsGameActive() bool {
	return r.CurrentGame != nil &&
		   r.CurrentGame.Status != GameStatusFinished &&
		   r.CurrentGame.Status != GameStatusAbandoned
}

// GetPlayerCount 获取当前玩家数量
func (r *Room) GetPlayerCount() int {
	if r.CurrentGame == nil {
		return 0
	}
	count := 0
	for _, player := range r.CurrentGame.Players {
		if player != nil {
			count++
		}
	}
	return count
}

// IsFull 检查房间是否已满
func (r *Room) IsFull() bool {
	return r.GetPlayerCount() >= r.MaxPlayers
}

// CanJoin 检查玩家是否可以加入
func (r *Room) CanJoin(password string) bool {
	if r.IsFull() {
		return false
	}
	if r.Type == RoomTypePrivate && r.Password != password {
		return false
	}
	return true
}

// HasPlayer 检查是否包含指定玩家
func (r *Room) HasPlayer(username string) bool {
	if r.CurrentGame == nil {
		return false
	}
	return r.CurrentGame.GetPlayerByName(username) != nil
}

// ToJSON 转换为JSON字符串
func (r *Room) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	return string(data), err
}

// FromJSON 从JSON字符串创建房间对象
func (r *Room) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), r)
}

// GetSafeRoom 获取安全的房间信息（不包含敏感信息）
func (r *Room) GetSafeRoom() *Room {
	safeRoom := &Room{
		ID:         r.ID,
		Name:       r.Name,
		Owner:      r.Owner,
		Type:       r.Type,
		Status:     r.Status,
		MaxPlayers: r.MaxPlayers,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
	// 不返回密码和游戏详情
	return safeRoom
}