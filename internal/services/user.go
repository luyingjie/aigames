package services

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"aigames/internal/models"

	"go.etcd.io/bbolt"
)

var db *bbolt.DB

// 加密密码
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// 保存用户到数据库
func SaveUser(user *models.User) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.Name), userJSON)
	})
}

// 从数据库获取用户
func GetUser(name string) (*models.User, error) {
	var user models.User
	err := db.View(func(tx *bbolt.Tx) error {
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

// 检查用户是否存在
func UserExists(name string) bool {
	_, err := GetUser(name)
	return err == nil
}

// 更新用户最后登录时间
func UpdateLastLogin(name string) error {
	user, err := GetUser(name)
	if err != nil {
		return err
	}
	user.LastLoginAt = time.Now()
	return SaveUser(user)
}
