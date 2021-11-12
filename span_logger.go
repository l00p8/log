package log

import (
	"github.com/l00p8/tracer"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger     *zap.Logger
	span       trace.Span
	spanFields []zapcore.Field
}

func (sl spanLogger) Warn(msg string, fields ...zapcore.Field) {
	sl.logToSpan("warn", msg, fields...)
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Debug(msg string, fields ...zapcore.Field) {
	sl.logToSpan("debug", msg, fields...)
	sl.logger.Debug(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Info(msg string, fields ...zapcore.Field) {
	sl.logToSpan("info", msg, fields...)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Error(msg string, fields ...zapcore.Field) {
	sl.logToSpan("error", msg, fields...)
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	sl.logToSpan("fatal", msg, fields...)
	sl.logger.Fatal(msg, append(sl.spanFields, fields...)...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (sl spanLogger) With(fields ...zapcore.Field) Logger {
	return spanLogger{
		logger:     sl.logger.With(fields...),
		span:       sl.span,
		spanFields: sl.spanFields,
	}
}

func (sl spanLogger) logToSpan(level string, msg string, fields ...zapcore.Field) {
	// TODO rather than always converting the fields, we could wrap them into a lazy logger
	tags := map[string]string{
		"event": msg,
		"level": level,
	}
	for _, field := range fields {
		tags[field.Key] = field.String
	}
	tracer.AddSpanEvents(sl.span, msg, tags)
}
