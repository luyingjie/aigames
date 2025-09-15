package handlers

import (
	"time"

	"aigames/internal/models"
	"aigames/internal/services"
	"aigames/pkg/logger"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
)

type (
	// Handler 处理器结构体
	User struct {
		component.Base
		userService *services.UserService
	}
	// LoginRequest 登录请求结构
	LoginRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// LoginResponse 登录响应结构
	LoginResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// SignupRequest 注册请求结构
	SignupRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Age      int    `json:"age"`
	}

	// SignupResponse 注册响应结构
	SignupResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Name string `json:"name"`
	}
)

func NewUser(userService *services.UserService) *User {
	return &User{userService: userService}
}

// Login 登录处理方法
func (h *User) Login(s *session.Session, req *LoginRequest) error {
	logger.Info("用户登录请求: %s", req.Name)

	// 验证输入
	if req.Name == "" {
		resp := &LoginResponse{
			Code: 400,
			Msg:  "用户名不能为空",
		}
		return s.Response(resp)
	}

	if req.Password == "" {
		resp := &LoginResponse{
			Code: 400,
			Msg:  "密码不能为空",
		}
		return s.Response(resp)
	}

	// 从数据库获取用户
	user, err := h.userService.GetUser(req.Name)
	if err != nil {
		resp := &LoginResponse{
			Code: 404,
			Msg:  "用户不存在",
		}
		return s.Response(resp)
	}

	// 验证密码
	hashedPassword := h.userService.HashPassword(req.Password)
	if user.Password != hashedPassword {
		resp := &LoginResponse{
			Code: 401,
			Msg:  "密码错误",
		}
		return s.Response(resp)
	}

	// 更新最后登录时间
	if err := h.userService.UpdateLastLogin(req.Name); err != nil {
		logger.Error("更新登录时间失败: %v", err)
	}

	// 登录成功
	resp := &LoginResponse{
		Code: 200,
		Msg:  "登录成功",
		Name: user.Name,
		Age:  user.Age,
	}

	logger.Info("用户 %s 登录成功", req.Name)
	return s.Response(resp)
}

// Signup 注册处理方法
func (h *User) Signup(s *session.Session, req *SignupRequest) error {
	logger.Info("用户注册请求: %s", req.Name)

	// 验证输入
	if req.Name == "" {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "用户名不能为空",
		}
		return s.Response(resp)
	}

	if req.Password == "" {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "密码不能为空",
		}
		return s.Response(resp)
	}

	if req.Age <= 0 || req.Age > 150 {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "年龄必须在1-150之间",
		}
		return s.Response(resp)
	}

	// 检查用户是否已存在
	if h.userService.UserExists(req.Name) {
		resp := &SignupResponse{
			Code: 409,
			Msg:  "用户名已存在",
		}
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
		resp := &SignupResponse{
			Code: 500,
			Msg:  "注册失败，请稍后重试",
		}
		return s.Response(resp)
	}

	// 注册成功
	resp := &SignupResponse{
		Code: 200,
		Msg:  "注册成功",
		Name: user.Name,
	}

	logger.Info("用户 %s 注册成功", req.Name)
	return s.Response(resp)
}
