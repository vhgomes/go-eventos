package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs(), "stdout"},
		Level:       zap.NewAtomicLevelAt(getLevelLogs()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "timestamp",
			EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	log, _ = logConfig.Build()
}

// TODO: logging saindo em um json na raiz da pasta mas será trocado para algo mais viavel.
func getOutputLogs() string {
	output := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_OUTPUT")))

	if output == "" {
		return "logs.json"
	}

	return output
}

func getLevelLogs() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
	log.Sync()
}

func Error(msg string, err error, fields ...zap.Field) {
	tags := append(fields, zap.NamedError("error", err))
	log.Error(msg, tags...)
	log.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
	log.Sync()
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
	log.Sync()
}

func Fatal(msg string, err error, fields ...zap.Field) {
	if err != nil {
		tags := append(fields, zap.NamedError("error", err))
		log.Fatal(msg, tags...)
	} else {
		log.Fatal(msg, fields...)
	}
	log.Sync()
}
