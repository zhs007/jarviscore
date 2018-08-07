package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var onceLogger sync.Once

func initLogger(logpath string, errlogpath string) (*zap.Logger, error) {
	var cfgZap zap.Config
	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	cfgZap.Encoding = "json"
	cfgZap.OutputPaths = []string{logpath}
	cfgZap.ErrorOutputPaths = []string{errlogpath}
	cfgZap.Level = atom
	cfgZap.EncoderConfig = zap.NewProductionEncoderConfig()
	cfgZap.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfgZap.Build()
}

// ReleaseLogger -
func ReleaseLogger() error {
	if logger != nil {
		logger.Sync()
	}

	return nil
}

// InitLogger - initializes a thread-safe singleton logger
func InitLogger(logpath string, errlogpath string) (err error) {
	// once ensures the singleton is initialized only once
	onceLogger.Do(func() {
		logger, err = initLogger(logpath, errlogpath)
	})

	return
}

// // Log a message at the given level with given fields
// func Log(level zap.Level, message string, fields ...zap.Field) {
// 	singleton.Log(level, message, fields...)
// }

// Debug logs a debug message with the given fields
func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

// Info logs a debug message with the given fields
func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

// Error logs a debug message with the given fields
func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
