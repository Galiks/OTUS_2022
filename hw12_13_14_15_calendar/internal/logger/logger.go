package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *logger
	once sync.Once
)

type logger struct {
	zapLog *zap.Logger
}

func InitLog(level string, isPrintStackTrace bool, logPath string) error {
	var (
		err    error
		newLog *zap.Logger
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
		} else {
			writer = zapcore.AddSync(os.Stdout)
		}

		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writer, zapLevel),
		)
		var options []zap.Option
		if isPrintStackTrace {
			options = append(options, zap.AddStacktrace(zapLevel))
		}
		options = append(options, zap.AddCaller())
		newLog = zap.New(core, options...)

		log = &logger{
			zapLog: newLog,
		}
	})

	return err
}

func Info(msg ...any) {
	log.zapLog.Sugar().Info(msg)
}

func Error(msg ...any) {
	log.zapLog.Sugar().Error(msg)
}

func Fatal(msg ...any) {

	log.zapLog.Sugar().Fatal(msg)
}
