package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log           *logger
	once          sync.Once
	logStackTrace bool
)

type logger struct {
	log *zap.Logger
}

func InitLog(level string, isPrintStackTrace bool, logPath string) (*logger, error) {
	var (
		err error
		log *zap.Logger
	)
	once.Do(func() {
		var zapLevel zapcore.Level
		if isPrintStackTrace {
			switch level {
			case "INFO":
				zapLevel = zapcore.InfoLevel
			case "ERROR":
				zapLevel = zapcore.ErrorLevel
			default:
				zapLevel = zap.DebugLevel
			}
		}

		var writer zapcore.WriteSyncer = nil
		var logFile *os.File = nil
		if logPath != "" {
			logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return
			}
			writer = zapcore.AddSync(logFile)
		}

		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writer, zapLevel),
		)
		var options []zap.Option
		if isPrintStackTrace {
			options = append(options, zap.AddStacktrace(zapLevel))
		}
		options = append(options, zap.AddCaller())

		log = zap.New(core, options...)

		Log = &logger{
			log: log,
		}
	})

	return Log, err
}

func (l logger) Info(msg string) {
	l.log.Info(msg)
}

func (l logger) Error(msg string) {
	l.log.Error(msg)
}

func (l logger) Fatal(msg string) {
	l.log.Fatal(msg)
}
