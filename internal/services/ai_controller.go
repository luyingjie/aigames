package services

import (
	"fmt"
	"time"

	"aigames/internal/config"
	"aigames/internal/models"
	"aigames/pkg/logger"
)

// AIController AI控制器
type AIController struct {
	player      *models.GamePlayer
	gameService *GameService
	actionChan  chan bool
	stopChan    chan bool
	roomID      string
	aiService   *AIService
}

// NewAIController 创建AI控制器
func NewAIController(player *models.GamePlayer, gameService *GameService, roomID string) *AIController {
	// 创建AI服务
	aiService := NewAIService()

	return &AIController{
		player:      player,
		gameService: gameService,
		actionChan:  make(chan bool, 1),
		stopChan:    make(chan bool, 1),
		roomID:      roomID,
		aiService:   aiService,
	}
}

// Start 启动AI控制器
func (c *AIController) Start() {
	logger.Info("AI玩家 %s 控制器启动", c.player.UserName)

	for {
		select {
		case <-c.actionChan:
			// 模拟思考时间
			time.Sleep(time.Duration(config.GetConfig().AI.DefaultThinkTime) * time.Second)

			// 根据游戏状态执行不同的操作
			if err := c.executeAction(); err != nil {
				logger.Error("AI玩家 %s 执行操作失败: %v", c.player.UserName, err)
			}

		case <-c.stopChan:
			logger.Info("AI玩家 %s 控制器停止", c.player.UserName)
			return
		}
	}
}

// Stop 停止AI控制器
func (c *AIController) Stop() {
	select {
	case c.stopChan <- true:
	default:
	}
}

// NotifyTurn 通知轮到AI行动
func (c *AIController) NotifyTurn() {
	select {
	case c.actionChan <- true:
	default:
	}
}

// executeAction 执行AI操作
func (c *AIController) executeAction() error {
	// 获取当前游戏对象
	game, err := c.gameService.GetGameByRoom(c.roomID)
	if err != nil {
		return fmt.Errorf("获取游戏对象失败: %w", err)
	}

	switch game.Status {
	case models.GameStatusCalling:
		// 叫地主阶段：使用AI服务决策
		return c.aiService.CallLandlord(c.player, game, c.gameService, c.roomID)

	case models.GameStatusPlaying:
		// 出牌阶段：使用AI服务决策
		return c.aiService.PlayCards(c.player, game, c.gameService, c.roomID)

	default:
		logger.Info("AI玩家 %s 在状态 %d 下无需操作", c.player.UserName, game.Status)
		return nil
	}
}

// GetPlayer 获取AI玩家信息
func (c *AIController) GetPlayer() *models.GamePlayer {
	return c.player
}

// GetGameService 获取游戏服务
func (c *AIController) GetGameService() *GameService {
	return c.gameService
}

// GetRoomID 获取房间ID
func (c *AIController) GetRoomID() string {
	return c.roomID
}
