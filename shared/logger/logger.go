package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, data map[string]any)
	Warn(msg string, data map[string]any)
	Error(msg string, err error, data map[string]any)
}

type implLogger struct {
	log *slog.Logger
}

func New() Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &implLogger{
		log: slog.New(handler),
	}
}

func (l *implLogger) Info(msg string, data map[string]any) {
	l.log.Info(msg, slog.Any("data", data))
}

func (l *implLogger) Warn(msg string, data map[string]any) {
	l.log.Warn(msg, slog.Any("data", data))
}

func (l *implLogger) Error(msg string, err error, data map[string]any) {
	l.log.Error(
		msg,
		slog.Any("error", err),
		slog.Any("data", data),
	)
}
