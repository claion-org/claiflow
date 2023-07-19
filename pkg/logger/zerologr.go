package logger

import (
	"io"
	"os"
	"path"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type ZerologConfig struct {
	Verbose       int       `env:"LOG_VERBOSE"        yaml:"verbose"` // enum(-1: only error, 0: info, 1: debug, 2: trace, n>2: further more)
	DisableCaller bool      `env:"LOG_CALLER_DISABLE" yaml:"callerDisable"`
	Writer        io.Writer `json:"_"`
}

type OpZerologConfig func(cfg *ZerologConfig)

func NewZerologr(opts ...OpZerologConfig) logr.Logger {
	var cfg = ZerologConfig{
		Verbose:       0,
		DisableCaller: true,
		Writer:        os.Stdout,
	}
	for _, op := range opts {
		op(&cfg)
	}

	// zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999999"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	// zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zl := zerolog.New(cfg.Writer)
	zl = zl.With().Timestamp().Logger()

	if cfg.DisableCaller {
		zl = zl.With().Stack().Logger()
	}

	// zev := zl.WithLevel(zeroLevel)

	// zerologr.NameSeparator = "."
	zerologr.SetMaxV(cfg.Verbose)
	logger := zerologr.New(&zl)
	logger = logger.WithCallDepth(0)
	logger = logger.WithName(path.Base(os.Args[0]))

	return logger
}
