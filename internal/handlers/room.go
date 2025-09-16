package handlers

import (
	"fmt"
	"time"

	"aigames/internal/models"
	"aigames/internal/services"
	"aigames/pkg/logger"
	"aigames/pkg/protocol"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

// Room 房间处理器
type Room struct {
	component.Base
	roomService *services.RoomService
	gameService *services.GameService
}

// NewRoom 创建房间处理器实例
func NewRoom(roomService *services.RoomService, gameService *services.GameService) *Room {
	return &Room{
		roomService: roomService,
		gameService: gameService,
	}
}

// CreateRoom 创建房间
func (h *Room) CreateRoom(s *session.Session, req *protocol.CreateRoomRequest) error {
	logger.Info("创建房间请求: %s", req.Name)

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 获取用户名（假设从session中获取）
	username := s.String("username")
	if username == "" {
		resp := protocol.Unauthorized("请先登录")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 生成房间ID
	roomID := fmt.Sprintf("room_%d", time.Now().Unix())

	// 创建房间
	room, err := h.roomService.CreateRoom(roomID, req.Name, username, req.Type, req.Password, req.AICount)
	if err != nil {
		logger.Error("创建房间失败: %v", err)
		resp := protocol.InternalServerError("创建房间失败")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 房主自动加入房间
	_, err = h.roomService.JoinRoom(roomID, username, req.Password)
	if err != nil {
		logger.Error("房主加入房间失败: %v", err)
		// 删除刚创建的房间
		h.roomService.DeleteRoom(roomID)
		resp := protocol.InternalServerError("加入房间失败")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.CreateRoomSuccess(room)
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 创建房间成功: %s", username, roomID)
	return s.Response(resp)
}

// JoinRoom 加入房间
func (h *Room) JoinRoom(s *session.Session, req *protocol.JoinRoomRequest) error {
	logger.Info("加入房间请求: %s", req.RoomID)

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

	// 加入房间
	room, err := h.roomService.JoinRoom(req.RoomID, username, req.Password)
	if err != nil {
		logger.Error("加入房间失败: %v", err)

		// 根据错误类型返回不同响应
		var resp protocol.BaseResponse
		if err.Error() == "房间不存在" {
			resp = protocol.RoomNotFound()
		} else if err.Error() == "房间已满" {
			resp = protocol.RoomFull()
		} else if err.Error() == "房间密码错误" {
			resp = protocol.Unauthorized("房间密码错误")
		} else {
			resp = protocol.InternalServerError("加入房间失败")
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 保存房间ID到session
	s.Set("room_id", req.RoomID)

	resp := protocol.JoinRoomSuccess(room)
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 加入房间成功: %s", username, req.RoomID)
	return s.Response(resp)
}

// LeaveRoom 离开房间
func (h *Room) LeaveRoom(s *session.Session, req *protocol.LeaveRoomRequest) error {
	logger.Info("离开房间请求: %s", req.RoomID)

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

	// 离开房间
	err := h.roomService.LeaveRoom(req.RoomID, username)
	if err != nil {
		logger.Error("离开房间失败: %v", err)
		resp := protocol.InternalServerError("离开房间失败")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 清除session中的房间ID
	s.Remove("room_id")

	resp := protocol.LeaveRoomSuccess()
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 离开房间成功: %s", username, req.RoomID)
	return s.Response(resp)
}

// GetRoomList 获取房间列表
func (h *Room) GetRoomList(s *session.Session, req *protocol.GetRoomListRequest) error {
	logger.Info("获取房间列表请求")

	// 获取房间列表
	var rooms []*models.Room
	if req.Type == models.RoomTypePublic {
		rooms = h.roomService.GetPublicRooms()
	} else {
		rooms = h.roomService.GetAllRooms()
	}

	// 分页处理
	total := len(rooms)
	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start > total {
		rooms = []*models.Room{}
	} else if end > total {
		rooms = rooms[start:]
	} else {
		rooms = rooms[start:end]
	}

	resp := protocol.RoomListSuccess(rooms, total, req.Page, req.Size)
	resp.SetRequestId(req.RequestId)

	return s.Response(resp)
}

// DeleteRoom 删除房间
func (h *Room) DeleteRoom(s *session.Session, req *protocol.DeleteRoomRequest) error {
	logger.Info("删除房间请求: %s", req.RoomID)

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

	// 获取房间信息
	room, err := h.roomService.GetRoom(req.RoomID)
	if err != nil {
		resp := protocol.RoomNotFound()
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 检查是否是房主
	if room.Owner != username {
		resp := protocol.Forbidden("只有房主可以删除房间")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 删除房间
	err = h.roomService.DeleteRoom(req.RoomID)
	if err != nil {
		logger.Error("删除房间失败: %v", err)
		resp := protocol.InternalServerError("删除房间失败")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.DeleteRoomSuccess()
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 删除房间成功: %s", username, req.RoomID)
	return s.Response(resp)
}

// SetReady 设置准备状态
func (h *Room) SetReady(s *session.Session, req *protocol.SetReadyRequest) error {
	logger.Info("设置准备状态请求: %s, ready=%t", req.RoomID, req.Ready)

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

	// 设置准备状态
	err := h.roomService.SetPlayerReady(req.RoomID, username, req.Ready)
	if err != nil {
		logger.Error("设置准备状态失败: %v", err)

		var resp protocol.BaseResponse
		if err.Error() == "房间不存在" {
			resp = protocol.RoomNotFound()
		} else if err.Error() == "玩家不在游戏中" {
			resp = protocol.PlayerNotInRoom()
		} else {
			resp = protocol.InternalServerError("设置准备状态失败")
		}

		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	resp := protocol.SetReadySuccess()
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 设置准备状态成功: %t", username, req.Ready)
	return s.Response(resp)
}

// StartGame 开始游戏
func (h *Room) StartGame(s *session.Session, req *protocol.StartGameRequest) error {
	logger.Info("开始游戏请求: %s", req.RoomID)

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

	// 检查是否为房主
	room, err := h.roomService.GetRoom(req.RoomID)
	if err != nil {
		resp := protocol.RoomNotFound()
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	if room.Owner != username {
		resp := protocol.Forbidden("只有房主可以开始游戏")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 开始游戏
	_, err = h.roomService.StartGame(req.RoomID)
	if err != nil {
		logger.Error("开始游戏失败: %v", err)
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 启动AI控制器
	if err := h.gameService.StartAIControllers(req.RoomID); err != nil {
		logger.Error("启动AI控制器失败: %v", err)
	}

	// 检查是否第一个玩家是AI，如果是则通知其行动
	game, err := h.gameService.GetGameByRoom(req.RoomID)
	if err == nil {
		currentPlayer := game.GetPlayer(game.CurrentTurn)
		if currentPlayer != nil && currentPlayer.IsAI {
			h.gameService.NotifyAITurn(req.RoomID, currentPlayer.UserName)
		}
	}

	resp := protocol.StartGameSuccess()
	resp.SetRequestId(req.RequestId)

	logger.Info("房间 %s 游戏开始", req.RoomID)
	return s.Response(resp)
}
