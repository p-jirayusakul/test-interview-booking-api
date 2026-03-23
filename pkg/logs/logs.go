package logs

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(env string) (logger *zap.Logger, err error) {
	var cfg zap.Config

	if strings.ToUpper(env) == "LOCAL" || strings.ToUpper(env) == "DEV" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	return cfg.Build()
}

func ShouldSkip(path string) bool {
	switch {
	case strings.HasPrefix(path, "/health/liveness"):
		return true
	case strings.HasPrefix(path, "/health/readiness"):
		return true
	case strings.HasPrefix(path, "/swagger"):
		return true
	default:
		return false
	}
}
