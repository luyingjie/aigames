package handlers

import (
	"aigames/internal/models"
	"aigames/internal/services"
	"aigames/pkg/logger"
	"aigames/pkg/protocol"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

// Game 游戏处理器
type Game struct {
	component.Base
	gameService *services.GameService
	roomService *services.RoomService
}

// NewGame 创建游戏处理器实例
func NewGame(gameService *services.GameService, roomService *services.RoomService) *Game {
	return &Game{
		gameService: gameService,
		roomService: roomService,
	}
}

// CallLandlord 叫地主
func (h *Game) CallLandlord(s *session.Session, req *protocol.CallLandlordRequest) error {
	logger.Info("叫地主请求: %s, call=%t", req.RoomID, req.Call)

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 叫地主
	err := h.gameService.CallLandlord(req.RoomID, username, req.Call)
	if err != nil {
		logger.Error("叫地主失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间没有活跃的游戏" {
			resp = protocol.GameNotFound()
		} else if err.Error() == "玩家不在游戏中" {
			resp = protocol.PlayerNotInRoom()
		} else if err.Error() == "不是该玩家的回合" {
			resp = protocol.NotPlayerTurn()
		} else if err.Error() == "当前不是叫地主阶段" {
			resp = protocol.InvalidMove()
		} else {
			resp = protocol.InternalServerError("操作失败")
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.CallLandlordSuccess()
	resp.SetRequestId(req.RequestId)

	action := "不叫"
	if req.Call {
		action = "叫地主"
	}
	logger.Info("用户 %s %s", username, action)
	return s.Response(resp)
}

// PlayCards 出牌
func (h *Game) PlayCards(s *session.Session, req *protocol.PlayCardsRequest) error {
	logger.Info("出牌请求: %s, 牌数=%d", req.RoomID, len(req.Cards))

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	if len(req.Cards) == 0 {
		resp := protocol.BadRequest("请选择要出的牌")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 出牌
	err := h.gameService.PlayCards(req.RoomID, username, req.Cards)
	if err != nil {
		logger.Error("出牌失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间没有活跃的游戏" {
			resp = protocol.GameNotFound()
		} else if err.Error() == "玩家不在游戏中" {
			resp = protocol.PlayerNotInRoom()
		} else if err.Error() == "不是该玩家的回合" {
			resp = protocol.NotPlayerTurn()
		} else if err.Error() == "游戏状态不正确" {
			resp = protocol.GameNotStarted()
		} else {
			resp = protocol.BadRequest(err.Error())
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.PlayCardsSuccess()
	resp.SetRequestId(req.RequestId)

	// 分析牌型用于日志
	handPattern := models.AnalyzeHand(req.Cards)
	logger.Info("用户 %s 出牌成功: %s", username, models.HandTypeNames[handPattern.Type])
	return s.Response(resp)
}

// PassTurn 过牌
func (h *Game) PassTurn(s *session.Session, req *protocol.PassTurnRequest) error {
	logger.Info("过牌请求: %s", req.RoomID)

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 过牌
	err := h.gameService.PassTurn(req.RoomID, username)
	if err != nil {
		logger.Error("过牌失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间没有活跃的游戏" {
			resp = protocol.GameNotFound()
		} else if err.Error() == "玩家不在游戏中" {
			resp = protocol.PlayerNotInRoom()
		} else if err.Error() == "不是该玩家的回合" {
			resp = protocol.NotPlayerTurn()
		} else if err.Error() == "游戏状态不正确" {
			resp = protocol.GameNotStarted()
		} else {
			resp = protocol.BadRequest(err.Error())
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.PassTurnSuccess()
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 过牌", username)
	return s.Response(resp)
}

// GetGameState 获取游戏状态
func (h *Game) GetGameState(s *session.Session, req *protocol.GetGameStateRequest) error {
	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取游戏状态
	gameState, err := h.gameService.GetGameState(req.RoomID, username)
	if err != nil {
		logger.Error("获取游戏状态失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间没有活跃的游戏" {
			resp = protocol.GameNotFound()
		} else if err.Error() == "房间不存在" {
			resp = protocol.RoomNotFound()
		} else {
			resp = protocol.InternalServerError("获取游戏状态失败")
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.GameStateSuccess(gameState)
	resp.SetRequestId(req.RequestId)

	return s.Response(resp)
}

// GetPlayerHand 获取玩家手牌
func (h *Game) GetPlayerHand(s *session.Session, req *protocol.GetPlayerHandRequest) error {
	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取玩家手牌
	cards, err := h.gameService.GetPlayerHand(req.RoomID, username)
	if err != nil {
		logger.Error("获取玩家手牌失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间没有活跃的游戏" {
			resp = protocol.GameNotFound()
		} else if err.Error() == "房间不存在" {
			resp = protocol.RoomNotFound()
		} else if err.Error() == "玩家不在游戏中" {
			resp = protocol.PlayerNotInRoom()
		} else {
			resp = protocol.InternalServerError("获取手牌失败")
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.PlayerHandSuccess(cards)
	resp.SetRequestId(req.RequestId)

	return s.Response(resp)
}