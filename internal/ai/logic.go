package ai

import (
	"aigames/pkg/logger"
)

// PlayerInterface 定义玩家接口
type PlayerInterface interface {
	GetUserName() string
}

// PlayerWrapper 包装器结构
type PlayerWrapper struct {
	UserName string
}

// GetUserName 获取玩家用户名
func (p *PlayerWrapper) GetUserName() string {
	return p.UserName
}

// CallLandlord AI叫地主操作
func CallLandlord(player *PlayerWrapper, gameService interface {
	CallLandlord(roomID, username string, call bool) error
}, roomID string, call bool) error {
	logger.Info("AI玩家 %s 叫地主: %t", player.GetUserName(), call)
	return gameService.CallLandlord(roomID, player.GetUserName(), call)
}

// PassTurn AI过牌操作
func PassTurn(player *PlayerWrapper, gameService interface {
	PassTurn(roomID, username string) error
}, roomID string) error {
	logger.Info("AI玩家 %s 过牌", player.GetUserName())
	return gameService.PassTurn(roomID, player.GetUserName())
}
