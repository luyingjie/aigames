package models

import (
	"time"
)

// User 用户数据结构
type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
	Age         int       `json:"age"`
	IsAI        bool      `json:"is_ai"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}
