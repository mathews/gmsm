package log

import (
	"go.uber.org/zap"
)

//Logger SugaredLogger
var Logger *zap.SugaredLogger

func init() {
	// logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	Logger = logger.Sugar()
}
