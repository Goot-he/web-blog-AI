package clog

import (
	"core-server/internal/config"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ModeConsole = "console"
	ModeFile    = "file"
)

type Log struct {
	zl *zap.Logger
}

func NewLog(cfg *config.Config) (*Log, error) {
	level, err := parseLevel(cfg.Log.Level)
	if err != nil {
		return nil, err
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.Log.Mode))
	if mode == "" {
		mode = ModeConsole
	}

	var core zapcore.Core
	switch mode {
	case ModeConsole:
		core = newConsoleCore(level)
	case ModeFile:
		fileCore, err := newFileCore(level, cfg.Log.FileName)
		if err != nil {
			return nil, err
		}
		core = fileCore
	default:
		return nil, fmt.Errorf("unsupported log mde %q, want console|file", cfg.Log.Mode)
	}

	zl := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return &Log{zl: zl}, nil
}

func (l *Log) Zap() *zap.Logger {
	if l == nil || l.zl == nil {
		return zap.NewNop()
	}
	return l.zl
}

func (l *Log) Sugar() *zap.SugaredLogger {
	return l.Zap().Sugar()
}

// With 绑定字段，返回子 Log（便于按 component 注入）
func (l *Log) With(fields ...zap.Field) *Log {
	return &Log{zl: l.Zap().With(fields...)}
}

// Sync 刷盘，进程退出前调用
func (l *Log) Sync() error {
	if l == nil || l.zl == nil {
		return nil
	}
	return l.zl.Sync()
}

func (l *Log) Debug(msg string, fields ...zap.Field) {
	l.zl.Debug(msg, fields...)
}
func (l *Log) Info(msg string, fields ...zap.Field) {
	l.zl.Info(msg, fields...)
}
func (l *Log) Warn(msg string, fields ...zap.Field) {
	l.zl.Warn(msg, fields...)
}
func (l *Log) Error(msg string, fields ...zap.Field) {
	l.zl.Error(msg, fields...)
}
func (l *Log) Fatal(msg string, fields ...zap.Field) {
	l.zl.Fatal(msg, fields...)
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "info":
		return zapcore.InfoLevel, nil
	case "debug":
		return zapcore.DebugLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		var l zapcore.Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			return zapcore.InfoLevel, fmt.Errorf("invalid log level %q: %w", level, err)
		}
		return l, nil
	}
}
