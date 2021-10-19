package logger

import (
	"os"

	logger_interface "base/src/applications/interfaces/logger"

	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

func (l Logger) Info(message string, optional interface{}) {
	if optional == nil {
		l.zap.Info(message)
	} else {
		l.zap.Info(message, zap.Any("event", optional))
	}
}

func (l Logger) Debug(message string, optional interface{}) {
	if optional == nil {
		l.zap.Debug(message)
	} else {
		l.zap.Debug(message, zap.Any("event", optional))
	}
}

func (l Logger) Warn(message string, optional interface{}) {
	if optional == nil {
		l.zap.Warn(message)
	} else {
		l.zap.Warn(message, zap.Any("event", optional))
	}
}

func (l Logger) Error(message string, optional interface{}) {
	if optional == nil {
		l.zap.Error(message)
	} else {
		l.zap.Error(message, zap.Any("event", optional))
	}
}

func _setDevLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic("error trying to create new dev logger")
	}

	return logger
}

func _setProdLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		panic("error trying to create new prod logger")
	}

	return logger
}

func NewLogger() logger_interface.ILogger {

	var logger *zap.Logger

	switch os.Getenv("GO_ENV") {
	case "development":
		logger = _setDevLogger()
	case "production":
		logger = _setProdLogger()
	default:
		logger = _setDevLogger()
	}

	return &Logger{
		zap: logger,
	}
}
