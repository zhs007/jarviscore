package jarviscore

import (
	"path"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
)

// Config - config
type Config struct {
	RunPath     string
	DefPeerAddr string
}

const normalLogFilename = "output.log"
const errLogFilename = "error.log"

var config = Config{RunPath: "./", DefPeerAddr: "jarvis.heyalgo.io:7788"}

// InitJarvisCore -
func InitJarvisCore(cfg Config) {
	config.RunPath = cfg.RunPath
	config.DefPeerAddr = cfg.DefPeerAddr

	log.InitLogger(getRealPath(normalLogFilename), getRealPath(errLogFilename))

	log.Info("InitJarvisCore", zap.String("RunPath", cfg.RunPath), zap.String("DefPeerAddr", cfg.DefPeerAddr))

	return
}

// ReleaseJarvisCore -
func ReleaseJarvisCore() error {
	log.ReleaseLogger()

	return nil
}

func getRealPath(filename string) string {
	return path.Join(config.RunPath, filename)
}
