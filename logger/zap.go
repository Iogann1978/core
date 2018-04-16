package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(target zapcore.WriteSyncer, options ...zap.Option) Logger {
	l := new(ZapLogger)
	if target == nil {
		target = os.Stdout
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		target,
		zapcore.DebugLevel,
	)
	l.l = zap.New(core, options...)
	return l
}

func (l *ZapLogger) Debug(msg string, args map[string]interface{}) {
	l.l.Debug(msg, zap.Reflect(`args`, args))
}

func (l *ZapLogger) Info(msg string, args map[string]interface{}) {
	l.l.Info(msg, zap.Reflect(`args`, args))
}

func (l *ZapLogger) Fatal(msg string, args map[string]interface{}) {
	l.l.Fatal(msg, zap.Reflect(`args`, args))
}

func (l *ZapLogger) Panic(msg string, args map[string]interface{}) {
	l.l.Panic(msg, zap.Reflect(`args`, args))
}

func (l *ZapLogger) Warn(msg string, args map[string]interface{}) {
	l.l.Warn(msg, zap.Reflect(`args`, args))
}
