package coredb

import (
	"testing"

	"go.uber.org/zap/zapcore"

	jarvisbase "github.com/zhs007/jarviscore/base"
)

func TestBackup06(t *testing.T) {
	jarvisbase.InitLogger(zapcore.DebugLevel, true, "", "")

	//------------------------------------------------------------------------
	// initial CoreDB

	cdb, err := NewCoreDB("../test/backup-v0.6", "", "leveldb", nil)
	if err != nil {
		t.Fatalf("TestBackup06 NewCoreDB err! %v", err)

		return
	}

	err = cdb.Init()
	if err != nil {
		t.Fatalf("TestBackup06 Init err! %v", err)

		return
	}

	if cdb.GetPrivateKey().ToAddress() != "13VHbRCxFsiFk6qmVZjkDGXHzEJ553yMHd" {
		t.Fatalf("TestBackup06 GetPrivateKey err! (%v - %v)", cdb.GetPrivateKey().ToAddress(), "13VHbRCxFsiFk6qmVZjkDGXHzEJ553yMHd")

		return
	}

	numsNode := cdb.CountNodeNums()
	if numsNode != 30 {
		t.Fatalf("TestBackup06 CountNodeNums err! (%v - %v)", numsNode, 30)

		return
	}

	t.Logf("TestBackup06 is OK")
}
