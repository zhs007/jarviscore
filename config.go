package jarviscore

import (
	"path"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
)

// Config - config
type Config struct {
	DBPath         string
	LogPath        string
	DefPeerAddr    string
	AnkaDBHttpServ string
	AnkaDBEngine   string
}

const normalLogFilename = "output.log"
const errLogFilename = "error.log"

var config = Config{
	DBPath:         "./data",
	LogPath:        "./log",
	DefPeerAddr:    "jarvis.heyalgo.io:7788",
	AnkaDBHttpServ: "8888",
	AnkaDBEngine:   "leveldb",
}

// InitJarvisCore -
func InitJarvisCore(cfg Config) {
	config.DBPath = cfg.DBPath
	config.DefPeerAddr = cfg.DefPeerAddr

	log.InitLogger(path.Join(config.LogPath, normalLogFilename), path.Join(config.LogPath, errLogFilename))

	log.Info("InitJarvisCore",
		zap.String("DBPath", cfg.DBPath),
		zap.String("DefPeerAddr", cfg.DefPeerAddr),
		zap.String("LogPath", cfg.LogPath))

	return
}

// ReleaseJarvisCore -
func ReleaseJarvisCore() error {
	log.ReleaseLogger()

	return nil
}

// func getRealPath(filename string) string {
// 	return path.Join(config.RunPath, filename)
// }
