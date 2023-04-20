package loggerCdc

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var DebugLevel string

// NewLogger Create a new logger instance
func NewLogger(name string, isDevelopment bool) (Logger, error) {
	var l zapcore.Level
	switch strings.ToLower(DebugLevel) {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	default:
		l = zap.DebugLevel

	}

	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(l),
			Development:       false,
			DisableCaller:     true,
			DisableStacktrace: true,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "name",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stack",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.FullCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
			InitialFields:    nil,
		}
	}
	zapLogger, _ := config.Build()
	log := Logger{
		SugaredLogger: zapLogger.Sugar().With("source", name),
	}
	return log, nil
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Infof(format, v)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.Info(str)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.Debug(str)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	l.DPanicf(str)
}

func (l *Logger) NewWith(p1, p2 string) *Logger {
	return &Logger{SugaredLogger: l.SugaredLogger.With(p1, p2)}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.Info(string(p))
	return len(p), nil
}
