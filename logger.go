// Logger это набор для упрощенного подключения логирования к проекту или пакету
// Его использование ведет к уменьшению повторения boilerplate кода в проектах

// Пакет в своей основе использует Zap - невероятно быстрый, структурированный,
// с поддержкой уровней логгер, разработанный в Uber: https://github.com/uber-go/zap

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	Warn(msg string, fields ...zapcore.Field)
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
}

// logger delegates all calls to the underlying zap.Logger
type logger struct {
	logger *zap.Logger
}

// Info logs an info msg with fields
func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs an info msg with fields
func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// Debug logs an info msg with fields
func (l logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Error logs an error msg with fields
func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l logger) With(fields ...zapcore.Field) Logger {
	return logger{logger: l.logger.With(fields...)}
}
