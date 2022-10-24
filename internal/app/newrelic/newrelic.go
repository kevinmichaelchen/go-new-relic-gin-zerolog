package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"os"
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

	// Toggles whether the agent enriches local logs printed to console, so they
	// can be sent to new relic for ingestion
	DecorateLogs bool `env:"DECORATE_LOGS,default=true"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewNewRelicApplication(cfg *Config) (*newrelic.Application, error) {
	enabled := cfg.Enabled
	if enabled && len(cfg.Key) == 0 {
		return nil, errors.New("missing key for New Relic")
	}

	if !enabled {
		fmt.Println("[WARN] New Relic monitoring is not enabled.")
	}

	opts := buildNewRelicConfigOptions(cfg)

	app, err := newrelic.NewApplication(opts...)

	if err != nil {
		return app, err
	}

	return app, nil
}

func buildNewRelicConfigOptions(cfg *Config) []newrelic.ConfigOption {
	enabled := cfg.Enabled
	env := cfg.Env
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
		newrelic.ConfigInfoLogger(os.Stdout),
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
	}

	// if none local.. do log decoration of NR attributes
	if cfg.DecorateLogs {
		opts = append(opts, newrelic.ConfigAppLogDecoratingEnabled(true))
	}

	if debugLoggerEnabled {
		opts = append(opts, newrelic.ConfigDebugLogger(os.Stdout))
	}

	return opts
}
