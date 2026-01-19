package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.SugaredLogger
}

func NewLogger(logPath string) (*Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		zap: logger.Sugar(),
	}, nil
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.zap.Infow(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.zap.Errorw(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.zap.Warnw(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.zap.Debugw(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}
