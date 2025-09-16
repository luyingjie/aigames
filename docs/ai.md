
```go
package services

import (
	"aigames/internal/models"
	"aigames/pkg/logger"
	"time"
)

// AIController 控制一个AI玩家的行为
type AIController struct {
	player      *models.GamePlayer
	gameService *GameService // 用来向游戏主逻辑发送操作
	actionChan  chan bool    // 用来接收游戏主逻辑的通知
	stopChan    chan bool    // 用来停止goroutine
}

// NewAIController 创建AI控制器实例
func NewAIController(player *models.GamePlayer, gameService *GameService) *AIController {
	return &AIController{
		player:      player,
		gameService: gameService,
		actionChan:  make(chan bool),
		stopChan:    make(chan bool),
	}
}

// Start 启动AI的监听和行动循环
func (c *AIController) Start() {
	logger.Info("启动AI控制器: %s", c.player.UserName)
	go func() {
		for {
			select {
			case <-c.actionChan:
				c.performAction()
			case <-c.stopChan:
				logger.Info("停止AI控制器: %s", c.player.UserName)
				return
			}
		}
	}()
}

// Stop 停止AI控制器
func (c *AIController) Stop() {
	close(c.stopChan)
}

// NotifyTurn 由 GameService 调用，通知AI行动
func (c *AIController) NotifyTurn() {
	// 使用非阻塞发送，以防channel阻塞
	select {
	case c.actionChan <- true:
	default:
	}
}

// performAction 执行AI的行动逻辑
func (c *AIController) performAction() {
	// 模拟AI思考
	time.Sleep(1 * time.Second)

	room, err := c.gameService.roomService.GetRoom(c.gameService.roomID)
	if err != nil {
		logger.Error("AI %s 获取房间失败: %v", c.player.UserName, err)
		return
	}

	game := room.CurrentGame
	if game == nil {
		logger.Error("AI %s 无法获取当前游戏", c.player.UserName)
		return
	}

	// 检查当前是否是自己的回合
	if game.CurrentTurn != c.player.Position {
		return // 不是自己的回合，忽略
	}

	switch game.Status {
	case models.GameStatusCalling:
		// 永远不叫地主
		logger.Info("AI %s 选择不叫地主", c.player.UserName)
		c.gameService.CallLandlord(c.player.UserName, false)
	case models.GameStatusPlaying:
		// 永远过牌
		logger.Info("AI %s 选择过牌", c.player.UserName)
		c.gameService.PassTurn(c.player.UserName)
	}
}
