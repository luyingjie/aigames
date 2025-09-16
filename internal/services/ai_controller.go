package services

import (
	"fmt"
	"time"

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
}

// NewAIController 创建AI控制器
func NewAIController(player *models.GamePlayer, gameService *GameService, roomID string) *AIController {
	return &AIController{
		player:      player,
		gameService: gameService,
		actionChan:  make(chan bool, 1),
		stopChan:    make(chan bool, 1),
		roomID:      roomID,
	}
}

// Start 启动AI控制器
func (c *AIController) Start() {
	logger.Info("AI玩家 %s 控制器启动", c.player.UserName)

	for {
		select {
		case <-c.actionChan:
			// 模拟思考时间
			time.Sleep(1 * time.Second)

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
		// 叫地主阶段：AI永远不叫地主
		return c.callLandlord(false)

	case models.GameStatusPlaying:
		// 出牌阶段：AI永远过牌
		return c.passTurn()

	default:
		logger.Info("AI玩家 %s 在状态 %d 下无需操作", c.player.UserName, game.Status)
		return nil
	}
}

// callLandlord AI叫地主操作
func (c *AIController) callLandlord(call bool) error {
	logger.Info("AI玩家 %s 叫地主: %t", c.player.UserName, call)

	// 调用GameService的叫地主方法
	return c.gameService.CallLandlord(c.roomID, c.player.UserName, call)
}

// passTurn AI过牌操作
func (c *AIController) passTurn() error {
	logger.Info("AI玩家 %s 过牌", c.player.UserName)

	// 调用GameService的过牌方法
	return c.gameService.PassTurn(c.roomID, c.player.UserName)
}