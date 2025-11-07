# Logger 使用说明

这是一个支持日志等级和行号显示的Go语言日志记录器，具有以下特性：

## 主要特性

- ✅ **日志等级支持**: DEBUG, INFO, WARN, ERROR, FATAL
- ✅ **行号显示**: 自动显示调用日志的文件名和行号
- ✅ **彩色输出**: 不同日志等级使用不同颜色（可开关）
- ✅ **时间戳**: 精确到毫秒的时间戳
- ✅ **等级过滤**: 只显示指定等级及以上的日志
- ✅ **全局和实例模式**: 支持全局日志和自定义logger实例

## 快速开始

### 1. 基本使用（全局Logger）

```go
package main

import "data-analyze/util"

func main() {
    // 简单日志（无需格式化参数）
    util.Debug("调试信息")
    util.Info("程序启动成功")
    util.Warn("这是一个警告")
    util.Error("发生了错误")
    util.Fatal("致命错误")

    // 格式化日志（需要参数）
    util.Debugf("调试信息，用户ID: %d", 12345)
    util.Infof("程序启动成功，版本: %s", "v1.0.0")
    util.Warnf("警告: 剩余空间 %.1f%%", 15.5)
    util.Errorf("错误: %s", "数据库连接失败")
    util.Fatalf("致命错误，代码: %d", 500)
}
```

### 2. 创建自定义Logger

```go
package main

import "data-analyze/util"

func main() {
    // 创建自定义logger，设置为DEBUG等级
    logger := util.NewLogger(util.DEBUG)

    // 简单日志
    logger.Info("这是自定义logger的信息")
    logger.Warn("这是警告信息")

    // 格式化日志
    logger.Infof("格式化信息: %s", "成功")
    logger.Warnf("格式化警告: %d 次重试", 3)
}
```

## 日志等级

| 等级 | 数值 | 颜色 | 说明 |
|------|------|------|------|
| DEBUG | 0 | 青色 | 调试信息 |
| INFO | 1 | 绿色 | 一般信息 |
| WARN | 2 | 黄色 | 警告信息 |
| ERROR | 3 | 红色 | 错误信息 |
| FATAL | 4 | 紫色 | 致命错误 |

## 配置选项

### 设置日志等级

```go
// 全局设置
util.SetLogLevel(util.WARN) // 只显示WARN及以上等级

// 自定义logger设置
logger := util.NewLogger(util.INFO)
logger.SetLevel(util.ERROR) // 只显示ERROR及以上等级
```

### 开关颜色输出

```go
// 全局设置
util.SetColorEnabled(false) // 禁用颜色

// 自定义logger设置
logger := util.NewLogger(util.INFO)
logger.SetColorEnabled(false) // 禁用颜色
```

## 输出格式

### 带颜色输出
```
[DEBUG] 2025-11-03 14:58:10.440 logger_example.go:46 - 这是调试信息
[INFO] 2025-11-03 14:58:10.443 logger_example.go:47 - 应用程序启动成功
[WARN] 2025-11-03 14:58:10.443 logger_example.go:48 - 配置文件未找到
[ERROR] 2025-11-03 14:58:10.443 logger_example.go:49 - 数据库连接失败
```

### 无颜色输出
```
[DEBUG] 2025-11-03 14:58:10.440 logger_example.go:46 - 这是调试信息
[INFO] 2025-11-03 14:58:10.443 logger_example.go:47 - 应用程序启动成功
[WARN] 2025-11-03 14:58:10.443 logger_example.go:48 - 配置文件未找到
[ERROR] 2025-11-03 14:58:10.443 logger_example.go:49 - 数据库连接失败
```

## 两套API设计

### 简单日志 API (无格式化参数)
用于输出固定文本，无需参数：
```go
util.Debug("调试信息")
util.Info("普通信息")
util.Warn("警告信息")
util.Error("错误信息")
util.Fatal("致命错误")
```

### 格式化日志 API (带f后缀)
用于需要参数格式化的场景：
```go
util.Debugf("调试信息: %s", variable)
util.Infof("用户 %s 登录成功，ID: %d", "张三", 12345)
util.Warnf("内存使用率: %.1f%%", 85.6)
util.Errorf("连接超时，重试次数: %d，错误: %v", 3, err)
util.Fatalf("系统错误，代码: %d", 500)
```

## 运行示例

```bash
# 运行测试
go test ./util -v

# 运行示例程序
go run examples/logger_example.go
```

## API 参考

### 全局函数

#### 设置函数
- `SetLogLevel(level LogLevel)` - 设置全局日志等级
- `SetColorEnabled(enabled bool)` - 设置全局颜色开关

#### 简单日志函数
- `Debug(message string)` - 调试日志
- `Info(message string)` - 信息日志
- `Warn(message string)` - 警告日志
- `Error(message string)` - 错误日志
- `Fatal(message string)` - 致命错误日志

#### 格式化日志函数
- `Debugf(format string, args ...interface{})` - 格式化调试日志
- `Infof(format string, args ...interface{})` - 格式化信息日志
- `Warnf(format string, args ...interface{})` - 格式化警告日志
- `Errorf(format string, args ...interface{})` - 格式化错误日志
- `Fatalf(format string, args ...interface{})` - 格式化致命错误日志

### Logger 方法

#### 创建和配置
- `NewLogger(level LogLevel) *Logger` - 创建新的logger实例
- `SetLevel(level LogLevel)` - 设置日志等级
- `SetColorEnabled(enabled bool)` - 设置颜色开关

#### 简单日志方法
- `Debug(message string)` - 调试日志
- `Info(message string)` - 信息日志
- `Warn(message string)` - 警告日志
- `Error(message string)` - 错误日志
- `Fatal(message string)` - 致命错误日志

#### 格式化日志方法
- `Debugf(format string, args ...interface{})` - 格式化调试日志
- `Infof(format string, args ...interface{})` - 格式化信息日志
- `Warnf(format string, args ...interface{})` - 格式化警告日志
- `Errorf(format string, args ...interface{})` - 格式化错误日志
- `Fatalf(format string, args ...interface{})` - 格式化致命错误日志

## 最佳实践

1. **生产环境建议设置为INFO等级**，避免过多调试信息
2. **开发环境可设置为DEBUG等级**，便于调试
3. **在CI/CD环境中建议禁用颜色输出**
4. **使用有意义的日志信息**，包含必要的上下文
5. **合理使用不同等级**：
   - DEBUG: 详细的执行流程
   - INFO: 重要的业务事件
   - WARN: 潜在问题但不影响运行
   - ERROR: 错误但程序可以继续
   - FATAL: 致命错误，程序需要停止

6. **选择合适的API**：
   - 使用简单API（如`Info()`）处理固定文本
   - 使用格式化API（如`Infof()`）处理需要参数的场景
   - 避免在简单API中使用格式化字符串