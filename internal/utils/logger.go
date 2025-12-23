package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

)

var Log *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}
}