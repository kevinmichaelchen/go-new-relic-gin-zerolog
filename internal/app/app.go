package app

import (
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/handler"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/loggingzerolog"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app/newrelic"
	"go.uber.org/fx"
)

var Module = fx.Options(
	//loggingzap.Module,
	loggingzerolog.Module,
	handler.Module,
	newrelic.Module,
)
