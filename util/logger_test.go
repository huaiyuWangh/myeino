package util

import (
	"testing"
)

// TestLogger 测试logger基本功能
func TestLogger(t *testing.T) {
	// 创建一个新的logger
	logger := NewLogger(DEBUG)

	// 测试简单日志
	logger.Debug("这是一条调试信息")
	logger.Info("应用程序启动成功")
	logger.Warn("配置文件未找到，使用默认配置")
	logger.Error("数据库连接失败")
	logger.Fatal("系统遇到致命错误")

	// 测试格式化日志
	logger.Debugf("这是一条调试信息，用户ID: %d", 12345)
	logger.Infof("应用程序启动成功，版本: %s", "v1.0.0")
	logger.Warnf("重试次数: %d", 3)
	logger.Errorf("数据库连接失败: %s", "connection timeout")
	logger.Fatalf("系统遇到致命错误，错误码: %d", 500)
}

// TestLoggerWithLevelFilter 测试日志等级过滤
func TestLoggerWithLevelFilter(t *testing.T) {
	logger := NewLogger(WARN) // 只显示WARN及以上等级的日志

	t.Log("=== 设置日志等级为WARN，以下只会显示WARN、ERROR、FATAL ===")
	logger.Debug("这条调试信息不会显示")
	logger.Info("这条信息也不会显示")
	logger.Warn("这条警告会显示")
	logger.Error("这条错误会显示")
	logger.Fatal("这条致命错误会显示")
}

// TestLoggerWithoutColor 测试无颜色输出
func TestLoggerWithoutColor(t *testing.T) {
	logger := NewLogger(DEBUG)
	logger.SetColorEnabled(false) // 禁用颜色输出

	t.Log("=== 无颜色输出测试 ===")
	logger.Debug("调试信息（无颜色）")
	logger.Info("普通信息（无颜色）")
	logger.Warn("警告信息（无颜色）")
	logger.Error("错误信息（无颜色）")

	// 测试格式化版本
	logger.Debugf("格式化调试信息（无颜色）: %s", "测试")
	logger.Infof("格式化普通信息（无颜色）: %d", 100)
}

// TestGlobalLogger 测试全局logger
func TestGlobalLogger(t *testing.T) {
	// 设置全局日志等级
	SetLogLevel(DEBUG)

	t.Log("=== 全局Logger测试 ===")
	// 测试简单日志
	Debug("全局调试信息")
	Info("全局普通信息")
	Warn("全局警告信息")
	Error("全局错误信息")
	Fatal("全局致命错误信息")

	// 测试格式化日志
	Debugf("全局格式化调试信息: %s", "测试")
	Infof("全局格式化普通信息: %d", 42)
	Warnf("全局格式化警告信息: %.2f", 3.14)
	Errorf("全局格式化错误信息: %v", "错误详情")
	Fatalf("全局格式化致命错误: %s", "系统崩溃")
}

// TestLoggerInFunction 在不同函数中测试，验证行号显示
func TestLoggerInFunction(t *testing.T) {
	testFunction1()
	testFunction2()
}

func testFunction1() {
	Info("这是在testFunction1中的日志")
	Warn("testFunction1中的警告")
	Infof("testFunction1中的格式化信息: %s", "测试参数")
}

func testFunction2() {
	Error("这是在testFunction2中的错误日志")
	Debug("testFunction2中的调试信息")
	Errorf("testFunction2中的格式化错误: %d", 404)
}