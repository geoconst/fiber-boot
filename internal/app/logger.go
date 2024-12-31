package app

import (
	"fiber-boot/internal/utils"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/phsym/console-slog"
)

func GetLogWriter() io.Writer {
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
func InitLogger() *slog.Logger {
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
			GetLogWriter(),
			options,
		),
	)
	slog.SetDefault(logger)
	return logger
}
