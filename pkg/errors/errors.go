package errors

import (
	"aigames/pkg/constants"
	"fmt"
)

// 自定义错误类型
type AppError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误消息
	Details string `json:"details"` // 详细信息
}

// 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// 创建新的应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// 创建带详细信息的应用错误
func NewWithDetails(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// 预定义的常见错误
var (
	// 通用错误
	ErrInternalServer = &AppError{
		Code:    constants.StatusServerError,
		Message: "服务器内部错误",
	}

	ErrBadRequest = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "请求参数错误",
	}

	ErrUnauthorized = &AppError{
		Code:    constants.StatusUnauthorized,
		Message: "未授权访问",
	}

	ErrNotFound = &AppError{
		Code:    constants.StatusNotFound,
		Message: "资源未找到",
	}

	// 用户相关错误
	ErrUserNotFound = &AppError{
		Code:    constants.StatusNotFound,
		Message: "用户不存在",
	}

	ErrUserExists = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "用户已存在",
	}

	ErrInvalidCredentials = &AppError{
		Code:    constants.StatusUnauthorized,
		Message: "用户名或密码错误",
	}

	ErrInvalidToken = &AppError{
		Code:    constants.StatusUnauthorized,
		Message: "无效的访问令牌",
	}

	// 房间相关错误
	ErrRoomNotFound = &AppError{
		Code:    constants.StatusNotFound,
		Message: "房间不存在",
	}

	ErrRoomFull = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "房间已满",
	}

	ErrNotInRoom = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "玩家不在房间中",
	}

	ErrAlreadyInRoom = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "玩家已在房间中",
	}

	// 游戏相关错误
	ErrGameNotFound = &AppError{
		Code:    constants.StatusNotFound,
		Message: "游戏不存在",
	}

	ErrGameNotStarted = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "游戏尚未开始",
	}

	ErrGameAlreadyStarted = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "游戏已经开始",
	}

	ErrNotYourTurn = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "当前不是您的回合",
	}

	ErrInvalidMove = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "无效的出牌",
	}

	ErrInvalidCards = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "无效的牌组",
	}

	// AI相关错误
	ErrAIPlayerNotFound = &AppError{
		Code:    constants.StatusNotFound,
		Message: "AI玩家不存在",
	}

	ErrAIConfigInvalid = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "AI配置无效",
	}

	ErrAIAPIError = &AppError{
		Code:    constants.StatusServerError,
		Message: "AI API调用失败",
	}

	// 数据库相关错误
	ErrDatabaseConnection = &AppError{
		Code:    constants.StatusServerError,
		Message: "数据库连接失败",
	}

	ErrDatabaseOperation = &AppError{
		Code:    constants.StatusServerError,
		Message: "数据库操作失败",
	}

	// WebSocket相关错误
	ErrWebSocketConnection = &AppError{
		Code:    constants.StatusServerError,
		Message: "WebSocket连接失败",
	}

	ErrInvalidMessage = &AppError{
		Code:    constants.StatusBadRequest,
		Message: "无效的消息格式",
	}
)

// 判断是否为指定的错误类型
func Is(err error, target *AppError) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == target.Code && appErr.Message == target.Message
	}
	return false
}

// 从error转换为AppError
func FromError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return NewWithDetails(constants.StatusServerError, "未知错误", err.Error())
}
