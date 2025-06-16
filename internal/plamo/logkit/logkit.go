package logkit

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	logger     *slog.Logger
	levelVar   = &slog.LevelVar{} // 動的にレベルを変更可能
	onceForLog sync.Once
)

// どのloggerを使うかは自由です。
// 用途にあったloggerを定義してください。
// ここで作成されているloggerをそのまま使っても問題ありませんが、
// 用途に合わせてloggerを作成・カスタマイズしてください。
func Logger() *slog.Logger {
	onceForLog.Do(func() {
		levelVar.Set(slog.LevelDebug)
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: levelVar,
		})
		logger = slog.New(handler)
	})
	return logger
}

func SetLogLevel(level string) {
	var slogLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	levelVar.Set(slogLevel)
}

func Debug(msg string, args ...any) {
	Logger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Logger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger().Error(msg, args...)
}
