package models

import (
	"time"
)

// User 用户数据结构
type User struct {
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Age         int       `json:"age"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}
