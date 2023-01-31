package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *Logger

func init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	Log = &Logger{logger.Sugar()}
}

type Logger struct {
	Sugar *zap.SugaredLogger
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.Sugar.Debugf(msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.Sugar.Infof(msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.Sugar.Warnf(msg, args...)
}

func (l *Logger) Errorf(err error, msg string, args ...interface{}) {
	l.Sugar.With("err", err).Errorf(msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.Sugar.Fatalf(msg, args...)
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{l.Sugar.With(args...)}
}
