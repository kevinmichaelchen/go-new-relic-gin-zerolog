package logging

import (
	"context"
	"github.com/rs/zerolog"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	nullLogger   = zerolog.Nop()
)

func Extract(ctx context.Context) zerolog.Logger {
	l, ok := ctx.Value(ctxMarkerKey).(zerolog.Logger)
	//if !ok || l == nil {
	if !ok {
		return nullLogger
	}
	return l
}

// ToContext adds the zap.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}
