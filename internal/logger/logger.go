package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func ExitWithError(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)

	if err := logger.Sync(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
	}

	os.Exit(1)
}

func NewLogger(debug bool) (*zap.Logger, error) {
	if debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
