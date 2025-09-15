package services

import (
	"fmt"
	"sync"
	"time"

	"aigames/internal/models"
	"aigames/pkg/logger"

	"go.etcd.io/bbolt"
)

// GameService 游戏服务
type GameService struct {
	db         *bbolt.DB
	roomService *RoomService
	games      map[string]*models.Game // 内存中的游戏缓存
	mutex      sync.RWMutex           // 读写锁
}

// NewGameService 创建游戏服务实例
func NewGameService(db *bbolt.DB, roomService *RoomService) *GameService {
	return &GameService{
		db:         db,
		roomService: roomService,
		games:      make(map[string]*models.Game),
	}
}

// GetGame 获取游戏
func (gs *GameService) GetGame(gameID string) (*models.Game, error) {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	game, exists := gs.games[gameID]
	if !exists {
		return nil, fmt.Errorf("游戏不存在")
	}

	return game, nil
}

// GetGameByRoom 通过房间ID获取游戏
func (gs *GameService) GetGameByRoom(roomID string) (*models.Game, error) {
	room, err := gs.roomService.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	if room.CurrentGame == nil {
		return nil, fmt.Errorf("房间没有活跃的游戏")
	}

	// 将游戏添加到缓存中
	gs.mutex.Lock()
	gs.games[room.CurrentGame.ID] = room.CurrentGame
	gs.mutex.Unlock()

	return room.CurrentGame, nil
}

// CallLandlord 叫地主
func (gs *GameService) CallLandlord(roomID, username string, call bool) error {
	game, err := gs.GetGameByRoom(roomID)
	if err != nil {
		return err
	}

	position, found := game.GetPlayerPosition(username)
	if !found {
		return fmt.Errorf("玩家不在游戏中")
	}

	gameLogic := models.NewGameLogic(game)
	if err := gameLogic.CallLandlord(position, call); err != nil {
		return err
	}

	// 更新房间状态
	room, _ := gs.roomService.GetRoom(roomID)
	return gs.roomService.UpdateRoom(room)
}

// PlayCards 出牌
func (gs *GameService) PlayCards(roomID, username string, cards []models.Card) error {
	game, err := gs.GetGameByRoom(roomID)
	if err != nil {
		return err
	}

	if game.Status != models.GameStatusPlaying {
		return fmt.Errorf("游戏状态不正确")
	}

	position, found := game.GetPlayerPosition(username)
	if !found {
		return fmt.Errorf("玩家不在游戏中")
	}

	if position != game.CurrentTurn {
		return fmt.Errorf("不是该玩家的回合")
	}

	player := game.GetPlayer(position)
	if player == nil {
		return fmt.Errorf("玩家不存在")
	}

	// 检查玩家是否有这些牌
	if !player.HasCards(cards) {
		return fmt.Errorf("玩家没有这些牌")
	}

	// 分析牌型
	handPattern := models.AnalyzeHand(cards)
	if !handPattern.IsValid {
		return fmt.Errorf("无效的牌型")
	}

	// 检查是否能压过上一手牌
	if len(game.LastPlayCards) > 0 && game.LastPlayer != position {
		lastPattern := models.AnalyzeHand(game.LastPlayCards)
		if !models.CanBeat(handPattern, lastPattern) {
			return fmt.Errorf("无法压过上一手牌")
		}
	}

	// 出牌
	if !player.RemoveCards(cards) {
		return fmt.Errorf("移除手牌失败")
	}

	// 更新游戏状态
	game.LastPlayCards = cards
	game.LastPlayer = position
	game.NextTurn()

	// 添加游戏日志
	game.AddLog("play_cards", position, cards, fmt.Sprintf("%s 出牌 %s", player.UserName, handPattern.Type))

	// 检查是否获胜
	if player.GetCardCount() == 0 {
		game.Status = models.GameStatusFinished
		game.Winner = position
		now := time.Now()
		game.FinishedAt = &now

		// 计算分数
		gs.calculateScore(game)

		game.AddLog("win", position, nil, fmt.Sprintf("%s 获胜", player.UserName))

		// 结束房间游戏
		room, _ := gs.roomService.GetRoom(roomID)
		room.EndGame()
	}

	// 更新房间状态
	room, _ := gs.roomService.GetRoom(roomID)
	return gs.roomService.UpdateRoom(room)
}

// PassTurn 过牌
func (gs *GameService) PassTurn(roomID, username string) error {
	game, err := gs.GetGameByRoom(roomID)
	if err != nil {
		return err
	}

	if game.Status != models.GameStatusPlaying {
		return fmt.Errorf("游戏状态不正确")
	}

	position, found := game.GetPlayerPosition(username)
	if !found {
		return fmt.Errorf("玩家不在游戏中")
	}

	if position != game.CurrentTurn {
		return fmt.Errorf("不是该玩家的回合")
	}

	// 如果上一次出牌的是自己，不能过牌
	if game.LastPlayer == position {
		return fmt.Errorf("上次出牌是你自己，不能过牌")
	}

	player := game.GetPlayer(position)
	game.NextTurn()

	// 添加游戏日志
	game.AddLog("pass", position, nil, fmt.Sprintf("%s 过牌", player.UserName))

	// 检查是否一圈都过了
	if gs.checkAllPass(game) {
		// 清除上一手牌，让最后出牌的玩家继续出牌
		game.LastPlayCards = nil
		game.CurrentTurn = game.LastPlayer
		game.AddLog("new_round", game.LastPlayer, nil, "新一轮开始")
	}

	// 更新房间状态
	room, _ := gs.roomService.GetRoom(roomID)
	return gs.roomService.UpdateRoom(room)
}

// checkAllPass 检查是否所有其他玩家都过牌了
func (gs *GameService) checkAllPass(game *models.Game) bool {
	// 这里简化处理，实际需要记录每轮的过牌状态
	// 如果当前回合回到了上次出牌的玩家，说明其他人都过了
	return game.CurrentTurn == game.LastPlayer
}

// calculateScore 计算分数
func (gs *GameService) calculateScore(game *models.Game) {
	if game.Status != models.GameStatusFinished {
		return
	}

	winner := game.GetPlayer(game.Winner)
	if winner == nil {
		return
	}

	baseScore := 1

	// 地主获胜
	if winner.Role == models.RoleLandlord {
		// 地主获胜，地主得2分，农民各扣1分
		winner.Score = baseScore * 2
		for _, player := range game.Players {
			if player != nil && player.Role == models.RoleFarmer {
				player.Score = -baseScore
			}
		}
	} else {
		// 农民获胜，每个农民得1分，地主扣2分
		for _, player := range game.Players {
			if player != nil {
				if player.Role == models.RoleFarmer {
					player.Score = baseScore
				} else if player.Role == models.RoleLandlord {
					player.Score = -baseScore * 2
				}
			}
		}
	}

	logger.Info("游戏 %s 结束，获胜者: %s", game.ID, winner.UserName)
}

// GetPlayerHand 获取玩家手牌（只能获取自己的手牌）
func (gs *GameService) GetPlayerHand(roomID, username string) ([]models.Card, error) {
	game, err := gs.GetGameByRoom(roomID)
	if err != nil {
		return nil, err
	}

	player := game.GetPlayerByName(username)
	if player == nil {
		return nil, fmt.Errorf("玩家不在游戏中")
	}

	return player.Cards, nil
}

// GetGameState 获取游戏状态（公开信息）
func (gs *GameService) GetGameState(roomID, username string) (map[string]interface{}, error) {
	game, err := gs.GetGameByRoom(roomID)
	if err != nil {
		return nil, err
	}

	// 构建游戏状态信息
	state := map[string]interface{}{
		"game_id":     game.ID,
		"status":      game.Status,
		"status_name": models.GameStatusNames[game.Status],
		"current_turn": game.CurrentTurn,
		"last_play_cards": game.LastPlayCards,
		"last_player": game.LastPlayer,
		"created_at":  game.CreatedAt,
		"started_at":  game.StartedAt,
		"finished_at": game.FinishedAt,
		"winner":      game.Winner,
	}

	// 玩家信息（隐藏其他玩家的手牌）
	players := make([]map[string]interface{}, 3)
	for i, player := range game.Players {
		if player == nil {
			players[i] = nil
			continue
		}

		playerInfo := map[string]interface{}{
			"username":      player.UserName,
			"position":      player.Position,
			"role":          player.Role,
			"role_name":     models.RoleNames[player.Role],
			"card_count":    player.GetCardCount(),
			"is_ready":      player.IsReady,
			"is_online":     player.IsOnline,
			"score":         player.Score,
			"call_landlord": player.CallLandlord,
		}

		// 只有自己能看到自己的手牌
		if player.UserName == username {
			playerInfo["cards"] = player.Cards
		}

		players[i] = playerInfo
	}
	state["players"] = players

	// 地主牌（只有地主确定后才显示）
	if game.Status >= models.GameStatusPlaying {
		state["landlord_cards"] = game.LandlordCards
	}

	return state, nil
}