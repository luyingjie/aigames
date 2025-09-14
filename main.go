package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	jsonSerializer "github.com/lonng/nano/serialize/json"
	"github.com/lonng/nano/session"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB

type (
	// Handler 处理器结构体
	Handler struct {
		component.Base
	}

	// User 用户数据结构
	User struct {
		Name        string    `json:"name"`
		Password    string    `json:"password"`
		Age         int       `json:"age"`
		CreatedAt   time.Time `json:"created_at"`
		LastLoginAt time.Time `json:"last_login_at"`
	}

	// LoginRequest 登录请求结构
	LoginRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// LoginResponse 登录响应结构  
	LoginResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// SignupRequest 注册请求结构
	SignupRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Age      int    `json:"age"`
	}

	// SignupResponse 注册响应结构
	SignupResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Name string `json:"name"`
	}
)

// 数据库初始化
func initDB() error {
	var err error
	db, err = bbolt.Open("game.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// 创建用户表（bucket）
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})
}

// 关闭数据库
func closeDB() {
	if db != nil {
		db.Close()
	}
}

// 加密密码
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// 保存用户到数据库
func saveUser(user *User) error {
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
func getUser(name string) (*User, error) {
	var user User
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
func userExists(name string) bool {
	_, err := getUser(name)
	return err == nil
}

// 更新用户最后登录时间
func updateLastLogin(name string) error {
	user, err := getUser(name)
	if err != nil {
		return err
	}
	user.LastLoginAt = time.Now()
	return saveUser(user)
}

// Login 登录处理方法
func (h *Handler) Login(s *session.Session, req *LoginRequest) error {
	log.Printf("用户登录请求: %s", req.Name)
	
	// 验证输入
	if req.Name == "" {
		resp := &LoginResponse{
			Code: 400,
			Msg:  "用户名不能为空",
		}
		return s.Response(resp)
	}
	
	if req.Password == "" {
		resp := &LoginResponse{
			Code: 400,
			Msg:  "密码不能为空",
		}
		return s.Response(resp)
	}

	// 从数据库获取用户
	user, err := getUser(req.Name)
	if err != nil {
		resp := &LoginResponse{
			Code: 404,
			Msg:  "用户不存在",
		}
		return s.Response(resp)
	}

	// 验证密码
	hashedPassword := hashPassword(req.Password)
	if user.Password != hashedPassword {
		resp := &LoginResponse{
			Code: 401,
			Msg:  "密码错误",
		}
		return s.Response(resp)
	}

	// 更新最后登录时间
	if err := updateLastLogin(req.Name); err != nil {
		log.Printf("更新登录时间失败: %v", err)
	}

	// 登录成功
	resp := &LoginResponse{
		Code: 200,
		Msg:  "登录成功",
		Name: user.Name,
		Age:  user.Age,
	}
	
	log.Printf("用户 %s 登录成功", req.Name)
	return s.Response(resp)
}

// Signup 注册处理方法
func (h *Handler) Signup(s *session.Session, req *SignupRequest) error {
	log.Printf("用户注册请求: %s", req.Name)
	
	// 验证输入
	if req.Name == "" {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "用户名不能为空",
		}
		return s.Response(resp)
	}
	
	if req.Password == "" {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "密码不能为空",
		}
		return s.Response(resp)
	}
	
	if req.Age <= 0 || req.Age > 150 {
		resp := &SignupResponse{
			Code: 400,
			Msg:  "年龄必须在1-150之间",
		}
		return s.Response(resp)
	}

	// 检查用户是否已存在
	if userExists(req.Name) {
		resp := &SignupResponse{
			Code: 409,
			Msg:  "用户名已存在",
		}
		return s.Response(resp)
	}

	// 创建新用户
	user := &User{
		Name:        req.Name,
		Password:    hashPassword(req.Password),
		Age:         req.Age,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}

	// 保存用户到数据库
	if err := saveUser(user); err != nil {
		log.Printf("保存用户失败: %v", err)
		resp := &SignupResponse{
			Code: 500,
			Msg:  "注册失败，请稍后重试",
		}
		return s.Response(resp)
	}

	// 注册成功
	resp := &SignupResponse{
		Code: 200,
		Msg:  "注册成功",
		Name: user.Name,
	}
	
	log.Printf("用户 %s 注册成功", req.Name)
	return s.Response(resp)
}

func main() {
	// 初始化数据库
	if err := initDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer closeDB()

	log.Println("Nano 游戏服务器启动中...")
	log.Println("WebSocket 服务监听端口: 3250")
	log.Println("Web 静态文件服务监听端口: 8080")
	log.Println("数据库文件: game.db")
	
	// 启动静态文件服务器为前端页面提供服务
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./web/")))
		log.Println("静态文件服务器启动在 http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("静态文件服务器启动失败:", err)
		}
	}()

	// 创建组件容器并注册处理器
	components := &component.Components{}
	components.Register(&Handler{},
		component.WithName("gate"),
	)

	// 启动nano WebSocket服务器
	nano.Listen(":3250",
		nano.WithIsWebsocket(true),
		nano.WithComponents(components),
		nano.WithSerializer(jsonSerializer.NewSerializer()),
		nano.WithWSPath("/nano"),
	)
}