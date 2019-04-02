package jarviscore

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("./cfg/config.yaml")
	if err != nil {
		t.Fatalf("TestLoadConfig LoadConfig err %v", err)

		return
	}

	if cfg.RootServAddr != "jarvis.heyalgo.io:7788" {
		t.Fatalf("TestLoadConfig RootServAddr %v", cfg.RootServAddr)
	}

	if len(cfg.LstTrustNode) != 1 {
		t.Fatalf("TestLoadConfig len(cfg.LstTrustNode) %v", len(cfg.LstTrustNode))
	}

	if cfg.LstTrustNode[0] != "1JJaKpZGhYPuVHc1EKiiHZEswPAB5SybW5" {
		t.Fatalf("TestLoadConfig LstTrustNode %v", cfg.LstTrustNode[0])
	}

	if cfg.TimeRequestChild != 180 {
		t.Fatalf("TestLoadConfig TimeRequestChild %v", cfg.TimeRequestChild)
	}

	if cfg.MaxMsgLength != 4194304 {
		t.Fatalf("TestLoadConfig MaxMsgLength %v", cfg.MaxMsgLength)
	}

	if cfg.AnkaDB.DBPath != "./dat" {
		t.Fatalf("TestLoadConfig AnkaDB.DBPath %v", cfg.AnkaDB.DBPath)
	}

	if cfg.AnkaDB.HTTPServ != ":8880" {
		t.Fatalf("TestLoadConfig AnkaDB.HTTPServ %v", cfg.AnkaDB.HTTPServ)
	}

	if cfg.AnkaDB.Engine != "leveldb" {
		t.Fatalf("TestLoadConfig AnkaDB.Engine %v", cfg.AnkaDB.Engine)
	}

	if cfg.Log.LogPath != "./logs" {
		t.Fatalf("TestLoadConfig Log.LogPath %v", cfg.Log.LogPath)
	}

	if cfg.Log.LogLevel != "debug" {
		t.Fatalf("TestLoadConfig Log.LogLevel %v", cfg.Log.LogLevel)
	}

	if cfg.Log.LogConsole != true {
		t.Fatalf("TestLoadConfig Log.LogConsole %v", cfg.Log.LogConsole)
	}

	if cfg.BaseNodeInfo.NodeName != "node001" {
		t.Fatalf("TestLoadConfig BaseNodeInfo.NodeName %v", cfg.BaseNodeInfo.NodeName)
	}

	if cfg.BaseNodeInfo.BindAddr != ":7788" {
		t.Fatalf("TestLoadConfig BaseNodeInfo.BindAddr %v", cfg.BaseNodeInfo.BindAddr)
	}

	if cfg.BaseNodeInfo.ServAddr != "127.0.0.1:7788" {
		t.Fatalf("TestLoadConfig BaseNodeInfo.ServAddr %v", cfg.BaseNodeInfo.ServAddr)
	}

	if cfg.AutoUpdate != true {
		t.Fatalf("TestLoadConfig AutoUpdate %v", cfg.AutoUpdate)
	}

	// 	if cfg.UpdateScript != `cd ..
	// rm -rf jarvissh.tar.gz
	// wget https://github.com/zhs007/jarvissh/releases/download/{{.NewVersion}}/jarvissh.tar.gz
	// tar zxvf jarvissh.tar.gz
	// cd jarvissh
	// ./jarvissh stop
	// ./jarvissh start -d` {
	// 		t.Fatalf("TestLoadConfig UpdateScript %v", cfg.UpdateScript)
	// 	}

	t.Logf("TestLoadConfig is OK")
}
