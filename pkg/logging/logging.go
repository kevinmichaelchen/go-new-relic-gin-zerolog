package logging

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	nullLogger   = zap.NewNop()
)

func Extract(ctx context.Context) otelzap.LoggerWithCtx {
	logger := extractLogger(ctx)
	return otelzap.New(logger).Ctx(ctx)
}

func extractLogger(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(ctxMarkerKey).(*zap.Logger)
	if !ok || l == nil {
		return nullLogger
	}
	return l
}

// ToContext adds the zap.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}
