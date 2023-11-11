package worker

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger
}

func NewLogger(zap *zap.Logger) *Logger {
	return &Logger{
		zapLogger: zap,
	}
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, v...)
}

func (logger *Logger) Print(level zapcore.Level, args ...interface{}) {
	logger.zapLogger.Sugar().Info(fmt.Sprint(args...))
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(zapcore.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Print(zapcore.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(zapcore.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Print(zapcore.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(zapcore.FatalLevel, args...)
}
