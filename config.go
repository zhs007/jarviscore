package jarviscore

import (
	"io/ioutil"
	"os"

	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"

	jarvisbase "github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
)

// Config - config
type Config struct {

	//------------------------------------------------------------------
	// base configuration

	RootServAddr string
	LstTrustNode []string

	// TimeRequestChild - RequestChild time
	//					- default 30s
	TimeRequestChild int64

	// MaxMsgLength - default 4mb
	MaxMsgLength int32

	//------------------------------------------------------------------
	// pprof configuration

	Pprof struct {
		BaseURL string
	}

	//------------------------------------------------------------------
	// task server configuration

	TaskServ struct {
		BindAddr string
		ServAddr string
	}

	//------------------------------------------------------------------
	// http server configuration

	HTTPServ struct {
		BindAddr string
		ServAddr string
	}

	//------------------------------------------------------------------
	// ankadb configuration

	AnkaDB struct {
		DBPath   string
		HTTPServ string
		Engine   string
	}

	//------------------------------------------------------------------
	// logger configuration

	Log struct {
		LogPath        string
		LogLevel       string
		LogConsole     bool
		LogSubFileName string
	}

	//------------------------------------------------------------------
	// basenodeinfo configuration

	BaseNodeInfo struct {
		NodeName string
		BindAddr string
		ServAddr string
	}

	//------------------------------------------------------------------
	// auto update

	AutoUpdate    bool
	UpdateScript  string
	RestartScript string
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
func InitJarvisCore(cfg *Config, nodeType string, version string) {
	// config = cfg

	cfg.Log.LogSubFileName = jarvisbase.BuildLogSubFilename(nodeType, version)

	jarvisbase.InitLogger(getLogLevel(cfg.Log.LogLevel), cfg.Log.LogConsole,
		cfg.Log.LogPath, cfg.Log.LogSubFileName)

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

	if cfg.AutoUpdate {
		if cfg.UpdateScript == "" || cfg.RestartScript == "" {
			return ErrCfgInvalidUpdateScript
		}
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
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(fd, cfg)
	if err != nil {
		return nil, err
	}

	err = checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
