package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 日志级别
type Level int

const (
	DEBUG Level = iota // 调试级别
	INFO               // 信息级别
	WARN               // 警告级别
	ERROR              // 错误级别
	FATAL              // 致命错误级别
)

// 日志级别字符串表示
var levelStrings = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger 结构体
type Logger struct {
	level  Level
	logger *log.Logger
	file   *os.File
}

// 全局logger实例
var defaultLogger *Logger

// 初始化默认logger
func init() {
	defaultLogger = New(INFO, "")
}

// 创建新的Logger实例
func New(level Level, logFile string) *Logger {
	l := &Logger{
		level: level,
	}

	var writer io.Writer = os.Stdout

	// 如果指定了日志文件，则写入文件
	if logFile != "" {
		// 确保日志目录存在
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
		} else {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("打开日志文件失败: %v\n", err)
			} else {
				l.file = file
				// 同时写入控制台和文件
				writer = io.MultiWriter(os.Stdout, file)
			}
		}
	}

	l.logger = log.New(writer, "", 0)
	return l
}

// 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// 格式化日志消息
func (l *Logger) formatMessage(level Level, format string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelStrings[level]
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, levelStr, message)
}

// 写入日志
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level >= l.level {
		message := l.formatMessage(level, format, args...)
		l.logger.Println(message)

		// 如果是FATAL级别，程序退出
		if level == FATAL {
			if l.file != nil {
				l.file.Close()
			}
			os.Exit(1)
		}
	}
}

// Debug 输出调试信息
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 输出信息
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 输出警告
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 输出错误
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 输出致命错误并退出程序
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// 以下是全局函数，使用默认logger实例

// SetLevel 设置全局日志级别
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetLogFile 设置全局日志文件
func SetLogFile(logFile string) {
	newLogger := New(defaultLogger.level, logFile)
	if defaultLogger.file != nil {
		defaultLogger.file.Close()
	}
	defaultLogger = newLogger
}

// Debug 输出调试信息
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info 输出信息
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn 输出警告
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error 输出错误
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal 输出致命错误并退出程序
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// Close 关闭全局logger
func Close() {
	defaultLogger.Close()
}
