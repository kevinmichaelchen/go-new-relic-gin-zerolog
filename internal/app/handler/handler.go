package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/pkg/logging"
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Module("handler",
	fx.Provide(
		NewConfig,
		NewGinEngine,
	),
	fx.Invoke(
		RegisterHandler,
	),
)

type Config struct {
	Port int `env:"PORT,default=8081"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewGinEngine(
	lc fx.Lifecycle,
	cfg *Config,
	logger *zap.Logger,
	nrapp *newrelic.Application,
) *gin.Engine {

	// Create Gin router
	r := gin.Default()

	// Instrument requests with New Relic telemetry
	r.Use(nrgin.Middleware(nrapp))

	// Middleware to add the Trace ID to the logger
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()

		newRelicFields := getNrMetadataFields(ctx)

		newLogger := logger.With(newRelicFields...)

		newLogger.Info("Before injecting logger and calling Next()")

		// Update context
		newCtx := logging.ToContext(ctx, newLogger)

		//c.Request = c.Request.Clone(newCtx)
		c.Request = c.Request.WithContext(newCtx)

		// TODO the zap.Logger here has TraceContext, but it's not
		//  getting passed down correctly??

		// Invoke next handler in chain
		c.Next()
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.Port)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				// TODO get proper logger
				logger.Info("Serving HTTP", zap.String("addr", addr))
				err := r.Run(addr)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// TODO see https://github.com/gin-gonic/examples/tree/master/graceful-shutdown/graceful-shutdown
			return nil
		},
	})

	return r
}

func RegisterHandler(r *gin.Engine, nrapp *newrelic.Application) {
	// for k8s health check
	r.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()

		logger := logging.Extract(ctx)
		logger.Info("this is a zap log")
		c.Writer.Write([]byte("ok"))
	})
}

func getNrMetadataFields(ctx context.Context) []zap.Field {
	// Extract the New Relic transaction from context
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return []zap.Field{}
	}
	// The fields needed to link data to a trace
	md := txn.GetLinkingMetadata()
	return []zap.Field{
		zap.String(logcontext.KeyTraceID, md.TraceID),
		zap.String(logcontext.KeySpanID, md.SpanID),
		zap.String(logcontext.KeyEntityName, md.EntityName),
		zap.String(logcontext.KeyEntityType, md.EntityType),
		zap.String(logcontext.KeyEntityGUID, md.EntityGUID),
		zap.String(logcontext.KeyHostname, md.Hostname),
	}
}
