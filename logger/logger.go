package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomLogger struct {
	logger *zap.Logger
}

func NewCustomLogger(objectName string, isDevelopment bool) (*CustomLogger, error) {
	var config zap.Config

	if isDevelopment {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		//config.EncoderConfig.EncodeCaller = customCallerEncoder
	} else {
		config = zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:     "message",
				LevelKey:       "level",
				TimeKey:        "time",
				NameKey:        "objectName",
				CallerKey:      "caller",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}
		config.EncoderConfig.EncodeCaller = customCallerEncoder
	}

	logger, err := config.Build(zap.AddCaller(), zap.Fields(zap.String("invoker", objectName)))
	if err != nil {
		return nil, err
	}

	return &CustomLogger{
		logger: logger,
	}, nil
}

func (c *CustomLogger) logWithCaller(level zapcore.Level, msg string, tags ...zap.Field) {
	entry := c.logger.WithOptions(zap.AddCallerSkip(2)).With(tags...).Check(level, msg)
	if entry != nil {
		entry.Write()
	}
}

func (c *CustomLogger) Info(msg string, tags ...zap.Field) {
	c.logWithCaller(zapcore.InfoLevel, msg, tags...)
}

func (c *CustomLogger) Error(msg string, tags ...zap.Field) {
	c.logWithCaller(zapcore.ErrorLevel, msg, tags...)
}

func (c *CustomLogger) Warn(msg string, tags ...zap.Field) {
	c.logWithCaller(zapcore.WarnLevel, msg, tags...)
}

func (c *CustomLogger) Debug(msg string, tags ...zap.Field) {
	c.logWithCaller(zapcore.DebugLevel, msg, tags...)
}

func (c *CustomLogger) Fatal(msg string, tags ...zap.Field) {
	c.logWithCaller(zapcore.FatalLevel, msg, tags...)
}

func (c *CustomLogger) Sync() error {
	return c.logger.Sync()
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}
	enc.AppendString(fmt.Sprintf("%s.%s", caller.TrimmedPath(), caller.Function))
}
