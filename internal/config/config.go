package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // 服务器配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库配置
	AI       AIConfig       `mapstructure:"ai"`       // AI配置
	Log      LogConfig      `mapstructure:"log"`      // 日志配置
	Game     GameConfig     `mapstructure:"game"`     // 游戏配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"` // 服务器主机
	Port int    `mapstructure:"port"` // 服务器端口
	Mode string `mapstructure:"mode"` // 运行模式：debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path    string `mapstructure:"path"`    // BoltDB文件路径
	Timeout int    `mapstructure:"timeout"` // 操作超时(秒)
}

// AIConfig AI配置
type AIConfig struct {
	DefaultModel       string  `mapstructure:"default_model"`       // 默认AI模型
	MaxConcurrent      int     `mapstructure:"max_concurrent"`      // 最大并发请求数
	Timeout            int     `mapstructure:"timeout"`             // 请求超时(秒)
	DefaultThinkTime   int     `mapstructure:"default_think_time"`  // 默认思考时间(秒)
	DefaultTemperature float64 `mapstructure:"default_temperature"` // 默认创造性参数
	MaxTokens          int     `mapstructure:"max_tokens"`          // 最大token数
	// 新增AI API配置
	APIURL   string `mapstructure:"api_url"`  // AI API地址
	APIKey   string `mapstructure:"api_key"`  // API密钥
	Provider string `mapstructure:"provider"` // AI服务提供商
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"`     // 日志级别：debug, info, warn, error
	FilePath string `mapstructure:"file_path"` // 日志文件路径
	MaxSize  int    `mapstructure:"max_size"`  // 最大文件大小(MB)
	MaxAge   int    `mapstructure:"max_age"`   // 最大保存天数
	Compress bool   `mapstructure:"compress"`  // 是否压缩
}

// GameConfig 游戏配置
type GameConfig struct {
	DefaultRoomCapacity   int `mapstructure:"default_room_capacity"`   // 默认房间容量
	DefaultGameTimeout    int `mapstructure:"default_game_timeout"`    // 默认游戏超时(秒)
	DefaultBiddingTimeout int `mapstructure:"default_bidding_timeout"` // 默认叫地主超时(秒)
	DefaultPlayTimeout    int `mapstructure:"default_play_timeout"`    // 默认出牌超时(秒)
	MaxRoomsPerUser       int `mapstructure:"max_rooms_per_user"`      // 用户最大房间数
	MaxAIPlayersPerRoom   int `mapstructure:"max_ai_players_per_room"` // 房间最大AI玩家数
}

var (
	// 全局配置实例
	AppConfig *Config
)

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 如果指定了配置文件路径
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 默认配置文件搜索路径
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../configs")
		viper.AddConfigPath("../../configs")
		viper.AddConfigPath(".")
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("AIGAME")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认配置
			fmt.Println("配置文件未找到，使用默认配置")
		} else {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 解析配置到结构体
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 处理相对路径
	if config.Database.Path != "" && !filepath.IsAbs(config.Database.Path) {
		config.Database.Path = filepath.Join(".", config.Database.Path)
	}

	if config.Log.FilePath != "" && !filepath.IsAbs(config.Log.FilePath) {
		config.Log.FilePath = filepath.Join(".", config.Log.FilePath)
	}

	AppConfig = config
	return nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 120)

	// 数据库默认配置
	viper.SetDefault("database.path", "./data/game.db")
	viper.SetDefault("database.timeout", 5)

	// JWT默认配置
	viper.SetDefault("jwt.secret", "aigame-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_time", 24)
	viper.SetDefault("jwt.issuer", "ai-game")

	// AI默认配置
	viper.SetDefault("ai.default_model", "gpt-3.5-turbo")
	viper.SetDefault("ai.max_concurrent", 10)
	viper.SetDefault("ai.timeout", 30)
	viper.SetDefault("ai.default_think_time", 3)
	viper.SetDefault("ai.default_temperature", 0.7)
	viper.SetDefault("ai.max_tokens", 1000)
	// AI API默认配置
	viper.SetDefault("ai.api_url", "https://api.openai.com/v1/chat/completions")
	viper.SetDefault("ai.api_key", "")
	viper.SetDefault("ai.provider", "openai")

	// WebSocket默认配置
	viper.SetDefault("websocket.read_buffer_size", 1024)
	viper.SetDefault("websocket.write_buffer_size", 1024)
	viper.SetDefault("websocket.heartbeat_interval", 30)
	viper.SetDefault("websocket.pong_wait", 60)
	viper.SetDefault("websocket.write_wait", 10)
	viper.SetDefault("websocket.max_message_size", 512)

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file_path", "./logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)

	// 游戏默认配置
	viper.SetDefault("game.default_room_capacity", 3)
	viper.SetDefault("game.default_game_timeout", 1800)
	viper.SetDefault("game.default_bidding_timeout", 30)
	viper.SetDefault("game.default_play_timeout", 60)
	viper.SetDefault("game.max_rooms_per_user", 10)
	viper.SetDefault("game.max_ai_players_per_room", 2)
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return AppConfig
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment 是否为开发模式
func (c *Config) IsDevelopment() bool {
	return c.Server.Mode == "debug"
}

// IsProduction 是否为生产模式
func (c *Config) IsProduction() bool {
	return c.Server.Mode == "release"
}

// GetDatabaseDir 获取数据库目录
func (c *Config) GetDatabaseDir() string {
	return filepath.Dir(c.Database.Path)
}

// GetLogDir 获取日志目录
func (c *Config) GetLogDir() string {
	return filepath.Dir(c.Log.FilePath)
}

// EnsureDirs 确保必要的目录存在
func (c *Config) EnsureDirs() error {
	// 创建数据库目录
	if err := os.MkdirAll(c.GetDatabaseDir(), 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 创建日志目录
	if err := os.MkdirAll(c.GetLogDir(), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	return nil
}
