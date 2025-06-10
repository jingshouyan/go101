package logger

import (
	"go101/config"

	"github.com/timandy/routine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfg = config.Conf.Logger

func init() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	level := new(zapcore.Level)
	level.UnmarshalText([]byte(cfg.Level))
	core := zapcore.NewCore(encoder, writeSyncer, level)
	wrapCore := &WarpCore{core}
	log := zap.New(wrapCore, zap.AddCaller())
	zap.ReplaceGlobals(log)
}

type WarpCore struct {
	zapcore.Core
}

func (c *WarpCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *WarpCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	fields = append(fields, zap.Uint64("_goroutine", routine.Goid()))
	return c.Core.Write(ent, fields)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
