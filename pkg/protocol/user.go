package protocol

// 用户相关的请求和响应结构体

// LoginRequest 登录请求
type LoginRequest struct {
	BaseRequest
	Name     string `json:"name" validate:"required,min=1,max=50"`     // 用户名
	Password string `json:"password" validate:"required,min=6,max=50"` // 密码
}

// LoginResponse 登录响应数据
type LoginData struct {
	Name string `json:"name"` // 用户名
	Age  int    `json:"age"`  // 年龄
}

// SignupRequest 注册请求
type SignupRequest struct {
	BaseRequest
	Name     string `json:"name" validate:"required,min=1,max=50"`     // 用户名
	Password string `json:"password" validate:"required,min=6,max=50"` // 密码
	Age      int    `json:"age" validate:"required,min=1,max=150"`     // 年龄
}

// SignupResponse 注册响应数据
type SignupData struct {
	Name string `json:"name"` // 用户名
}

// RestoreSessionRequest 恢复会话请求
type RestoreSessionRequest struct {
	BaseRequest
	Name string `json:"name" validate:"required,min=1,max=50"` // 用户名
}

// NewLoginRequest 创建登录请求
func NewLoginRequest(name, password string) LoginRequest {
	return LoginRequest{
		BaseRequest: NewBaseRequest(),
		Name:        name,
		Password:    password,
	}
}

// NewSignupRequest 创建注册请求
func NewSignupRequest(name, password string, age int) SignupRequest {
	return SignupRequest{
		BaseRequest: NewBaseRequest(),
		Name:        name,
		Password:    password,
		Age:         age,
	}
}

// NewRestoreSessionRequest 创建恢复会话请求
func NewRestoreSessionRequest(name string) RestoreSessionRequest {
	return RestoreSessionRequest{
		BaseRequest: NewBaseRequest(),
		Name:        name,
	}
}

// LoginSuccess 创建登录成功响应
func LoginSuccess(name string, age int) BaseResponse {
	data := LoginData{
		Name: name,
		Age:  age,
	}
	return SuccessWithMessage(data, "登录成功")
}

// SignupSuccess 创建注册成功响应
func SignupSuccess(name string) BaseResponse {
	data := SignupData{
		Name: name,
	}
	return SuccessWithMessage(data, "注册成功")
}

// RestoreSessionSuccess 创建恢复会话成功响应
func RestoreSessionSuccess(name string, age int) BaseResponse {
	data := LoginData{
		Name: name,
		Age:  age,
	}
	return SuccessWithMessage(data, "会话恢复成功")
}
