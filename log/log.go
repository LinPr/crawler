package log

// zapcore is a core package of the Zap logger, which provides the fundamental interfaces and types for creating a logger in Go. It includes the following components:

// Encoder: defines the encoding format of the logs such as JSON, console, and custom encoders.
// WriteSyncer: specifies the destination for the logs such as a file, standard output, and network socket.
// Level: defines the log level such as DEBUG, INFO, WARN, ERROR, and FATAL.
// Core: composes the above components to form a logger.

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Plugin zapcore.Core

// 根据插件类型，创建对应的logger
func NewLogger(plugin zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(plugin, append(DefaultOption(), options...)...)
}

// 将Core的Encoder设置为JSON，writer 和 enabler待定，作为插件使用
func NewPlugin(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) Plugin {
	return zapcore.NewCore(DefaultEncoder(), writer, enabler)
}

func NewStdoutPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stdout)), enabler)
}

func NewStderrPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stderr)), enabler)
}

// Lumberjack logger虽然持有File但没有暴露sync方法，所以没办法利用zap的sync特性
// 所以额外返回一个closer，需要保证在进程退出前close以保证写入的内容可以全部刷到到磁盘
func NewFilePlugin(filePath string, enabler zapcore.LevelEnabler) (Plugin, io.Closer) {
	var writer = DefaultLumberjackLogger()
	writer.Filename = filePath
	return NewPlugin(zapcore.AddSync(writer), enabler), writer
}
