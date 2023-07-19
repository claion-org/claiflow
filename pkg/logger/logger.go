package logger

import "github.com/go-logr/logr"

var Logger func() logr.Logger = func() logr.Logger {
	return NewZapr(
		func(config *ZaprConfig) { config.Verbose = 0 },
		func(config *ZaprConfig) { config.DisableCaller = true },
		func(config *ZaprConfig) { config.DisableStacktrace = true },
	)
}

func Info(msg string, keysAndValues ...interface{}) {
	Logger().WithCallDepth(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	Logger().WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func V(level int) logr.Logger {
	return Logger().V(level)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return Logger().WithValues(keysAndValues...)
}

func WithName(name string) logr.Logger {
	return Logger().WithName(name)
}
