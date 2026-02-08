package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger — интерфейс логгера для приложения.
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	With(fields ...any) Logger
	Sync() error
}

type zapLogger struct {
	*zap.SugaredLogger
}

func (z *zapLogger) Debug(msg string, keysAndValues ...any) {
	z.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (z *zapLogger) Info(msg string, keysAndValues ...any) {
	z.SugaredLogger.Infow(msg, keysAndValues...)
}

func (z *zapLogger) Warn(msg string, keysAndValues ...any) {
	z.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (z *zapLogger) Error(msg string, keysAndValues ...any) {
	z.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (z *zapLogger) With(fields ...any) Logger {
	return &zapLogger{z.SugaredLogger.With(fields...)}
}

func (z *zapLogger) Sync() error {
	return z.SugaredLogger.Sync()
}

// New создаёт Logger: format "console" → dev-режим, иначе → prod.
func NewLogger(level, format string) (Logger, error) {
	if format == "console" {
		return NewDev(level)
	}
	return NewProd(level)
}

// NewDev — режим разработки: консоль, читаемый вывод, stacktrace на Error, по умолчанию level=debug.
func NewDev(level string) (Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = parseLevel(level, zapcore.DebugLevel)
	cfg.DisableStacktrace = false
	zl, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{zl.Sugar()}, nil
}

// NewProd — продакшен: JSON, без stacktrace, по умолчанию level=info.
func NewProd(level string) (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = parseLevel(level, zapcore.InfoLevel)
	cfg.DisableStacktrace = true
	zl, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{zl.Sugar()}, nil
}

func parseLevel(s string, defaultLvl zapcore.Level) zap.AtomicLevel {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return zap.NewAtomicLevelAt(defaultLvl)
	}
	return zap.NewAtomicLevelAt(lvl)
}

// defaultLogger — глобальный логгер для пакетных вызовов (log.Info(...)).
var defaultLogger Logger

// SetDefault задаёт логгер по умолчанию. Вызвать один раз при старте (например в main после New).
func SetDefault(l Logger) {
	defaultLogger = l
}

// Sync сбрасывает буферы default-логгера. Вызвать перед выходом (defer log.Sync()).
func Sync() {
	if defaultLogger != nil {
		_ = defaultLogger.Sync()
	}
}

func defaultOrNoop() Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	return noopLogger{}
}

type noopLogger struct{}

func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (n noopLogger) With(...any) Logger { return n }
func (noopLogger) Sync() error          { return nil }

// Пакетные вызовы — используют default-логгер. Удобно: log.Info("started"), log.With("id", id).Info("request").
func Debug(msg string, keysAndValues ...any) { defaultOrNoop().Debug(msg, keysAndValues...) }
func Info(msg string, keysAndValues ...any)  { defaultOrNoop().Info(msg, keysAndValues...) }
func Warn(msg string, keysAndValues ...any)  { defaultOrNoop().Warn(msg, keysAndValues...) }
func Error(msg string, keysAndValues ...any) { defaultOrNoop().Error(msg, keysAndValues...) }
func With(fields ...any) Logger              { return defaultOrNoop().With(fields...) }
