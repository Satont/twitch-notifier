package logger

import (
	"log/slog"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
	Warn(input string, fields ...any)
	WithComponent(name string) Logger
	GetSlog() *slog.Logger
}
