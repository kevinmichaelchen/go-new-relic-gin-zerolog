package app

import (
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/handler"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/newrelic"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Options(
	zerolog.Module,
	handler.Module,
	newrelic.Module,
)
