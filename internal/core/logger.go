package core

import (
	"fiber-boot/internal/utils"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/phsym/console-slog"
)

// 配置日志writer
func logWriter() io.Writer {
	var output io.Writer = os.Stderr
	if utils.IsProd() {
		output = &lumberjack.Logger{
			Filename:   "logs/server.log",
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
	}
	return output
}

// 初始化日志配置
func InitLogger(config *Config) *slog.Logger {
	options := &console.HandlerOptions{Level: slog.LevelInfo, AddSource: true}
	if utils.IsProd() {
		options = &console.HandlerOptions{
			Level:     slog.LevelWarn,
			AddSource: true,
			NoColor:   true,
		}
	}
	logger := slog.New(
		console.NewHandler(
			logWriter(),
			options,
		),
	)
	slog.SetDefault(logger)
	return logger
}
