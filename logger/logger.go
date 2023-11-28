package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(isDebug bool) *zap.SugaredLogger {

	developmentConfig := zap.NewDevelopmentConfig()
	if isDebug {
		developmentConfig.Level.SetLevel(zapcore.DebugLevel)
	} else {
		developmentConfig.Level.SetLevel(zapcore.InfoLevel)
	}
	logger := zap.Must(developmentConfig.Build())
	logger.Info("日志等级为:", zap.String("level", logger.Level().String()))
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}
