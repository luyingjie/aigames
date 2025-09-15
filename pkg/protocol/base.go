package protocol

import (
	"time"

	"github.com/google/uuid"
)

// BaseRequest 基础请求结构
type BaseRequest struct {
	RequestId string    `json:"request_id"`          // 请求追踪ID
	Timestamp time.Time `json:"timestamp"`           // 请求时间戳
	ClientId  string    `json:"client_id,omitempty"` // 客户端ID（可选）
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code      int         `json:"code"`                // 业务状态码
	Message   string      `json:"message"`             // 状态消息
	Data      interface{} `json:"data,omitempty"`      // 业务数据
	Timestamp time.Time   `json:"timestamp"`           // 响应时间戳
	RequestId string      `json:"request_id"`          // 对应的请求ID
}

// PageRequest 分页请求结构
type PageRequest struct {
	BaseRequest
	Page int `json:"page"` // 页码，从1开始
	Size int `json:"size"` // 每页大小，默认10
}

// PageResponse 分页响应结构
type PageResponse struct {
	BaseResponse
	Total int `json:"total"` // 总记录数
	Page  int `json:"page"`  // 当前页码
	Size  int `json:"size"`  // 每页大小
}

// NewBaseRequest 创建基础请求
func NewBaseRequest() BaseRequest {
	return BaseRequest{
		RequestId: uuid.New().String(),
		Timestamp: time.Now(),
	}
}

// NewPageRequest 创建分页请求
func NewPageRequest(page, size int) PageRequest {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return PageRequest{
		BaseRequest: NewBaseRequest(),
		Page:        page,
		Size:        size,
	}
}

// Success 创建成功响应
func Success(data interface{}) BaseResponse {
	return BaseResponse{
		Code:      StatusOK,
		Message:   "操作成功",
		Data:      data,
		Timestamp: time.Now(),
	}
}

// SuccessWithMessage 创建带自定义消息的成功响应
func SuccessWithMessage(data interface{}, message string) BaseResponse {
	return BaseResponse{
		Code:      StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Error 创建错误响应
func Error(code int, message string) BaseResponse {
	return BaseResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ErrorWithData 创建带数据的错误响应
func ErrorWithData(code int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// SuccessPage 创建分页成功响应
func SuccessPage(data interface{}, total, page, size int) PageResponse {
	return PageResponse{
		BaseResponse: BaseResponse{
			Code:      StatusOK,
			Message:   "查询成功",
			Data:      data,
			Timestamp: time.Now(),
		},
		Total: total,
		Page:  page,
		Size:  size,
	}
}

// ErrorPage 创建分页错误响应
func ErrorPage(code int, message string, page, size int) PageResponse {
	return PageResponse{
		BaseResponse: BaseResponse{
			Code:      code,
			Message:   message,
			Timestamp: time.Now(),
		},
		Total: 0,
		Page:  page,
		Size:  size,
	}
}

// SetRequestId 设置响应的请求ID
func (r *BaseResponse) SetRequestId(requestId string) {
	r.RequestId = requestId
}

// SetRequestId 设置分页响应的请求ID
func (r *PageResponse) SetRequestId(requestId string) {
	r.BaseResponse.RequestId = requestId
}

// IsSuccess 判断响应是否成功
func (r *BaseResponse) IsSuccess() bool {
	return r.Code == StatusOK
}

// IsSuccess 判断分页响应是否成功
func (r *PageResponse) IsSuccess() bool {
	return r.BaseResponse.Code == StatusOK
}