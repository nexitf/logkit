package logkit

import (
	"context"
	"io"
	"time"

	"github.com/nexitf/logkit/errors"
	"go.uber.org/zap"
)

const (
	DefaultErrorFieldName = "error"
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

type FieldOption func(func(key string, value any))

// Field
func Field(key string, value any) FieldOption {
	return func(fn func(string, any)) { fn(key, value) }
}

type LogHandler interface {
	Debug(message string, fields map[string]any, keys []string)
	Info(message string, fields map[string]any, keys []string)
	Warn(message string, fields map[string]any, keys []string)
	Error(message string, fields map[string]any, keys []string)
	Panic(message string, fields map[string]any, keys []string)
	Fatal(message string, fields map[string]any, keys []string)
	Flush() (err error)
}

type InitOption func(*logkit)

// WithCustomHandler
func WithCustomHandler(handler LogHandler) InitOption {
	return func(l *logkit) { l.handler = handler }
}

// WithSeparateErrorField
func WithSeparateErrorField(name string) InitOption {
	return func(l *logkit) { l.errorName = name }
}

// WithContextKeys
func WithContextKeys(keys ...string) InitOption {
	return func(l *logkit) { l.ctxKeys = keys }
}

// WithNolog
func WithNolog() InitOption {
	return func(l *logkit) { l.nolog = true }
}

// WithZapLoggerName
func WithZapLoggerName(name string) InitOption {
	return func(l *logkit) { l.zapCfg.name = name }
}

// WithZapLevel sets a level.
func WithZapLevel(level Level) InitOption {
	return func(l *logkit) {
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
	return func(l *logkit) {
		if w != nil {
			l.zapCfg.writer = w
		}
	}
}

// WithZapCaller sets caller in log.
func WithZapCaller(enabled bool) InitOption {
	return func(l *logkit) {
		l.zapCfg.zloggerOpts = append(l.zapCfg.zloggerOpts, zap.WithCaller(enabled), zap.AddCallerSkip(3))
	}
}

// WithZapPanicHappened
func WithZapPanicHappened(fn func(time time.Time, logger string, message string, caller string)) InitOption {
	return func(l *logkit) { l.zapCfg.panicFn = fn }
}

// WithZapFatalHappened
func WithZapFatalHappened(fn func(time time.Time, logger string, message string, caller string)) InitOption {
	return func(l *logkit) { l.zapCfg.fatalFn = fn }
}

// WithZapEncoder sets a function to new an encoder.
func WithZapEncoder(fun Encoder) InitOption {
	return func(l *logkit) {
		l.zapCfg.encoder = fun
	}
}

// WithZapEncoderFieldKey sets a key name with specific key.
func WithZapEncoderFieldKey(key, name string) InitOption {
	return func(l *logkit) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderFieldKey(key, name))
	}
}

// WithZapEncoderRemoveField removes the specified field.
func WithZapEncoderRemoveField(key string) InitOption {
	return func(l *logkit) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderFieldKey(key, ""))
	}
}

// WithZapEncoderTimeLayout sets a time encoder.
func WithZapEncoderTimeLayout(loc *time.Location, layout string) InitOption {
	return func(l *logkit) {
		l.zapCfg.encoderOpts = append(l.zapCfg.encoderOpts, withEncoderTimeLayout(loc, layout))
	}
}

type logkit struct {
	handler   LogHandler
	errorName string
	nolog     bool
	ctxKeys   []string
	zapCfg    zapConfig
}

var (
	log logkit
)

// Init
func Init(opts ...InitOption) {
	log.zapCfg = defaultConfig
	// Set options
	for _, setOpt := range opts {
		setOpt(&log)
	}
	// Option: handler
	if log.handler == nil && !log.nolog {
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

func (l *logkit) Debug(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Debug(message, fieldsMap, keys)
	}
}

func Debug(message string, fields ...FieldOption) {
	log.Debug(message, fields...)
}

func (l *logkit) DebugCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Debug(message, fieldsMap, keys)
	}
}

func DebugCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.DebugCtx(ctx, message, fields...)
}

func (l *logkit) Info(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Info(message, fieldsMap, keys)
	}
}

func Info(message string, fields ...FieldOption) {
	log.Info(message, fields...)
}

func (l *logkit) InfoCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Info(message, fieldsMap, keys)
	}
}

func InfoCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.InfoCtx(ctx, message, fields...)
}

func (l *logkit) Warn(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Warn(message, fieldsMap, keys)
	}
}

func Warn(message string, fields ...FieldOption) {
	log.Warn(message, fields...)
}

func (l *logkit) WarnCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Warn(message, fieldsMap, keys)
	}
}

func WarnCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.WarnCtx(ctx, message, fields...)
}

func (l *logkit) Error(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Error(message, fieldsMap, keys)
	}
}

func Error(message string, fields ...FieldOption) {
	log.Error(message, fields...)
}

func (l *logkit) ErrorWrap(err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Error(message, fieldsMap, keys)
	}
}

func ErrorWrap(err error, message string, fields ...FieldOption) {
	log.ErrorWrap(err, message, fields...)
}

func (l *logkit) ErrorCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Error(message, fieldsMap, keys)
	}
}

func ErrorCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.ErrorCtx(ctx, message, fields...)
}

func (l *logkit) ErrorWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Error(message, fieldsMap, keys)
	}
}

func ErrorWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	log.ErrorWrapCtx(ctx, err, message, fields...)
}

func (l *logkit) Panic(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Panic(message, fieldsMap, keys)
	}
}

func Panic(message string, fields ...FieldOption) {
	log.Panic(message, fields...)
}

func (l *logkit) PanicWrap(err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Panic(message, fieldsMap, keys)
	}
}

func PanicWrap(err error, message string, fields ...FieldOption) {
	log.PanicWrap(err, message, fields...)
}

func (l *logkit) PanicCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Panic(message, fieldsMap, keys)
	}
}

func PanicCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.PanicCtx(ctx, message, fields...)
}

func (l *logkit) PanicWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Panic(message, fieldsMap, keys)
	}
}

func PanicWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	log.PanicWrapCtx(ctx, err, message, fields...)
}

func (l *logkit) Fatal(message string, fields ...FieldOption) {
	if l.haslog() {
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Fatal(message, fieldsMap, keys)
	}
}

func Fatal(message string, fields ...FieldOption) {
	log.Fatal(message, fields...)
}

func (l *logkit) FatalWrap(err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Fatal(message, fieldsMap, keys)
	}
}

func FatalWrap(err error, message string, fields ...FieldOption) {
	log.FatalWrap(err, message, fields...)
}

func (l *logkit) FatalCtx(ctx context.Context, message string, fields ...FieldOption) {
	if l.haslog() {
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Fatal(message, fieldsMap, keys)
	}
}

func FatalCtx(ctx context.Context, message string, fields ...FieldOption) {
	log.FatalCtx(ctx, message, fields...)
}

func (l *logkit) FatalWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	if l.haslog() {
		if l.errorName != "" {
			fields = append(fields, Field(l.errorName, err))
		} else {
			message = errors.Wrap(err, message).Error()
		}
		fields = append(fields, l.withCtxKeys(ctx)...)
		fieldsMap, keys := l.toFieldsMap(fields...)
		l.handler.Fatal(message, fieldsMap, keys)
	}
}

func FatalWrapCtx(ctx context.Context, err error, message string, fields ...FieldOption) {
	log.FatalWrapCtx(ctx, err, message, fields...)
}

// haslog
func (l *logkit) haslog() bool {
	return l.handler != nil
}

// toFieldsMap
func (l *logkit) toFieldsMap(opts ...FieldOption) (fields map[string]any, keys []string) {
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
func (l *logkit) withCtxKeys(ctx context.Context) (fields []FieldOption) {
	for _, key := range l.ctxKeys {
		if value := ctx.Value(key); value != nil {
			fields = append(fields, func(fn func(key string, value any)) { fn(key, value) })
		}
	}
	return
}
