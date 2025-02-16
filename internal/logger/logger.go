package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var logger *zap.Logger

// Init инициализирует логгер.
func Init(isDev bool) {
	var err error

	if isDev {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %v", err))
	}
}

// Destroy уничтожает логгер, записывая оставшиеся данные.
func Destroy() {
	logger.Sync()
}

// Info пишет лог уровня INFO.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
