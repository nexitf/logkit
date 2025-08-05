package logger

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultTimeLayout = "2006-01-02 15:04:05.000000 -0700 MST"
)

var (
	defaultConfig = zapConfig{
		level:   zapcore.DebugLevel,
		writer:  NewConsoleWriter(),
		encoder: NewConsoleEncoder,
		encoderOpts: []EncoderOption{
			withEncoderTimeLayout(time.Local, defaultTimeLayout),
			withEncoderFieldKey(EncoderFieldKeyFunction, ""), // Hide function field
		},
		zloggerOpts: []zap.Option{
			zap.WithCaller(true),
			zap.AddCallerSkip(3),
		},
	}
)

type zapHook struct {
	cfg   *zapConfig
	level zapcore.Level
}

// OnWrite implements zapcore.CheckWriteHook
func (zh *zapHook) OnWrite(ce *zapcore.CheckedEntry, _ []zap.Field) {
	if zh.level == zap.DebugLevel {
		if zh.cfg.panicFn != nil {
			zh.cfg.panicFn(ce.Time, ce.LoggerName, ce.Message, ce.Caller.TrimmedPath())
		}
	} else if zh.level == zap.FatalLevel {
		if zh.cfg.fatalFn != nil {
			zh.cfg.fatalFn(ce.Time, ce.LoggerName, ce.Message, ce.Caller.TrimmedPath())
		}
	}
}

type zapConfig struct {
	level       zapcore.Level
	name        string
	writer      io.Writer
	encoder     Encoder
	encoderOpts []EncoderOption
	zloggerOpts []zap.Option
	panicFn     func(time time.Time, logger string, message string, caller string)
	fatalFn     func(time time.Time, logger string, message string, caller string)
}

type zapHandler struct {
	log *zap.Logger
}

// newZapHandler
func newZapHandler(cfg *zapConfig) *zapHandler {
	newEnc := cfg.encoder
	// New zap logger
	zlog := zap.New(
		zapcore.NewCore(
			newEnc(cfg.encoderOpts...),
			zapcore.AddSync(newZapWriter(cfg.writer)),
			cfg.level,
		),
	)
	// Option: name
	if cfg.name != "" {
		zlog = zlog.Named(cfg.name)
	}
	// Options
	cfg.zloggerOpts = append(cfg.zloggerOpts,
		zap.WithPanicHook(&zapHook{cfg: cfg, level: zap.DebugLevel}),
		zap.WithFatalHook(&zapHook{cfg: cfg, level: zap.FatalLevel}),
	)
	zlog = zlog.WithOptions(cfg.zloggerOpts...)

	return &zapHandler{log: zlog}
}

// Debug implements LogHandler.
func (z *zapHandler) Debug(message string, fields map[string]any, keys []string) {
	z.log.Debug(message, z.toZapFields(fields, keys)...)
}

// Info implements LogHandler.
func (z *zapHandler) Info(message string, fields map[string]any, keys []string) {
	z.log.Info(message, z.toZapFields(fields, keys)...)
}

// Warn implements LogHandler.
func (z *zapHandler) Warn(message string, fields map[string]any, keys []string) {
	z.log.Warn(message, z.toZapFields(fields, keys)...)
}

// Error implements LogHandler.
func (z *zapHandler) Error(message string, fields map[string]any, keys []string) {
	z.log.Error(message, z.toZapFields(fields, keys)...)
}

// Panic implements LogHandler.
func (z *zapHandler) Panic(message string, fields map[string]any, keys []string) {
	z.log.Panic(message, z.toZapFields(fields, keys)...)
}

// Fatal implements LogHandler.
func (z *zapHandler) Fatal(message string, fields map[string]any, keys []string) {
	z.log.Fatal(message, z.toZapFields(fields, keys)...)
}

// Flush implements LogHandler.
func (z *zapHandler) Flush() (err error) {
	return z.log.Sync()
}

// toZapFields
func (z *zapHandler) toZapFields(fields map[string]any, keys []string) (zapFields []zap.Field) {
	for _, key := range keys {
		zapFields = append(zapFields, zap.Any(key, fields[key]))
	}
	return
}
