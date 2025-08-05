package logger

import (
	"context"
	"io"
	"time"

	"github.com/nexitf/logger/field"
	"go.uber.org/zap"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

type FieldOption = field.FieldOption

type LogHandler interface {
	Debug(message string, fields map[string]any, keys []string)
	Info(message string, fields map[string]any, keys []string)
	Warn(message string, fields map[string]any, keys []string)
	Error(message string, fields map[string]any, keys []string)
	Panic(message string, fields map[string]any, keys []string)
	Fatal(message string, fields map[string]any, keys []string)
	Flush() (err error)
}

type logInstance struct{}

var (
	Instance = logInstance{}
)

type InitOption func(*logger)

// WithCustomHandler
func WithCustomHandler(handler LogHandler) InitOption {
	return func(l *logger) { l.handler = handler }
}

// WithContextKeys
func WithContextKeys(keys ...string) InitOption {
	return func(l *logger) { l.ctxKeys = keys }
}

// WithZapLoggerName
func WithZapLoggerName(name string) InitOption {
	return func(l *logger) { l.zapCfg.name = name }
}

// WithZapLevel sets a level.
func WithZapLevel(level Level) InitOption {
	return func(l *logger) {
		switch level {
		case DebugLevel:
			l.zapCfg.level = zap.DebugLevel
		case InfoLevel:
			l.zapCfg.level = zap.InfoLevel
		case WarnLevel:
			l.zapCfg.level = zap.WarnLevel
		case ErrorLevel:
			l.zapCfg.level = zap.ErrorLevel
		case PanicLevel:
			l.zapCfg.level = zap.PanicLevel
		case FatalLevel:
			l.zapCfg.level = zap.FatalLevel
		}
	}
}

// WithZapWriter sets a log writer.
func WithZapWriter(w io.Writer) InitOption {
	return func(l *logger) {
		if w != nil {
			l.zapCfg.writer = w
		}
	}
}

// WithZapCaller sets caller in log.
func WithZapCaller(enabled bool) InitOption {
	return func(l *logger) {
		l.zapCfg.zloggerOpts = append(l.zapCfg.zloggerOpts, zap.WithCaller(enabled), zap.AddCallerSkip(3))
	}
}

// WithZapPanicHappened
func WithZapPanicHappened(fn func(time time.Time, logger string, message string, caller string)) InitOption {
	return func(l *logger) { l.zapCfg.panicFn = fn }
}

// WithZapFatalHappened
func WithZapFatalHappened(fn func(time time.Time, logger string, message string, caller string)) InitOption {
	return func(l *logger) { l.zapCfg.fatalFn = fn }
}

// WithZapEncoder sets a function to new an encoder.
func WithZapEncoder(fun Encoder) InitOption {
	return func(l *logger) {
		l.zapCfg.encoder = fun
	}
}

// WithZapEncoderFieldKey sets a key name with specific key.
func WithZapEncoderFieldKey(key, name string) InitOption {
	return func(l *logger) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderFieldKey(key, name))
	}
}

// WithZapEncoderRemoveField removes the specified field.
func WithZapEncoderRemoveField(key string) InitOption {
	return func(l *logger) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderFieldKey(key, ""))
	}
}

// WithZapEncoderTimeLayout sets a time encoder.
func WithZapEncoderTimeLayout(loc *time.Location, layout string) InitOption {
	return func(l *logger) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderTimeLayout(loc, layout))
	}
}

type logger struct {
	ctxKeys []string
	handler LogHandler
	zapCfg  zapConfig
}

var (
	log logger
)

// Init
func Init(opts ...InitOption) {
	log.zapCfg = defaultConfig
	// Set options
	for _, setOpt := range opts {
		setOpt(&log)
	}
	// Option: handler
	if log.handler == nil {
		log.handler = newZapHandler(&log.zapCfg)
	}
}

// Flush writes out all buffered logs.
func Flush() (err error) {
	if log.handler != nil {
		err = log.handler.Flush()
	}
	return
}

func (l *logger) Debug(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Debug(message, fieldsMap, keys)
	}
}

func Debug(message string, fields ...FieldOption) {
	log.Debug(message, fields...)
}

func (l *logger) DebugCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Debug(message, fieldsMap, keys)
	}
}

func DebugCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.DebugCtx(ctx, message, fields...)
}

func (l *logger) Info(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Info(message, fieldsMap, keys)
	}
}

func Info(message string, fields ...FieldOption) {
	log.Info(message, fields...)
}

func (l *logger) InfoCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Info(message, fieldsMap, keys)
	}
}

func InfoCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.InfoCtx(ctx, message, fields...)
}

func (l *logger) Warn(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Warn(message, fieldsMap, keys)
	}
}

func Warn(message string, fields ...FieldOption) {
	log.Warn(message, fields...)
}

func (l *logger) WarnCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Warn(message, fieldsMap, keys)
	}
}

func WarnCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.WarnCtx(ctx, message, fields...)
}

func (l *logger) Error(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Error(message, fieldsMap, keys)
	}
}

func Error(message string, fields ...FieldOption) {
	log.Error(message, fields...)
}

func (l *logger) ErrorCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Error(message, fieldsMap, keys)
	}
}

func ErrorCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.ErrorCtx(ctx, message, fields...)
}

func (l *logger) Panic(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Panic(message, fieldsMap, keys)
	}
}

func Panic(message string, fields ...FieldOption) {
	log.Panic(message, fields...)
}

func (l *logger) PanicCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Panic(message, fieldsMap, keys)
	}
}

func PanicCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.PanicCtx(ctx, message, fields...)
}

func (l *logger) Fatal(message string, fields ...FieldOption) {
	if h, ok := l.haslog(); ok {
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Fatal(message, fieldsMap, keys)
	}
}

func Fatal(message string, fields ...FieldOption) {
	log.Fatal(message, fields...)
}

func (l *logger) FatalCtx(ctx context.Context, message string, fields ...FieldOption) {
	if h, ok := l.haslogCtx(ctx); ok {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		h.Fatal(message, fieldsMap, keys)
	}
}

func FatalCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.FatalCtx(ctx, message, fields...)
}

// haslog
func (l *logger) haslog() (LogHandler, bool) {
	return l.handler, l.handler != nil
}

// haslogCtx
func (l *logger) haslogCtx(ctx context.Context) (LogHandler, bool) {
	log := ctx.Value(Instance)
	if log != nil {
		switch v := log.(type) {
		// LogHandler implemention
		case LogHandler:
			return v, true
		// zap logger
		case *zap.Logger:
			return &zapHandler{log: v}, true
		}
	}
	return l.haslog()
}

// toFieldsMap
func (l *logger) toFieldsMap(opts ...FieldOption) (fields map[string]any, keys []string) {
	fields = make(map[string]any)
	// Collect keys
	for _, fn := range opts {
		fn(func(key string, value any) {
			fields[key] = value
			keys = append(keys, key)
		})
	}
	return
}

// withCtxKeys
func (l *logger) withCtxKeys(ctx context.Context) (fields []FieldOption) {
	for _, key := range l.ctxKeys {
		if value := ctx.Value(key); value != nil {
			fields = append(fields, func(fn func(key string, value any)) { fn(key, value) })
		}
	}
	return
}
