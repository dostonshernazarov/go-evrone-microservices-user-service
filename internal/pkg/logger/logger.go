package logger

import (
	"github/user_service_evrone_microservces/internal/pkg/app"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func productionConfig(file string) zap.Config {
	configZap := zap.NewProductionConfig()
	configZap.OutputPaths = []string{"stdout", file}
	configZap.DisableStacktrace = true
	return configZap
}

func developmentConfig(file string) zap.Config {
	configZap := zap.NewDevelopmentConfig()
	configZap.OutputPaths = []string{"stdout", file}
	configZap.ErrorOutputPaths = []string{"stderr"}
	return configZap
}

func New(level, environment string, file_name string) (*zap.Logger, error) {
	file := filepath.Join("./" + file_name)

	configZap := productionConfig(file)

	if environment == app.EnvironmentDevelop {
		configZap = developmentConfig(file)
	}

	switch level {
	case "debug":
		configZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		configZap.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		configZap.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		configZap.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "dpanic":
		configZap.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "panic":
		configZap.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "fatal":
		configZap.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		configZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return configZap.Build()
}

func Error(err error) zapcore.Field {
	return zap.Error(err)
}
