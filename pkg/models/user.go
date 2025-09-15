package models

import "aigames/pkg/protocol"

// 用户相关的请求和响应结构体

// LoginRequest 登录请求
type LoginRequest struct {
	protocol.BaseRequest
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
	protocol.BaseRequest
	Name     string `json:"name" validate:"required,min=1,max=50"`     // 用户名
	Password string `json:"password" validate:"required,min=6,max=50"` // 密码
	Age      int    `json:"age" validate:"required,min=1,max=150"`     // 年龄
}

// SignupResponse 注册响应数据
type SignupData struct {
	Name string `json:"name"` // 用户名
}

// UserInfoRequest 获取用户信息请求
type UserInfoRequest struct {
	protocol.BaseRequest
	Name string `json:"name" validate:"required"` // 用户名
}

// UserInfoData 用户信息数据
type UserInfoData struct {
	Name        string `json:"name"`          // 用户名
	Age         int    `json:"age"`           // 年龄
	CreatedAt   string `json:"created_at"`    // 创建时间
	LastLoginAt string `json:"last_login_at"` // 最后登录时间
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	protocol.PageRequest
	Keyword string `json:"keyword,omitempty"` // 搜索关键词
}

// UserListData 用户列表数据
type UserListData struct {
	Users []UserInfoData `json:"users"` // 用户列表
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	protocol.BaseRequest
	Name string `json:"name" validate:"required"`               // 用户名
	Age  int    `json:"age,omitempty" validate:"min=1,max=150"` // 年龄（可选）
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	protocol.BaseRequest
	Name        string `json:"name" validate:"required"`                      // 用户名
	OldPassword string `json:"old_password" validate:"required"`              // 旧密码
	NewPassword string `json:"new_password" validate:"required,min=6,max=50"` // 新密码
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	protocol.BaseRequest
	Name string `json:"name" validate:"required"` // 用户名
}

// NewLoginRequest 创建登录请求
func NewLoginRequest(name, password string) LoginRequest {
	return LoginRequest{
		BaseRequest: protocol.NewBaseRequest(),
		Name:        name,
		Password:    password,
	}
}

// NewSignupRequest 创建注册请求
func NewSignupRequest(name, password string, age int) SignupRequest {
	return SignupRequest{
		BaseRequest: protocol.NewBaseRequest(),
		Name:        name,
		Password:    password,
		Age:         age,
	}
}

// NewUserInfoRequest 创建获取用户信息请求
func NewUserInfoRequest(name string) UserInfoRequest {
	return UserInfoRequest{
		BaseRequest: protocol.NewBaseRequest(),
		Name:        name,
	}
}

// NewUserListRequest 创建用户列表请求
func NewUserListRequest(page, size int, keyword string) UserListRequest {
	return UserListRequest{
		PageRequest: protocol.NewPageRequest(page, size),
		Keyword:     keyword,
	}
}

// LoginSuccess 创建登录成功响应
func LoginSuccess(name string, age int) protocol.BaseResponse {
	data := LoginData{
		Name: name,
		Age:  age,
	}
	return protocol.SuccessWithMessage(data, "登录成功")
}

// SignupSuccess 创建注册成功响应
func SignupSuccess(name string) protocol.BaseResponse {
	data := SignupData{
		Name: name,
	}
	return protocol.SuccessWithMessage(data, "注册成功")
}

// UserInfoSuccess 创建获取用户信息成功响应
func UserInfoSuccess(userInfo UserInfoData) protocol.BaseResponse {
	return protocol.SuccessWithMessage(userInfo, "获取用户信息成功")
}

// UserListSuccess 创建用户列表成功响应
func UserListSuccess(users []UserInfoData, total, page, size int) protocol.PageResponse {
	data := UserListData{
		Users: users,
	}
	return protocol.PageResponse{
		BaseResponse: protocol.SuccessWithMessage(data, "获取用户列表成功"),
		Total:        total,
		Page:         page,
		Size:         size,
	}
}

// UpdateUserSuccess 创建更新用户信息成功响应
func UpdateUserSuccess() protocol.BaseResponse {
	return protocol.SuccessWithMessage(nil, "更新用户信息成功")
}

// ChangePasswordSuccess 创建修改密码成功响应
func ChangePasswordSuccess() protocol.BaseResponse {
	return protocol.SuccessWithMessage(nil, "修改密码成功")
}

// DeleteUserSuccess 创建删除用户成功响应
func DeleteUserSuccess() protocol.BaseResponse {
	return protocol.SuccessWithMessage(nil, "删除用户成功")
}
