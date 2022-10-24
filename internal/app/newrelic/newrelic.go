package newrelic

import (
	"context"
	"errors"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strings"
)

var Module = fx.Module("newrelic",
	fx.Provide(
		NewConfig,
		NewNewRelicApplication,
	),
)

type Config struct {
	ServiceName    string `env:"SERVICE_NAME,required"`
	ServiceVersion string `env:"SERVICE_VERSION,default=not-found"`

	// Useful for grouping
	Team string `env:"TEAM"`

	// e.g., "local"
	Env string `env:"ENV,required"`

	*NestedConfig `env:",prefix=NEW_RELIC_"`
}

type NestedConfig struct {
	Enabled bool   `env:"ENABLED,default=true"`
	Key     string `env:"KEY"`

	DebugLogging bool `env:"DEBUG_LOGGING"`

	// Whether the New Relic agent forwards logs to New Relic.
	// Only do this when logs are in the JSON format.
	ForwardLogs bool `env:"FORWARD_LOGS,default=true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewNewRelicApplication(cfg *Config, logger *zap.Logger) (*newrelic.Application, error) {
	enabled := cfg.Enabled
	if enabled && len(cfg.Key) == 0 {
		return nil, errors.New("missing key for New Relic")
	}

	if !enabled {
		logger.Warn("New Relic not enabled. Check Environment Variables for NEW_RELIC_ENABLED=true and NEW_RELIC_KEY=[valid-license-key]")
	}

	logger.Info("Attempting to initialize New Relic...")

	opts := buildNewRelicConfigOptions(cfg, logger)

	app, err := newrelic.NewApplication(opts...)

	if err != nil {
		logger.Error("Failed to create New Relic application.", zap.Error(err))
		return app, err
	}

	logger.Info("New Relic initialized.")

	return app, nil
}

func buildNewRelicConfigOptions(cfg *Config, logger *zap.Logger) []newrelic.ConfigOption {
	enabled := cfg.Enabled
	env := cfg.Env
	isLocal := strings.ToLower(env) == "local"
	debugLoggerEnabled := enabled && cfg.DebugLogging

	labels := map[string]string{
		"team":    cfg.Team,
		"service": cfg.ServiceName,
		"env":     env,
		"version": cfg.ServiceVersion,
	}

	opts := []newrelic.ConfigOption{
		// A base set of configs read from the environment.
		// Latter ConfigOptions may overwrite the Config fields already set.
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppName(cfg.ServiceName + "-" + env),
		newrelic.ConfigLicense(cfg.Key),
		newrelic.ConfigEnabled(enabled),
		newrelic.ConfigDistributedTracerEnabled(true),
		nrzap.ConfigLogger(logger.Named("newrelic")),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
			cfg.CrossApplicationTracer.Enabled = false // this is legacy and is now  DistributedTracerEnabled
			cfg.CustomInsightsEvents.Enabled = true
			cfg.TransactionTracer.Attributes.Enabled = true
			cfg.Labels = labels
		},
	}

	if cfg.ForwardLogs {
		opts = append(opts, newrelic.ConfigAppLogForwardingEnabled(true))
		logger.Info("New Relic - ConfigAppLogForwardingEnabled")
	}

	// if none local.. do log decoration of NR attributes
	if !isLocal {
		opts = append(opts, newrelic.ConfigAppLogDecoratingEnabled(true))
		logger.Info("New Relic - ConfigAppLogDecoratingEnabled")
	}

	if debugLoggerEnabled {
		opts = append(opts, newrelic.ConfigDebugLogger(os.Stdout))
		logger.Info("New Relic - ConfigDebugLogger")
	}

	return opts
}
