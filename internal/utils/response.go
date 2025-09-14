package utils

import (
	"net/http"

	"ai-game/pkg/constants"
	appErrors "ai-game/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
	Page    PageInfo    `json:"page"`    // 分页信息
}

// PageInfo 分页信息
type PageInfo struct {
	Current int `json:"current"` // 当前页码
	Size    int `json:"size"`    // 每页数量
	Total   int `json:"total"`   // 总数量
	Pages   int `json:"pages"`   // 总页数
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    constants.StatusSuccess,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    constants.StatusSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	var appErr *appErrors.AppError

	// 判断是否为自定义错误
	if e, ok := err.(*appErrors.AppError); ok {
		appErr = e
	} else {
		// 转换为自定义错误
		appErr = appErrors.FromError(err)
	}

	// 根据错误码设置HTTP状态码
	var httpStatus int
	switch appErr.Code {
	case constants.StatusBadRequest:
		httpStatus = http.StatusBadRequest
	case constants.StatusUnauthorized:
		httpStatus = http.StatusUnauthorized
	case constants.StatusNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Data:    nil,
	})
}

// ErrorWithCode 指定错误码的错误响应
func ErrorWithCode(c *gin.Context, code int, message string) {
	var httpStatus int
	switch code {
	case constants.StatusBadRequest:
		httpStatus = http.StatusBadRequest
	case constants.StatusUnauthorized:
		httpStatus = http.StatusUnauthorized
	case constants.StatusNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    constants.StatusBadRequest,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    constants.StatusUnauthorized,
		Message: message,
		Data:    nil,
	})
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    constants.StatusNotFound,
		Message: message,
		Data:    nil,
	})
}

// InternalError 500错误响应
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    constants.StatusServerError,
		Message: message,
		Data:    nil,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, data interface{}, page PageInfo) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    constants.StatusSuccess,
		Message: "操作成功",
		Data:    data,
		Page:    page,
	})
}

// PageSuccessWithMessage 带自定义消息的分页成功响应
func PageSuccessWithMessage(c *gin.Context, message string, data interface{}, page PageInfo) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    constants.StatusSuccess,
		Message: message,
		Data:    data,
		Page:    page,
	})
}

// CalculatePages 计算总页数
func CalculatePages(total, size int) int {
	if size <= 0 {
		return 0
	}
	return (total + size - 1) / size
}

// NewPageInfo 创建分页信息
func NewPageInfo(current, size, total int) PageInfo {
	return PageInfo{
		Current: current,
		Size:    size,
		Total:   total,
		Pages:   CalculatePages(total, size),
	}
}

// ValidatePagination 验证分页参数
func ValidatePagination(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

// GetPaginationFromQuery 从查询参数获取分页信息
func GetPaginationFromQuery(c *gin.Context) (int, int) {
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")

	pageInt := 1
	sizeInt := 10

	if p, err := parseInt(page); err == nil {
		pageInt = p
	}

	if s, err := parseInt(size); err == nil {
		sizeInt = s
	}

	return ValidatePagination(pageInt, sizeInt)
}

// parseInt 字符串转整数的辅助函数
func parseInt(s string) (int, error) {
	var result int
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, appErrors.New(constants.StatusBadRequest, "无效的数字格式")
		}
		result = result*10 + int(char-'0')
	}
	return result, nil
}
