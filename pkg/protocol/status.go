package protocol

// 通用状态码定义
const (
	// 成功状态码
	StatusOK = 200 // 操作成功

	// 客户端错误状态码 (4xx)
	StatusBadRequest          = 400 // 请求参数错误
	StatusUnauthorized        = 401 // 未授权
	StatusForbidden           = 403 // 禁止访问
	StatusNotFound            = 404 // 资源不存在
	StatusMethodNotAllowed    = 405 // 方法不允许
	StatusConflict            = 409 // 资源冲突
	StatusUnprocessableEntity = 422 // 参数验证失败
	StatusTooManyRequests     = 429 // 请求过于频繁

	// 服务器错误状态码 (5xx)
	StatusInternalServerError = 500 // 服务器内部错误
	StatusBadGateway          = 502 // 网关错误
	StatusServiceUnavailable  = 503 // 服务不可用
	StatusGatewayTimeout      = 504 // 网关超时

	// 业务错误状态码 (1xxx)
	StatusUserNotFound      = 1001 // 用户不存在
	StatusUserExists        = 1002 // 用户已存在
	StatusPasswordIncorrect = 1003 // 密码错误
	StatusUserLocked        = 1004 // 用户被锁定
	StatusTokenExpired      = 1005 // token过期
	StatusTokenInvalid      = 1006 // token无效

	// 数据库错误状态码 (2xxx)
	StatusDatabaseError    = 2001 // 数据库错误
	StatusDatabaseTimeout  = 2002 // 数据库超时
	StatusDatabaseConnFail = 2003 // 数据库连接失败
	StatusDataNotFound     = 2004 // 数据不存在
	StatusDataExists       = 2005 // 数据已存在
	StatusDataConstraint   = 2006 // 数据约束冲突

	// 第三方服务错误状态码 (3xxx)
	StatusThirdPartyError   = 3001 // 第三方服务错误
	StatusThirdPartyTimeout = 3002 // 第三方服务超时
	StatusThirdPartyLimit   = 3003 // 第三方服务限流

	// 游戏业务错误状态码 (4xxx)
	StatusGameNotFound    = 4001 // 游戏不存在
	StatusGameNotStarted  = 4002 // 游戏未开始
	StatusGameEnded       = 4003 // 游戏已结束
	StatusRoomFull        = 4004 // 房间已满
	StatusRoomNotFound    = 4005 // 房间不存在
	StatusPlayerNotInRoom = 4006 // 玩家不在房间内
	StatusNotPlayerTurn   = 4007 // 不是玩家回合
	StatusInvalidMove     = 4008 // 无效操作
)

// StatusMessages 状态码对应的消息
var StatusMessages = map[int]string{
	// 成功状态
	StatusOK: "操作成功",

	// 客户端错误
	StatusBadRequest:          "请求参数错误",
	StatusUnauthorized:        "未授权访问",
	StatusForbidden:           "禁止访问",
	StatusNotFound:            "资源不存在",
	StatusMethodNotAllowed:    "方法不允许",
	StatusConflict:            "资源冲突",
	StatusUnprocessableEntity: "参数验证失败",
	StatusTooManyRequests:     "请求过于频繁",

	// 服务器错误
	StatusInternalServerError: "服务器内部错误",
	StatusBadGateway:          "网关错误",
	StatusServiceUnavailable:  "服务不可用",
	StatusGatewayTimeout:      "网关超时",

	// 业务错误
	StatusUserNotFound:      "用户不存在",
	StatusUserExists:        "用户已存在",
	StatusPasswordIncorrect: "密码错误",
	StatusUserLocked:        "用户被锁定",
	StatusTokenExpired:      "登录已过期",
	StatusTokenInvalid:      "登录凭证无效",

	// 数据库错误
	StatusDatabaseError:    "数据库错误",
	StatusDatabaseTimeout:  "数据库操作超时",
	StatusDatabaseConnFail: "数据库连接失败",
	StatusDataNotFound:     "数据不存在",
	StatusDataExists:       "数据已存在",
	StatusDataConstraint:   "数据约束冲突",

	// 第三方服务错误
	StatusThirdPartyError:   "第三方服务错误",
	StatusThirdPartyTimeout: "第三方服务超时",
	StatusThirdPartyLimit:   "第三方服务限流",

	// 游戏业务错误
	StatusGameNotFound:    "游戏不存在",
	StatusGameNotStarted:  "游戏未开始",
	StatusGameEnded:       "游戏已结束",
	StatusRoomFull:        "房间已满",
	StatusRoomNotFound:    "房间不存在",
	StatusPlayerNotInRoom: "玩家不在房间内",
	StatusNotPlayerTurn:   "不是你的回合",
	StatusInvalidMove:     "无效操作",
}

// GetStatusMessage 获取状态码对应的消息
func GetStatusMessage(code int) string {
	if msg, exists := StatusMessages[code]; exists {
		return msg
	}
	return "未知错误"
}

// ErrorWithCode 根据状态码创建错误响应
func ErrorWithCode(code int) BaseResponse {
	return Error(code, GetStatusMessage(code))
}

// 快捷错误响应方法
func BadRequest(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusBadRequest)
	}
	return Error(StatusBadRequest, message)
}

func Unauthorized(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusUnauthorized)
	}
	return Error(StatusUnauthorized, message)
}

func Forbidden(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusForbidden)
	}
	return Error(StatusForbidden, message)
}

func NotFound(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusNotFound)
	}
	return Error(StatusNotFound, message)
}

func Conflict(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusConflict)
	}
	return Error(StatusConflict, message)
}

func InternalServerError(message string) BaseResponse {
	if message == "" {
		message = GetStatusMessage(StatusInternalServerError)
	}
	return Error(StatusInternalServerError, message)
}

// 业务错误快捷方法
func UserNotFound() BaseResponse {
	return ErrorWithCode(StatusUserNotFound)
}

func UserExists() BaseResponse {
	return ErrorWithCode(StatusUserExists)
}

func PasswordIncorrect() BaseResponse {
	return ErrorWithCode(StatusPasswordIncorrect)
}

func UserLocked() BaseResponse {
	return ErrorWithCode(StatusUserLocked)
}

func TokenExpired() BaseResponse {
	return ErrorWithCode(StatusTokenExpired)
}

func TokenInvalid() BaseResponse {
	return ErrorWithCode(StatusTokenInvalid)
}
