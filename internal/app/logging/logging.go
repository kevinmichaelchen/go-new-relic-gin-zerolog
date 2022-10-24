package logging

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("logging",
	fx.Provide(
		NewLogger,
	),
)

func NewLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}
