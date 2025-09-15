package services

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"aigames/internal/models"

	"go.etcd.io/bbolt"
)

// UserService 用户服务结构体
type UserService struct {
	db *bbolt.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *bbolt.DB) *UserService {
	return &UserService{db: db}
}

// HashPassword 加密密码
func (s *UserService) HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// SaveUser 保存用户到数据库
func (s *UserService) SaveUser(user *models.User) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.Name), userJSON)
	})
}

// GetUser 从数据库获取用户
func (s *UserService) GetUser(name string) (*models.User, error) {
	var user models.User
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		userData := b.Get([]byte(name))
		if userData == nil {
			return fmt.Errorf("用户不存在")
		}
		return json.Unmarshal(userData, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UserExists 检查用户是否存在
func (s *UserService) UserExists(name string) bool {
	_, err := s.GetUser(name)
	return err == nil
}

// UpdateLastLogin 更新用户最后登录时间
func (s *UserService) UpdateLastLogin(name string) error {
	user, err := s.GetUser(name)
	if err != nil {
		return err
	}
	user.LastLoginAt = time.Now()
	return s.SaveUser(user)
}
