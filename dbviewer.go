package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"aigames/internal/models"

	"go.etcd.io/bbolt"
)

func main_1() {
	// 打开数据库
	db, err := bbolt.Open("./data/game.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal("打开数据库失败:", err)
	}
	defer db.Close()

	fmt.Println("🔍 游戏数据库查看工具")
	fmt.Println("===================")

	// 查看所有用户
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			fmt.Println("暂无用户数据")
			return nil
		}

		fmt.Printf("📊 用户总数: %d\n\n", b.Stats().KeyN)

		return b.ForEach(func(k, v []byte) error {
			var user models.User
			if err := json.Unmarshal(v, &user); err != nil {
				fmt.Printf("❌ 解析用户数据失败: %v\n", err)
				return nil
			}

			fmt.Printf("👤 用户: %s\n", user.Name)
			fmt.Printf("   年龄: %d\n", user.Age)
			fmt.Printf("   注册时间: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("   最后登录: %s\n", user.LastLoginAt.Format("2006-01-02 15:04:05"))
			fmt.Println("   ---")
			return nil
		})
	})

	if err != nil {
		log.Fatal("查询数据库失败:", err)
	}
}
