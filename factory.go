package log

import (
	"context"
	"os"

	"github.com/l00p8/tracer"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

// NewLogger creates a new *zap.Logger.
func NewLogger(level string, fields ...zapcore.Field) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()
	err := atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	//encoderCfg.TimeKey = ""

	l := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).With(fields...)

	return l, nil
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger(b)
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	if span := tracer.SpanFromContext(ctx); span != nil {
		logger := spanLogger{span: span, logger: b.logger}
		logger.spanFields = []zapcore.Field{
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		}
		return logger
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
//func (b Factory) With(fields ...zapcore.Field) Factory {
//	return Factory{logger: b.logger.With(fields...)}
//}

// Info logs an info msg with fields
func (l Factory) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs an info msg with fields
func (l Factory) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

// Debug logs an info msg with fields
func (l Factory) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

// Error logs an error msg with fields
func (l Factory) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l Factory) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (l Factory) With(fields ...zapcore.Field) Logger {
	return logger{logger: l.logger.With(fields...)}
}
