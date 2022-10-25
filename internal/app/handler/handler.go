package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kevinmichaelchen/go-new-relic-gin-zap/pkg/logging"
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
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
	logger zerolog.Logger,
	nrapp *newrelic.Application,
	writer zerologWriter.ZerologWriter,
) *gin.Engine {

	// Create Gin router
	r := gin.Default()

	// Instrument requests with New Relic telemetry
	r.Use(nrgin.Middleware(nrapp))

	// Inject logger into context
	r.Use(func(c *gin.Context) {
		// Get the New Relic Transaction from the Gin context
		txn := nrgin.Transaction(c)

		// Always create a new logger in order to avoid changing the context of
		// the logger for other threads that may be logging external to this
		// transaction.
		newLogger := logger.Output(writer.WithTransaction(txn))

		md := txn.GetLinkingMetadata()
		sublogger := newLogger.With().
			Str(logcontext.KeyTraceID, md.TraceID).
			Str(logcontext.KeySpanID, md.SpanID).
			Str(logcontext.KeyEntityName, md.EntityName).
			Str(logcontext.KeyEntityType, md.EntityType).
			Str(logcontext.KeyEntityGUID, md.EntityGUID).
			Str(logcontext.KeyHostname, md.Hostname).
			Logger()

		//newCtx := sublogger.WithContext(c.Request.Context())
		newCtx := logging.ToContext(c.Request.Context(), sublogger)

		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.Port)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				err := r.Run(addr)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal().Err(err).Msg("server failed")
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
		logger.Info().Msg("this is a log")
		c.Writer.Write([]byte("ok"))
	})

	r.GET("/err", func(c *gin.Context) {
		ctx := c.Request.Context()

		logger := logging.Extract(ctx)

		c.Status(http.StatusInternalServerError)
		c.Writer.Write([]byte("mysterious error occurred"))

		txn := newrelic.FromContext(ctx)
		logger.
			Err(errors.New("something internal failed")).
			Bool("is_null", txn == nil).
			Msg("oops")
	})
}
