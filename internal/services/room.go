package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"aigames/internal/models"

	"go.etcd.io/bbolt"
)

// RoomService 房间服务
type RoomService struct {
	db    *bbolt.DB
	rooms map[string]*models.Room // 内存中的房间缓存
	mutex sync.RWMutex            // 读写锁
}

// NewRoomService 创建房间服务实例
func NewRoomService(db *bbolt.DB) *RoomService {
	service := &RoomService{
		db:    db,
		rooms: make(map[string]*models.Room),
	}
	// 加载已存在的房间
	service.loadRoomsFromDB()
	return service
}

// loadRoomsFromDB 从数据库加载房间
func (rs *RoomService) loadRoomsFromDB() {
	rs.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rooms"))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var room models.Room
			if err := json.Unmarshal(v, &room); err == nil {
				rs.rooms[room.ID] = &room
			}
		}
		return nil
	})
}

// CreateRoom 创建房间
func (rs *RoomService) CreateRoom(id, name, owner string, roomType models.RoomType, password string, aiCount int) (*models.Room, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	// 检查房间ID是否已存在
	if _, exists := rs.rooms[id]; exists {
		return nil, fmt.Errorf("房间ID已存在")
	}

	// 创建新房间
	room := models.NewRoom(id, name, owner, roomType, password)
	rs.rooms[id] = room

	// 如果指定了AI玩家数量，自动创建AI玩家
	if aiCount > 0 {
		// 启动游戏以便添加AI玩家
		room.StartGame()
		game := room.CurrentGame

		// 创建AI玩家
		for i := 0; i < aiCount && i < 2; i++ { // 最多2个AI玩家
			aiName := fmt.Sprintf("AI-%d", i+1)

			// 找到空位置加入AI玩家
			for pos := models.Position1; pos <= models.Position3; pos++ {
				if game.GetPlayer(pos) == nil {
					if game.AddPlayer(aiName, pos) {
						// 将AI玩家标记为AI
						player := game.GetPlayer(pos)
						if player != nil {
							player.IsAI = true
							player.IsReady = true // AI玩家默认准备
						}
						break
					}
				}
			}
		}

		// 更新房间状态
		if room.GetPlayerCount() > 0 {
			room.Status = models.RoomStatusWaiting
		}
	}

	// 保存到数据库
	if err := rs.saveRoomToDB(room); err != nil {
		delete(rs.rooms, id) // 如果保存失败，从内存中删除
		return nil, fmt.Errorf("保存房间失败: %w", err)
	}

	return room, nil
}

// GetRoom 获取房间
func (rs *RoomService) GetRoom(id string) (*models.Room, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	room, exists := rs.rooms[id]
	if !exists {
		return nil, fmt.Errorf("房间不存在")
	}

	return room, nil
}

// GetAllRooms 获取所有房间列表
func (rs *RoomService) GetAllRooms() []*models.Room {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	rooms := make([]*models.Room, 0, len(rs.rooms))
	for _, room := range rs.rooms {
		// 返回安全的房间信息
		rooms = append(rooms, room.GetSafeRoom())
	}

	return rooms
}

// GetPublicRooms 获取公开房间列表
func (rs *RoomService) GetPublicRooms() []*models.Room {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	rooms := make([]*models.Room, 0)
	for _, room := range rs.rooms {
		if room.Type == models.RoomTypePublic {
			rooms = append(rooms, room.GetSafeRoom())
		}
	}

	return rooms
}

// JoinRoom 加入房间
func (rs *RoomService) JoinRoom(roomID, username, password string) (*models.Room, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("房间不存在")
	}

	// 检查是否可以加入
	if !room.CanJoin(password) {
		if room.IsFull() {
			return nil, fmt.Errorf("房间已满")
		}
		return nil, fmt.Errorf("房间密码错误")
	}

	// 检查玩家是否已在房间中
	if room.HasPlayer(username) {
		return room, nil // 已经在房间中
	}

	// 如果房间没有活跃游戏，创建新游戏
	if !room.IsGameActive() {
		room.StartGame()
	}

	// 找到空位置加入游戏
	game := room.CurrentGame
	var joinedPosition models.PlayerPosition = -1

	for i := models.Position1; i <= models.Position3; i++ {
		if game.GetPlayer(i) == nil {
			if game.AddPlayer(username, i) {
				joinedPosition = i
				break
			}
		}
	}

	if joinedPosition == -1 {
		return nil, fmt.Errorf("无法加入游戏")
	}

	// 更新房间状态
	if room.GetPlayerCount() > 0 {
		room.Status = models.RoomStatusWaiting
	}

	room.UpdatedAt = time.Now()

	// 保存到数据库
	if err := rs.saveRoomToDB(room); err != nil {
		return nil, fmt.Errorf("保存房间状态失败: %w", err)
	}

	return room, nil
}

// LeaveRoom 离开房间
func (rs *RoomService) LeaveRoom(roomID, username string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return fmt.Errorf("房间不存在")
	}

	if !room.HasPlayer(username) {
		return nil // 玩家不在房间中
	}

	// 从游戏中移除玩家
	if room.CurrentGame != nil {
		if position, found := room.CurrentGame.GetPlayerPosition(username); found {
			room.CurrentGame.RemovePlayer(position)
		}
	}

	// 更新房间状态
	playerCount := room.GetPlayerCount()
	if playerCount == 0 {
		room.Status = models.RoomStatusIdle
		// 如果游戏已经结束，可以清除当前游戏
		if room.CurrentGame != nil &&
			(room.CurrentGame.Status == models.GameStatusFinished ||
				room.CurrentGame.Status == models.GameStatusAbandoned) {
			room.CurrentGame = nil
		}
	} else {
		room.Status = models.RoomStatusWaiting
	}

	room.UpdatedAt = time.Now()

	// 如果房间是私人房间且房主离开且没有其他玩家，删除房间
	if room.Type == models.RoomTypePrivate && room.Owner == username && playerCount == 0 {
		return rs.DeleteRoom(roomID)
	}

	// 保存到数据库
	return rs.saveRoomToDB(room)
}

// DeleteRoom 删除房间
func (rs *RoomService) DeleteRoom(id string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[id]
	if !exists {
		return fmt.Errorf("房间不存在")
	}

	// 结束当前游戏
	if room.IsGameActive() {
		room.EndGame()
	}

	// 从内存中删除
	delete(rs.rooms, id)

	// 从数据库中删除
	return rs.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rooms"))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(id))
	})
}

// UpdateRoom 更新房间
func (rs *RoomService) UpdateRoom(room *models.Room) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if _, exists := rs.rooms[room.ID]; !exists {
		return fmt.Errorf("房间不存在")
	}

	room.UpdatedAt = time.Now()
	rs.rooms[room.ID] = room

	return rs.saveRoomToDB(room)
}

// saveRoomToDB 保存房间到数据库
func (rs *RoomService) saveRoomToDB(room *models.Room) error {
	return rs.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("rooms"))
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(room)
		if err != nil {
			return fmt.Errorf("序列化房间失败: %w", err)
		}

		return b.Put([]byte(room.ID), encoded)
	})
}

// StartGame 开始游戏
func (rs *RoomService) StartGame(roomID string) (*models.Game, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("房间不存在")
	}

	if room.CurrentGame == nil {
		return nil, fmt.Errorf("没有活跃的游戏")
	}

	game := room.CurrentGame

	// 检查是否所有玩家都准备
	if !game.IsAllReady() {
		return nil, fmt.Errorf("不是所有玩家都准备")
	}

	// 开始游戏
	game.Status = models.GameStatusReady
	now := time.Now()
	game.StartedAt = &now

	// 发牌
	gameLogic := models.NewGameLogic(game)
	if err := gameLogic.DealCards(); err != nil {
		return nil, fmt.Errorf("发牌失败: %w", err)
	}

	room.Status = models.RoomStatusPlaying
	room.UpdatedAt = time.Now()

	// 保存到数据库
	if err := rs.saveRoomToDB(room); err != nil {
		return nil, fmt.Errorf("保存游戏状态失败: %w", err)
	}

	return game, nil
}

// SetPlayerReady 设置玩家准备状态
func (rs *RoomService) SetPlayerReady(roomID, username string, ready bool) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return fmt.Errorf("房间不存在")
	}

	if room.CurrentGame == nil {
		return fmt.Errorf("没有活跃的游戏")
	}

	player := room.CurrentGame.GetPlayerByName(username)
	if player == nil {
		return fmt.Errorf("玩家不在游戏中")
	}

	player.IsReady = ready
	room.UpdatedAt = time.Now()

	// 保存到数据库
	return rs.saveRoomToDB(room)
}

// GetPlayerCount 获取在线玩家总数
func (rs *RoomService) GetPlayerCount() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	count := 0
	for _, room := range rs.rooms {
		count += room.GetPlayerCount()
	}

	return count
}

// GetRoomCount 获取房间总数
func (rs *RoomService) GetRoomCount() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	return len(rs.rooms)
}
