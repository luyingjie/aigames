package main

import (
	"aigames/internal/config"
	"aigames/internal/database"
	"aigames/internal/handlers"
	"aigames/internal/services"
	"aigames/pkg/logger"
	"net/http"
	"strconv"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	jsonSerializer "github.com/lonng/nano/serialize/json"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(""); err != nil {
		logger.Fatal("加载配置失败: %v", err)
	}

	cfg := config.GetConfig()

	// 确保必要的目录存在
	if err := cfg.EnsureDirs(); err != nil {
		logger.Fatal("创建目录失败: %v", err)
	}

	// 设置日志
	logger.SetLogFile(cfg.Log.FilePath)
	if cfg.IsDevelopment() {
		logger.SetLevel(logger.DEBUG)
	} else {
		logger.SetLevel(logger.INFO)
	}

	logger.Info("启动AI游戏服务器...")
	logger.Info("配置信息: 模式=%s, 端口=%d", cfg.Server.Mode, cfg.Server.Port)

	logger.Info("Nano 游戏服务器启动中...")
	logger.Info("WebSocket 服务监听端口: 3250")
	logger.Info("Web 静态文件服务监听端口: 8080")
	logger.Info("数据库文件: game.db")

	// 连接数据库
	db, err := database.NewDB(cfg.Database.Path)
	if err != nil {
		logger.Fatal("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 创建服务实例
	userService := services.NewUserService(db.GetBoltDB())

	// 启动静态文件服务器为前端页面提供服务
	go func() {
		http.Handle("/", http.FileServer(http.Dir("./web/")))
		logger.Info("静态文件服务器启动在 http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			logger.Fatal("静态文件服务器启动失败:", err)
		}
	}()

	// 创建组件容器并注册处理器
	components := &component.Components{}
	components.Register(handlers.NewUser(userService),
		component.WithName("user"),
	)

	// 启动nano WebSocket服务器
	nano.Listen(":"+strconv.Itoa(cfg.Server.Port),
		nano.WithIsWebsocket(true),
		nano.WithComponents(components),
		nano.WithSerializer(jsonSerializer.NewSerializer()),
		nano.WithWSPath("/nano"),
	)
}
