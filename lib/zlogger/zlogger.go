package zlogger

import (
    "fmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger

func init() {
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
        EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
    }

    // 设置日志级别
    atom := zap.NewAtomicLevelAt(zap.DebugLevel)

    config := zap.Config{
        Level:            atom,                                                // 日志级别
        Development:      true,                                                // 开发模式，堆栈跟踪
        Encoding:         "json",                                              // 输出格式 console 或 json
        EncoderConfig:    encoderConfig,                                       // 编码器配置
        InitialFields:    map[string]interface{}{},                            // 初始化字段，如：添加一个服务器名称
        OutputPaths:      []string{"stdout"},                                  // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
        ErrorOutputPaths: []string{"stderr"},
    }


    // 构建日志
    logger, err := config.Build(zap.AddCallerSkip(1))
    if err != nil {
        panic(fmt.Sprintf("log 初始化失败: %v", err))
    }
    defer logger.Sync() // flushes buffer, if any

    sugaredLogger = logger.Sugar()
}

// Fatal logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatal(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panic logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panic(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Panicw(msg, keysAndValues...)
}

// Error logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Error(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Errorw(msg, keysAndValues...)
}

// Warn logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warn(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Warnw(msg, keysAndValues...)
}

// Info logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Info(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Infow(msg, keysAndValues...)
}

// Debug logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debug(msg string, keysAndValues ...interface{}) {
    defer handlePanic()
    sugaredLogger.Debugw(msg, keysAndValues...)
}



// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
    sugaredLogger.Fatalf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
    sugaredLogger.Panicf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
    sugaredLogger.Errorf(template, args...)
}

// Warningf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
    sugaredLogger.Warnf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
    sugaredLogger.Infof(template, args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
    sugaredLogger.Debugf(template, args...)
}

// handle panic when parameter's number is even
func handlePanic() {
    if err := recover(); err != nil {
        //fmt.Printf("%s\n", err)
    }
}
