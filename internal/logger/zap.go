package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wx-gw/config"
)

const TimeFormat = "2006-01-02 15:04:05"

type Logger struct {
	Zap        *zap.SugaredLogger
	DesugarZap *zap.Logger
}

func NewZapLogger(conf *config.Config) (*zap.SugaredLogger, *zap.Logger) {
	var options []zap.Option
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     localTimeEncoder,
	}

	if conf.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level := zap.NewAtomicLevelAt(toLevel(conf.Log.Level))

	core := zapcore.NewCore(encoder, toWriter(conf), level)

	stackLevel := zap.NewAtomicLevel()
	stackLevel.SetLevel(zap.WarnLevel)
	options = append(options,
		//zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(stackLevel),
	)

	logger := zap.New(core, options...)
	return logger.Sugar(), logger
}

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TimeFormat))
}

func toLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func toWriter(conf *config.Config) zapcore.WriteSyncer {
	fp := ""
	sp := string(filepath.Separator)

	fp, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
	fp += sp + "logs" + sp

	if conf.Log.Directory != "" {
		fp = conf.Log.Directory
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(fp, conf.Name) + ".log",
			MaxSize:    100,
			MaxAge:     0,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   true,
		}),
	)
}
