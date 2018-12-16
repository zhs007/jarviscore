package jarviscore

import (
	"io/ioutil"
	"os"

	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"

	"github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
)

// Config - config
type Config struct {
	RootServAddr string
	LstTrustNode []string

	// TimeRequestChild - RequestChild time
	//					- default 30s
	TimeRequestChild int64

	// MaxMsgLength - default 4mb
	MaxMsgLength int32

	AnkaDB struct {
		DBPath   string
		HTTPServ string
		Engine   string
	}

	Log struct {
		LogPath    string
		LogLevel   string
		LogConsole bool
	}

	BaseNodeInfo struct {
		NodeName string
		BindAddr string
		ServAddr string
	}
}

// const normalLogFilename = "output.log"
// const errLogFilename = "error.log"

// var config Config

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
func InitJarvisCore(cfg *Config) {
	// config = cfg

	jarvisbase.InitLogger(getLogLevel(cfg.Log.LogLevel), cfg.Log.LogConsole, cfg.Log.LogPath)

	jarvisbase.Info("InitJarvisCore",
		zap.String("DBPath", cfg.AnkaDB.DBPath),
		zap.String("RootServAddr", cfg.RootServAddr),
		zap.String("LogPath", cfg.Log.LogPath))

	return
}

// ReleaseJarvisCore -
func ReleaseJarvisCore() error {
	jarvisbase.SyncLogger()

	return nil
}

func checkConfig(cfg *Config) error {
	if cfg.TimeRequestChild <= 0 {
		cfg.TimeRequestChild = 180
	}

	if cfg.MaxMsgLength <= 0 {
		cfg.MaxMsgLength = 4194304
	}

	return nil
}

// func getRealPath(filename string) string {
// 	return path.Join(config.RunPath, filename)
// }

// LoadConfig - load config
func LoadConfig(filename string) (*Config, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fi.Close()
	fd, err1 := ioutil.ReadAll(fi)
	if err1 != nil {
		return nil, err1
	}

	cfg := &Config{}

	err2 := yaml.Unmarshal(fd, cfg)
	if err2 != nil {
		return nil, err2
	}

	err = checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
