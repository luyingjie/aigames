package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lonng/nano/session"
)

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`           // 状态码
	Message string      `json:"message"`        // 消息
	Data    interface{} `json:"data,omitempty"` // 数据
	Time    int64       `json:"time"`           // 时间戳
}

// Success 成功响应
func Success(s *session.Session, data interface{}) {
	s.Response(Response{
		Code:    200,
		Message: "success",
		Data:    data,
		Time:    time.Now().Unix(),
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
		Time:    time.Now().Unix(),
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: err.Error(),
		Time:    time.Now().Unix(),
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Time:    time.Now().Unix(),
	})
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
		Time:    time.Now().Unix(),
	})
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
		Time:    time.Now().Unix(),
	})
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
		Time:    time.Now().Unix(),
	})
}

// GenerateID 生成随机ID
func GenerateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateMessageID 生成消息ID
func GenerateMessageID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("msg_%d_%s", timestamp, hex.EncodeToString(randomBytes))
}

// ParseInt 解析整数参数
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultValue
}

// ParseBool 解析布尔参数
func ParseBool(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}
	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}
	return defaultValue
}

// Min 获取最小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max 获取最大值
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Contains 检查切片是否包含元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveFromSlice 从切片中移除元素
func RemoveFromSlice(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// SanitizeString 清理字符串中的危险字符
func SanitizeString(s string) string {
	// 移除危险字符
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
