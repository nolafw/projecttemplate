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

// TODO: 設定ファイルから読み込んだ値によって、自動的に出力レベルを変更する関数を作る
// 明示的にInfoなど指定するのではｎ
// slogにその機能があればそのまま使う
