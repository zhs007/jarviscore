package jarviscore

import (
	"path"
	"sync"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore/log"
)

// Config - config
type Config struct {
	RunPath      string
	PeerAddrFile string
	DefPeerAddr  string
}

var config *Config
var onceConfig sync.Once

// InitJarvisCore -
func InitJarvisCore(cfg Config) (err error) {
	onceConfig.Do(func() {
		config = &Config{RunPath: cfg.RunPath, PeerAddrFile: cfg.PeerAddrFile, DefPeerAddr: cfg.DefPeerAddr}

		log.InitLogger(path.Join(config.RunPath, "output.log"), path.Join(config.RunPath, "error.log"))

		log.Info("InitJarvisCore", zap.String("RunPath", cfg.RunPath), zap.String("PeerAddrFile", cfg.PeerAddrFile), zap.String("DefPeerAddr", cfg.DefPeerAddr))
	})

	return
}

// ReleaseJarvisCore -
func ReleaseJarvisCore() error {
	log.ReleaseLogger()

	return nil
}
