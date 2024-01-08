package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var std *zap.SugaredLogger

func New() (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = false
	config.DisableCaller = false

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	std = log.Sugar()
	return std, nil
}

func L(ctx context.Context) *zap.SugaredLogger {
	copy := *std
	lg := &copy

	if requestID := ctx.Value("requestId"); requestID != nil {
		lg = lg.With(zap.Any("requestId", requestID))
	}

	return lg
}
