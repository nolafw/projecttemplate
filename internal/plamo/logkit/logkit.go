package logkit

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
)

var (
	logger     *slog.Logger
	logLevel   = slog.LevelInfo
	onceForLog sync.Once
)

const LvDebug = "debug"
const LvInfo = "info"
const LvWarn = "warn"
const LvErr = "error"

// どのloggerを使うかは自由です。
// 用途にあったloggerを定義してください。
// ここで作成されているloggerをそのまま使っても問題ありませんが、
// 用途に合わせてloggerを作成・カスタマイズしてください。
func Logger() *slog.Logger {
	onceForLog.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return logger
}

// TEST: ログレベルに従った出力になるか確認
func SetLogLevel(level string) {
	lowerLevel := strings.ToLower(level)

	switch lowerLevel {
	case LvDebug:
		logLevel = slog.LevelDebug
	case LvInfo:
		logLevel = slog.LevelInfo
	case LvWarn:
		logLevel = slog.LevelWarn
	case LvErr:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	if !slices.Contains([]string{LvDebug, LvInfo, LvWarn, LvErr}, lowerLevel) {
		fmt.Printf("[NOLA LOG WARN]: Invalid log level: [%s]. Log level falls back to [info]\n", level)
	} else {
		fmt.Printf("[NOLA LOG INFO]: Log level set to [%s]\n", lowerLevel)
	}
}

func Debug(msg string, args ...any) {
	if logLevel <= slog.LevelDebug {
		Logger().Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if logLevel <= slog.LevelInfo {
		Logger().Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if logLevel <= slog.LevelWarn {
		Logger().Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if logLevel <= slog.LevelError {
		Logger().Error(msg, args...)
	}
}
