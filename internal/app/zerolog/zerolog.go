package zerolog

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"os"
)

var Module = fx.Module("logging",
	fx.Provide(
		NewZerolog,
		NewZerologWriter,
	),
)

func NewZerologWriter(nrapp *newrelic.Application) zerologWriter.ZerologWriter {
	return zerologWriter.New(os.Stdout, nrapp)
}

func NewZerolog(nrapp *newrelic.Application, writer zerologWriter.ZerologWriter) (zerolog.Logger, error) {
	logger := zerolog.New(writer)
	return logger, nil
}
