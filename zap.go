package kirkzap

import (
	"context"

	"github.com/pkg/errors"

	"go.uber.org/zap"

	"github.com/jimenezmaximiliano/kirk"
)

// Adapter adapts a zap SugaredLogger to implement kirk.Logger.
type Adapter struct {
	zap *zap.SugaredLogger
}

// NewLoggerFromSugaredZap returns an adapted zap SugaredLogger that implements Logger.
func NewLoggerFromSugaredZap(zap *zap.SugaredLogger) Adapter {
	return Adapter{zap: zap}
}

// NewLoggerFromDefaultZap returns an adapted zap SugaredLogger that implements Logger.
func NewLoggerFromDefaultZap() (Adapter, error) {
	baseLogger, err := zap.NewProduction()
	if err != nil {
		return Adapter{}, errors.Wrap(err, "failed to create a base zap logger")
	}

	sugaredLogger := baseLogger.Sugar()
	if sugaredLogger == nil {
		return Adapter{}, errors.Wrap(err, "failed to create a zap sugared logger")
	}

	return NewLoggerFromSugaredZap(sugaredLogger), nil
}

// Make sure Adapter implements kirk.Logger.
var _ kirk.Logger = Adapter{}

func (zap Adapter) withFields(ctx context.Context) *zap.SugaredLogger {
	if fields := kirk.FieldsFromCtx(ctx); len(fields) > 0 {

		return zap.zap.With("fields", fields)
	}

	return zap.zap
}

// Error logs an error with ERROR severity.
func (zap Adapter) Error(ctx context.Context, err error) {
	zap.withFields(ctx).Error(err)
}

// Panic logs an error with PANIC severity, and panics.
func (zap Adapter) Panic(ctx context.Context, err error) {
	zap.withFields(ctx).Panic(err)
}

// Debug logs a message with DEBUG severity.
func (zap Adapter) Debug(ctx context.Context, message string) {
	zap.withFields(ctx).Debug(message)
}

// Info logs a message with INFO severity.
func (zap Adapter) Info(ctx context.Context, message string) {
	zap.withFields(ctx).Info(message)
}

// Warn logs a message with WARN severity.
func (zap Adapter) Warn(ctx context.Context, message string) {
	zap.withFields(ctx).Warn(message)
}
