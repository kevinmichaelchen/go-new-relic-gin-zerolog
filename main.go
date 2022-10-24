package main

import (
	"github.com/ipfans/fxlogger"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/internal/app"
	"go.uber.org/fx"
)

func main() {
	a := fx.New(
		app.Module,
		fx.WithLogger(
			fxlogger.Default(),
		),
	)
	a.Run()
}
