package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.elastic.co/apm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger LoggerConf

type LoggerConf struct {
	// Logger instance
	dep *zap.Logger
}

type Logger interface {
	Info(ctx context.Context, msg string, meta interface{})
	Error(ctx context.Context, msg string, meta interface{})
}

func SetupLogger(serviceName string) *LoggerConf {
	lg := zap.NewProductionConfig()
	lg.DisableStacktrace = false
	lg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / int64(time.Millisecond))
		},
		EncodeName: zapcore.FullNameEncoder,
	})

	lg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	lg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	lg.EncoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(int64(d) / int64(time.Millisecond))
	}
	lg.EncoderConfig.EncodeName = zapcore.FullNameEncoder

	fileEncoder := zapcore.NewJSONEncoder(lg.EncoderConfig)
	fileName := fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02"))
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     14, // days
	})
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(fileEncoder, writer, zap.DebugLevel),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zapLogger = zapLogger.With(ZapString("serviceName", serviceName))
	return &LoggerConf{
		dep: zapLogger,
	}
}

func Init(l *LoggerConf) {
	logger = LoggerConf{
		dep: l.dep,
	}
}

func GetLogger() Logger {
	return &logger
}

func (l *LoggerConf) Info(ctx context.Context, msg string, meta interface{}) {
	metaField := zap.Any("meta", meta)
	l.withTraceInfo(ctx).dep.Info(msg, metaField)
}

func (l *LoggerConf) Error(ctx context.Context, msg string, meta interface{}) {
	metaField := zap.Any("meta", meta)
	l.withTraceInfo(ctx).dep.Error(msg, metaField)
}

func (l *LoggerConf) withTraceInfo(ctx context.Context) *LoggerConf {
	span := apm.SpanFromContext(ctx)
	if span == nil {
		return l.Clone(l.dep)
	}
	apmCtx := span.TraceContext()
	traceId := ZapString("trace.id", fmt.Sprintf("%s", apmCtx.Trace.String()))
	spanId := ZapString("span.id", fmt.Sprintf("%s", apmCtx.Span.String()))
	return l.Clone(l.dep.With(
		traceId,
		spanId,
	))
}

func ZapString(key, value string) zap.Field {
	if value == "" {
		return zap.Skip()
	}
	return zap.String(key, value)
}

// Clone will create new Logger instance with specified zap.Logger.
func (l *LoggerConf) Clone(logger *zap.Logger) *LoggerConf {
	return &LoggerConf{
		dep: logger,
	}
}
