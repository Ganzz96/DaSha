package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		zap.NewAtomicLevelAt(zap.DebugLevel),
	))
	return &Logger{
		Logger: l,
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Sugar().Infof(format, v...)
}

func (l *Logger) Verbose() bool {
	return true
}
