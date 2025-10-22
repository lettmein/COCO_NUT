package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() (*Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	sugar := zapLogger.Sugar()
	return &Logger{SugaredLogger: sugar}, nil
}

func (l *Logger) Close() error {
	return l.Sync()
}
