package logkit

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger     *slog.Logger
	onceForLog sync.Once
)

// どのloggerを使うかは自由です。
// 用途にあったloggerを定義してください。
func Logger() *slog.Logger {
	onceForLog.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return logger
}

func Info(msg string, args ...any) {
	Logger().Info(msg, args...)
}
