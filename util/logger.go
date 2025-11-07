package util

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// LogLevel 定义日志等级
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String 返回日志等级的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color 返回日志等级对应的颜色代码
func (l LogLevel) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // 青色
	case INFO:
		return "\033[32m" // 绿色
	case WARN:
		return "\033[33m" // 黄色
	case ERROR:
		return "\033[31m" // 红色
	case FATAL:
		return "\033[35m" // 紫色
	default:
		return "\033[0m" // 默认色
	}
}

// Logger 日志记录器结构体
type Logger struct {
	level        LogLevel
	colorEnabled bool
}

// NewLogger 创建新的日志记录器
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:        level,
		colorEnabled: true,
	}
}

// SetLevel 设置日志等级
func (logger *Logger) SetLevel(level LogLevel) {
	logger.level = level
}

// SetColorEnabled 设置是否启用颜色输出
func (logger *Logger) SetColorEnabled(enabled bool) {
	logger.colorEnabled = enabled
}

// getCallerInfo 获取调用者信息（文件名和行号）
func (logger *Logger) getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0
	}

	// 只显示文件名，不显示完整路径
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}

	return file, line
}

// logf 格式化日志记录方法
func (logger *Logger) logf(level LogLevel, format string, args ...interface{}) {
	if level < logger.level {
		return
	}

	// 获取调用者信息，skip=3 跳过当前方法、具体日志方法(Debugf/Infof等)和用户调用
	file, line := logger.getCallerInfo(3)

	// 格式化时间
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 构建日志消息
	message := fmt.Sprintf(format, args...)

	// 构建完整的日志行
	var logLine string
	if logger.colorEnabled {
		// 带颜色的输出
		colorCode := level.Color()
		resetCode := "\033[0m"
		logLine = fmt.Sprintf("%s[%s] %s%s %s:%d - %s",
			colorCode, level.String(), timestamp, resetCode, file, line, message)
	} else {
		// 无颜色输出
		logLine = fmt.Sprintf("[%s] %s %s:%d - %s",
			level.String(), timestamp, file, line, message)
	}

	fmt.Println(logLine)
}

// log 简单日志记录方法
func (logger *Logger) log(level LogLevel, message string) {
	if level < logger.level {
		return
	}

	// 获取调用者信息，skip=3 跳过当前方法、具体日志方法(Debug/Info等)和用户调用
	file, line := logger.getCallerInfo(3)

	// 格式化时间
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 构建完整的日志行
	var logLine string
	if logger.colorEnabled {
		// 带颜色的输出
		colorCode := level.Color()
		resetCode := "\033[0m"
		logLine = fmt.Sprintf("%s[%s] %s%s %s:%d - %s",
			colorCode, level.String(), timestamp, resetCode, file, line, message)
	} else {
		// 无颜色输出
		logLine = fmt.Sprintf("[%s] %s %s:%d - %s",
			level.String(), timestamp, file, line, message)
	}

	fmt.Println(logLine)
}

// logAny 接受任意类型参数的日志记录方法
func (logger *Logger) logAny(level LogLevel, args ...interface{}) {
	if level < logger.level {
		return
	}

	// 获取调用者信息，skip=3 跳过当前方法、具体日志方法和用户调用
	file, line := logger.getCallerInfo(3)

	// 格式化时间
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 将所有参数格式化为字符串
	message := fmt.Sprint(args...)

	// 构建完整的日志行
	var logLine string
	if logger.colorEnabled {
		// 带颜色的输出
		colorCode := level.Color()
		resetCode := "\033[0m"
		logLine = fmt.Sprintf("%s[%s] %s%s %s:%d - %s",
			colorCode, level.String(), timestamp, resetCode, file, line, message)
	} else {
		// 无颜色输出
		logLine = fmt.Sprintf("[%s] %s %s:%d - %s",
			level.String(), timestamp, file, line, message)
	}

	fmt.Println(logLine)
}

// Debug 记录调试日志
func (logger *Logger) Debug(message string) {
	logger.log(DEBUG, message)
}

// Debugf 记录格式化调试日志
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(DEBUG, format, args...)
}

// Info 记录信息日志
func (logger *Logger) Info(message string) {
	logger.log(INFO, message)
}

// Infof 记录格式化信息日志
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(INFO, format, args...)
}

// Warn 记录警告日志
func (logger *Logger) Warn(message string) {
	logger.log(WARN, message)
}

// Warnf 记录格式化警告日志
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(WARN, format, args...)
}

// Error 记录错误日志
func (logger *Logger) Error(message string) {
	logger.log(ERROR, message)
}

// Errorf 记录格式化错误日志
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(ERROR, format, args...)
}

// Fatal 记录致命错误日志
func (logger *Logger) Fatal(message string) {
	logger.log(FATAL, message)
}

// Fatalf 记录格式化致命错误日志
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logf(FATAL, format, args...)
}

// DebugAny 记录调试日志（接受任意类型参数）
func (logger *Logger) DebugAny(args ...interface{}) {
	logger.logAny(DEBUG, args...)
}

// InfoAny 记录信息日志（接受任意类型参数）
func (logger *Logger) InfoAny(args ...interface{}) {
	logger.logAny(INFO, args...)
}

// WarnAny 记录警告日志（接受任意类型参数）
func (logger *Logger) WarnAny(args ...interface{}) {
	logger.logAny(WARN, args...)
}

// ErrorAny 记录错误日志（接受任意类型参数）
func (logger *Logger) ErrorAny(args ...interface{}) {
	logger.logAny(ERROR, args...)
}

// FatalAny 记录致命错误日志（接受任意类型参数）
func (logger *Logger) FatalAny(args ...interface{}) {
	logger.logAny(FATAL, args...)
}

// 全局默认logger实例
var defaultLogger = NewLogger(INFO)

// SetLogLevel 设置全局日志等级
func SetLogLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetColorEnabled 设置全局颜色输出开关
func SetColorEnabled(enabled bool) {
	defaultLogger.SetColorEnabled(enabled)
}

// Debugf 全局格式化调试日志
func Debugf(format string, args ...interface{}) {
	defaultLogger.logf(DEBUG, format, args...)
}

// Infof 全局格式化信息日志
func Infof(format string, args ...interface{}) {
	defaultLogger.logf(INFO, format, args...)
}

// Warnf 全局格式化警告日志
func Warnf(format string, args ...interface{}) {
	defaultLogger.logf(WARN, format, args...)
}

// Errorf 全局格式化错误日志
func Errorf(format string, args ...interface{}) {
	defaultLogger.logf(ERROR, format, args...)
}

// Fatalf 全局格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	defaultLogger.logf(FATAL, format, args...)
}

// DebugAny 全局调试日志（接受任意类型参数）
func Debug(args ...interface{}) {
	defaultLogger.logAny(DEBUG, args...)
}

// InfoAny 全局信息日志（接受任意类型参数）
func Info(args ...interface{}) {
	defaultLogger.logAny(INFO, args...)
}

// WarnAny 全局警告日志（接受任意类型参数）
func Warn(args ...interface{}) {
	defaultLogger.logAny(WARN, args...)
}

// ErrorAny 全局错误日志（接受任意类型参数）
func Error(args ...interface{}) {
	defaultLogger.logAny(ERROR, args...)
}

// FatalAny 全局致命错误日志（接受任意类型参数）
func Fatal(args ...interface{}) {
	defaultLogger.logAny(FATAL, args...)
}
