package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZaprConfig struct {
	Verbose           int  `yaml:"verbose"` // enum(-1: only error, 0: info, 1: debug, 2: trace, n>2: further more)
	DisableCaller     bool `yaml:"disableCaller"`
	DisableStacktrace bool `yaml:"disableStacktrace"`
}

type OpZaprConfig func(cfg *ZaprConfig)

func NewZapr(opts ...OpZaprConfig) logr.Logger {
	var cfg = ZaprConfig{
		Verbose:           0,
		DisableCaller:     true,
		DisableStacktrace: true,
	}

	for _, op := range opts {
		op(&cfg)
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Verbose * -1))
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zapConfig.DisableCaller = cfg.DisableCaller
	zapConfig.DisableStacktrace = cfg.DisableStacktrace

	zapLog, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Errorf("%w: failed to build zap config", err))
	}

	return zapr.NewLogger(zapLog)

}
