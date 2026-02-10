package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --- Интерфейс ---

// Logger — контракт логгера приложения.
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	With(fields ...any) Logger
	Sync() error
}

// --- Реализация на zap ---

// normalizeKVs превращает нечётное число аргументов в пары: одиночное значение получает ключ "value".
// Так можно писать log.Info("msg"), log.Info("msg", addr) или log.Info("msg", "key", val).
func normalizeKVs(kv []any) []any {
	if len(kv)%2 == 0 {
		return kv
	}
	return append(append([]any{}, kv[:len(kv)-1]...), "value", kv[len(kv)-1])
}

type zapLogger struct{ *zap.SugaredLogger }

func (z *zapLogger) Debug(msg string, keysAndValues ...any) {
	z.SugaredLogger.Debugw(msg, normalizeKVs(keysAndValues)...)
}
func (z *zapLogger) Info(msg string, keysAndValues ...any) {
	z.SugaredLogger.Infow(msg, normalizeKVs(keysAndValues)...)
}
func (z *zapLogger) Warn(msg string, keysAndValues ...any) {
	z.SugaredLogger.Warnw(msg, normalizeKVs(keysAndValues)...)
}
func (z *zapLogger) Error(msg string, keysAndValues ...any) {
	z.SugaredLogger.Errorw(msg, normalizeKVs(keysAndValues)...)
}

func (z *zapLogger) With(fields ...any) Logger {
	return &zapLogger{z.SugaredLogger.With(normalizeKVs(fields)...)}
}

func (z *zapLogger) Sync() error { return z.SugaredLogger.Sync() }

// --- Конструкторы ---

// NewLogger создаёт Logger по level и format.
// format == "console" → dev (читаемый вывод, stacktrace на Error), иначе → prod (JSON, без stacktrace).
func NewLogger(level, format string) (Logger, error) {
	if format == "console" {
		return newDev(level)
	}
	return newProd(level)
}

func newDev(level string) (Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = levelAt(level, zapcore.DebugLevel)
	cfg.DisableStacktrace = false
	// Дата и время без мс: день.месяц.год часы:минуты:секунды
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02.01.2006 15:04:05")
	zl, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{zl.Sugar()}, nil
}

func newProd(level string) (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = levelAt(level, zapcore.InfoLevel)
	cfg.DisableStacktrace = true
	zl, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{zl.Sugar()}, nil
}

func levelAt(s string, defaultLvl zapcore.Level) zap.AtomicLevel {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return zap.NewAtomicLevelAt(defaultLvl)
	}
	return zap.NewAtomicLevelAt(lvl)
}

// --- Глобальный default (для log.Info/… без инъекции) ---

var defaultLogger Logger

// SetDefault задаёт логгер для пакетных вызовов. Вызвать один раз при старте (в main после NewLogger).
// Если не вызвать — log.Info/Debug/Warn/Error будут noop.
func SetDefault(l Logger) {
	defaultLogger = l
}

// Sync сбрасывает буферы глобального логгера. Вызывать при выходе: defer log.Sync().
func Sync() {
	if defaultLogger != nil {
		_ = defaultLogger.Sync()
	}
}

type noop struct{}

func (noop) Debug(string, ...any) {}
func (noop) Info(string, ...any)  {}
func (noop) Warn(string, ...any)  {}
func (noop) Error(string, ...any) {}
func (noop) With(...any) Logger   { return noop{} }
func (noop) Sync() error          { return nil }

func getDefault() Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	return noop{}
}

// Пакетные функции — пишут в глобальный логгер (после SetDefault) или в noop.
func Debug(msg string, keysAndValues ...any) { getDefault().Debug(msg, keysAndValues...) }
func Info(msg string, keysAndValues ...any)  { getDefault().Info(msg, keysAndValues...) }
func Warn(msg string, keysAndValues ...any)  { getDefault().Warn(msg, keysAndValues...) }
func Error(msg string, keysAndValues ...any) { getDefault().Error(msg, keysAndValues...) }
func With(fields ...any) Logger              { return getDefault().With(fields...) }
