package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(logFile *os.File) context.Context {
	consoleWriter := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	logger := zerolog.New(multi).With().Timestamp().Logger()

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	return ctx
}
