package protocol

import "aigames/internal/models"

// 房间相关的请求和响应

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	BaseRequest
	Name     string          `json:"name" validate:"required,min=1,max=50"` // 房间名称
	Type     models.RoomType `json:"type"`                                  // 房间类型
	Password string          `json:"password,omitempty" validate:"max=20"`  // 房间密码（可选）
	AICount  int             `json:"ai_count" validate:"min=0,max=2"`       // AI玩家数量
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	BaseRequest
	RoomID   string `json:"room_id" validate:"required"` // 房间ID
	Password string `json:"password,omitempty"`          // 房间密码（如果需要）
}

// LeaveRoomRequest 离开房间请求
type LeaveRoomRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// GetRoomListRequest 获取房间列表请求
type GetRoomListRequest struct {
	PageRequest
	Type models.RoomType `json:"type,omitempty"` // 房间类型过滤
}

// SetReadyRequest 设置准备状态请求
type SetReadyRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
	Ready  bool   `json:"ready"`                       // 准备状态
}

// StartGameRequest 开始游戏请求
type StartGameRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// DeleteRoomRequest 删除房间请求
type DeleteRoomRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// 游戏相关的请求和响应

// CallLandlordRequest 叫地主请求
type CallLandlordRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
	Call   bool   `json:"call"`                        // 是否叫地主
}

// PlayCardsRequest 出牌请求
type PlayCardsRequest struct {
	BaseRequest
	RoomID string        `json:"room_id" validate:"required"` // 房间ID
	Cards  []models.Card `json:"cards" validate:"required"`   // 出的牌
}

// PassTurnRequest 过牌请求
type PassTurnRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// GetGameStateRequest 获取游戏状态请求
type GetGameStateRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// GetPlayerHandRequest 获取玩家手牌请求
type GetPlayerHandRequest struct {
	BaseRequest
	RoomID string `json:"room_id" validate:"required"` // 房间ID
}

// 响应数据结构

// RoomData 房间数据
type RoomData struct {
	ID          string            `json:"id"`           // 房间ID
	Name        string            `json:"name"`         // 房间名称
	Owner       string            `json:"owner"`        // 房主
	Type        models.RoomType   `json:"type"`         // 房间类型
	TypeName    string            `json:"type_name"`    // 房间类型名称
	Status      models.RoomStatus `json:"status"`       // 房间状态
	StatusName  string            `json:"status_name"`  // 房间状态名称
	MaxPlayers  int               `json:"max_players"`  // 最大玩家数
	PlayerCount int               `json:"player_count"` // 当前玩家数
	HasPassword bool              `json:"has_password"` // 是否有密码
	CreatedAt   string            `json:"created_at"`   // 创建时间
	UpdatedAt   string            `json:"updated_at"`   // 更新时间
}

// GameStateData 游戏状态数据
type GameStateData struct {
	GameID        string                `json:"game_id"`         // 游戏ID
	Status        models.GameStatus     `json:"status"`          // 游戏状态
	StatusName    string                `json:"status_name"`     // 游戏状态名称
	CurrentTurn   models.PlayerPosition `json:"current_turn"`    // 当前回合
	LastPlayCards []models.Card         `json:"last_play_cards"` // 上一次出的牌
	LastPlayer    models.PlayerPosition `json:"last_player"`     // 上次出牌的玩家
	Players       []interface{}         `json:"players"`         // 玩家信息
	LandlordCards []models.Card         `json:"landlord_cards"`  // 地主牌
	CreatedAt     string                `json:"created_at"`      // 创建时间
	StartedAt     *string               `json:"started_at"`      // 开始时间
	FinishedAt    *string               `json:"finished_at"`     // 结束时间
	Winner        models.PlayerPosition `json:"winner"`          // 获胜者
}

// PlayerHandData 玩家手牌数据
type PlayerHandData struct {
	Cards []models.Card `json:"cards"` // 手牌
}

// RoomListData 房间列表数据
type RoomListData struct {
	Rooms []RoomData `json:"rooms"` // 房间列表
}

// 快捷响应方法

// CreateRoomSuccess 创建房间成功响应
func CreateRoomSuccess(room *models.Room) BaseResponse {
	data := RoomData{
		ID:          room.ID,
		Name:        room.Name,
		Owner:       room.Owner,
		Type:        room.Type,
		TypeName:    models.RoomTypeNames[room.Type],
		Status:      room.Status,
		StatusName:  models.RoomStatusNames[room.Status],
		MaxPlayers:  room.MaxPlayers,
		PlayerCount: room.GetPlayerCount(),
		HasPassword: room.Password != "",
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   room.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return SuccessWithMessage(data, "创建房间成功")
}

// JoinRoomSuccess 加入房间成功响应
func JoinRoomSuccess(room *models.Room) BaseResponse {
	data := RoomData{
		ID:          room.ID,
		Name:        room.Name,
		Owner:       room.Owner,
		Type:        room.Type,
		TypeName:    models.RoomTypeNames[room.Type],
		Status:      room.Status,
		StatusName:  models.RoomStatusNames[room.Status],
		MaxPlayers:  room.MaxPlayers,
		PlayerCount: room.GetPlayerCount(),
		HasPassword: room.Password != "",
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   room.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return SuccessWithMessage(data, "加入房间成功")
}

// LeaveRoomSuccess 离开房间成功响应
func LeaveRoomSuccess() BaseResponse {
	return SuccessWithMessage(nil, "离开房间成功")
}

// RoomListSuccess 房间列表成功响应
func RoomListSuccess(rooms []*models.Room, total, page, size int) PageResponse {
	roomDataList := make([]RoomData, len(rooms))
	for i, room := range rooms {
		roomDataList[i] = RoomData{
			ID:          room.ID,
			Name:        room.Name,
			Owner:       room.Owner,
			Type:        room.Type,
			TypeName:    models.RoomTypeNames[room.Type],
			Status:      room.Status,
			StatusName:  models.RoomStatusNames[room.Status],
			MaxPlayers:  room.MaxPlayers,
			PlayerCount: room.GetPlayerCount(),
			HasPassword: room.Password != "",
			CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   room.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	data := RoomListData{Rooms: roomDataList}

	return PageResponse{
		BaseResponse: SuccessWithMessage(data, "获取房间列表成功"),
		Total:        total,
		Page:         page,
		Size:         size,
	}
}

// SetReadySuccess 设置准备状态成功响应
func SetReadySuccess() BaseResponse {
	return SuccessWithMessage(nil, "设置准备状态成功")
}

// StartGameSuccess 开始游戏成功响应
func StartGameSuccess() BaseResponse {
	return SuccessWithMessage(nil, "游戏开始")
}

// DeleteRoomSuccess 删除房间成功响应
func DeleteRoomSuccess() BaseResponse {
	return SuccessWithMessage(nil, "删除房间成功")
}

// CallLandlordSuccess 叫地主成功响应
func CallLandlordSuccess() BaseResponse {
	return SuccessWithMessage(nil, "操作成功")
}

// PlayCardsSuccess 出牌成功响应
func PlayCardsSuccess() BaseResponse {
	return SuccessWithMessage(nil, "出牌成功")
}

// PassTurnSuccess 过牌成功响应
func PassTurnSuccess() BaseResponse {
	return SuccessWithMessage(nil, "过牌成功")
}

// GameStateSuccess 获取游戏状态成功响应
func GameStateSuccess(gameState map[string]interface{}) BaseResponse {
	return SuccessWithMessage(gameState, "获取游戏状态成功")
}

// PlayerHandSuccess 获取玩家手牌成功响应
func PlayerHandSuccess(cards []models.Card) BaseResponse {
	data := PlayerHandData{Cards: cards}
	return SuccessWithMessage(data, "获取手牌成功")
}

// 游戏相关错误响应

// RoomNotFound 房间不存在
func RoomNotFound() BaseResponse {
	return ErrorWithCode(StatusRoomNotFound)
}

// RoomFull 房间已满
func RoomFull() BaseResponse {
	return ErrorWithCode(StatusRoomFull)
}

// GameNotFound 游戏不存在
func GameNotFound() BaseResponse {
	return ErrorWithCode(StatusGameNotFound)
}

// GameNotStarted 游戏未开始
func GameNotStarted() BaseResponse {
	return ErrorWithCode(StatusGameNotStarted)
}

// GameEnded 游戏已结束
func GameEnded() BaseResponse {
	return ErrorWithCode(StatusGameEnded)
}

// PlayerNotInRoom 玩家不在房间内
func PlayerNotInRoom() BaseResponse {
	return ErrorWithCode(StatusPlayerNotInRoom)
}

// NotPlayerTurn 不是玩家回合
func NotPlayerTurn() BaseResponse {
	return ErrorWithCode(StatusNotPlayerTurn)
}

// InvalidMove 无效操作
func InvalidMove() BaseResponse {
	return ErrorWithCode(StatusInvalidMove)
}
