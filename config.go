package jarviscore

import (
	"go.uber.org/zap/zapcore"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
)

// Config - config
type Config struct {
	DBPath         string
	LogPath        string
	DefPeerAddr    string
	AnkaDBHttpServ string
	AnkaDBEngine   string
	LogLevel       string
	LogConsole     bool
}

const normalLogFilename = "output.log"
const errLogFilename = "error.log"

var config = Config{
	DBPath:         "./data",
	LogPath:        "./log",
	DefPeerAddr:    "jarvis.heyalgo.io:7788",
	AnkaDBHttpServ: ":8888",
	AnkaDBEngine:   "leveldb",
}

func getLogLevel(str string) zapcore.Level {
	switch str {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

// InitJarvisCore -
func InitJarvisCore(cfg Config) {
	config = cfg

	jarvisbase.InitLogger(getLogLevel(cfg.LogLevel), cfg.LogConsole, cfg.LogPath)

	jarvisbase.Info("InitJarvisCore",
		zap.String("DBPath", cfg.DBPath),
		zap.String("DefPeerAddr", cfg.DefPeerAddr),
		zap.String("LogPath", cfg.LogPath))

	return
}

// ReleaseJarvisCore -
func ReleaseJarvisCore() error {
	jarvisbase.SyncLogger()

	return nil
}

// func getRealPath(filename string) string {
// 	return path.Join(config.RunPath, filename)
// }
