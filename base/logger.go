package jarvisbase

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var onceLogger sync.Once

func initLogger(level zapcore.Level, isConsole bool, logpath string) (*zap.Logger, error) {
	loglevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	if isConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleDebugging := zapcore.Lock(os.Stdout)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleDebugging, loglevel),
		)

		cl := zap.New(core)
		// defer cl.Sync()

		return cl, nil
	}

	cfg := &zap.Config{}

	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.OutputPaths = []string{path.Join(logpath, "output.log")}
	cfg.ErrorOutputPaths = []string{path.Join(logpath, "error.log")}
	cfg.Encoding = "json"
	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:     "T",
		LevelKey:    "L",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}

	cl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// InitLogger - initializes a thread-safe singleton logger
func InitLogger(level zapcore.Level, isConsole bool, logpath string) {
	// once ensures the singleton is initialized only once
	onceLogger.Do(func() {
		cl, err := initLogger(level, isConsole, logpath)
		if err != nil {
			fmt.Printf("initLogger error! %v \n", err)

			os.Exit(-1)
		}

		logger = cl
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

// SyncLogger - sync logger
func SyncLogger() {
	logger.Sync()
}

// JSON - make json to field
func JSON(key string, obj interface{}) zap.Field {
	s, err := json.Marshal(obj)
	if err != nil {
		return zap.Error(err)
	}

	return zap.String(key, string(s))
}
