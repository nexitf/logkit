package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	EncoderFieldKeyMessage    = "message"
	EncoderFieldKeyLevel      = "level"
	EncoderFieldKeyTime       = "time"
	EncoderFieldKeyCaller     = "caller"
	EncoderFieldKeyFunction   = "function"
	EncoderFieldKeyStacktrace = "stacktrace"
)

var (
	defaultEncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

type Encoder func(...EncoderOption) zapcore.Encoder
type EncoderOption func(*zapcore.EncoderConfig)

// withEncoderFieldKey sets a new name for key.
func withEncoderFieldKey(key, name string) EncoderOption {
	return func(cfg *zapcore.EncoderConfig) {
		switch key {
		case EncoderFieldKeyTime:
			cfg.TimeKey = name
		case EncoderFieldKeyLevel:
			cfg.LevelKey = name
		case EncoderFieldKeyMessage:
			cfg.MessageKey = name
		case EncoderFieldKeyCaller:
			cfg.CallerKey = name
		case EncoderFieldKeyFunction:
			cfg.FunctionKey = name
		case EncoderFieldKeyStacktrace:
			cfg.StacktraceKey = name
		}
	}
}

// withEncoderTimeLayout sets a new time layout.
func withEncoderTimeLayout(loc *time.Location, layout string) EncoderOption {
	return func(cfg *zapcore.EncoderConfig) {
		cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(loc).Format(layout))
		}
	}
}

// NewJSONEncoder takes some options and then return an json encoder.
func NewJSONEncoder(opts ...EncoderOption) zapcore.Encoder {
	// Copy default encoder config
	cfg := defaultEncoderConfig
	// Set options
	for _, setOpt := range opts {
		setOpt(&cfg)
	}
	return zapcore.NewJSONEncoder(cfg)
}

// NewConsoleEncoder takes some options and then return an console encoder.
func NewConsoleEncoder(opts ...EncoderOption) zapcore.Encoder {
	// Copy default encoder config
	cfg := defaultEncoderConfig
	// Set options
	for _, setOpt := range opts {
		setOpt(&cfg)
	}
	return zapcore.NewConsoleEncoder(cfg)
}
