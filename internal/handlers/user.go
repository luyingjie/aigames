package handlers

import (
	"time"

	"aigames/internal/models"
	"aigames/internal/services"
	"aigames/pkg/logger"
	"aigames/pkg/protocol"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

type (
	// Handler 处理器结构体
	User struct {
		component.Base
		userService *services.UserService
	}
)

func NewUser(userService *services.UserService) *User {
	return &User{userService: userService}
}

// Login 登录处理方法
func (h *User) Login(s *session.Session, req *protocol.LoginRequest) error {
	logger.Info("用户登录请求: %s", req.Name)

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 从数据库获取用户
	user, err := h.userService.GetUser(req.Name)
	if err != nil {
		resp := protocol.UserNotFound()
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 验证密码
	hashedPassword := h.userService.HashPassword(req.Password)
	if user.Password != hashedPassword {
		resp := protocol.PasswordIncorrect()
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 更新最后登录时间
	if err := h.userService.UpdateLastLogin(req.Name); err != nil {
		logger.Error("更新登录时间失败: %v", err)
	}

	// 登录成功，保存用户信息到session
	s.Set("username", user.Name)

	// 登录成功
	resp := protocol.LoginSuccess(user.Name, user.Age)
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 登录成功", req.Name)
	return s.Response(resp)
}

// Signup 注册处理方法
func (h *User) Signup(s *session.Session, req *protocol.SignupRequest) error {
	logger.Info("用户注册请求: %s", req.Name)

	// 验证请求参数
	if err := protocol.ValidateRequest(req); err != nil {
		resp := protocol.BadRequest(err.Error())
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 检查用户是否已存在
	if h.userService.UserExists(req.Name) {
		resp := protocol.UserExists()
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 创建新用户
	user := &models.User{
		Name:        req.Name,
		Password:    h.userService.HashPassword(req.Password),
		Age:         req.Age,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}

	// 保存用户到数据库
	if err := h.userService.SaveUser(user); err != nil {
		logger.Error("保存用户失败: %v", err)
		resp := protocol.InternalServerError("注册失败，请稍后重试")
		resp.SetRequestId(req.RequestId)
		return s.Response(resp)
	}

	// 注册成功
	resp := protocol.SignupSuccess(user.Name)
	resp.SetRequestId(req.RequestId)

	logger.Info("用户 %s 注册成功", req.Name)
	return s.Response(resp)
}
