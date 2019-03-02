package coredb

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
	"github.com/zhs007/jarviscore/proto"
)

func TestCoreDB(t *testing.T) {
	jarvisbase.InitLogger(zapcore.DebugLevel, true, "")

	//------------------------------------------------------------------------
	// initial CoreDB

	cdb, err := NewCoreDB("../test/testcoredb", "", "leveldb")
	if err != nil {
		t.Fatalf("TestCoreDB NewCoreDB err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// loadPrivateKey

	err = cdb.loadPrivateKey()
	if err != nil {
		t.Fatalf("TestCoreDB loadPrivateKey err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// loadAllNodes

	err = cdb.loadAllNodes()
	if err != nil {
		t.Fatalf("TestCoreDB loadAllNodes err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// GetMyData

	mydata, err := cdb.GetMyData()
	if err != nil {
		t.Fatalf("TestCoreDB GetMyData err! %v", err)

		return
	}

	if mydata.Addr != cdb.GetPrivateKey().ToAddress() {
		t.Fatalf("TestCoreDB GetMyData addr err! (%v - %v)", mydata.Addr, cdb.GetPrivateKey().ToAddress())

		return
	}

	if mydata.CreateTime <= 0 {
		t.Fatalf("TestCoreDB GetMyData CreateTime err! %v", mydata.CreateTime)

		return
	}

	if mydata.StrPubKey != jarviscrypto.Base58Encode(cdb.GetPrivateKey().ToPublicBytes()) {
		t.Fatalf("TestCoreDB GetMyData StrPubKey err! (%v - %v)", mydata.StrPubKey, jarviscrypto.Base58Encode(cdb.GetPrivateKey().ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// UpdNodeBaseInfo 001

	privKey001 := jarviscrypto.GenerateKey()

	nbi001 := &jarviscorepb.NodeBaseInfo{
		ServAddr:        "127.0.0.1:9001",
		Addr:            privKey001.ToAddress(),
		Name:            "node001",
		NodeTypeVersion: "0.1.0",
		NodeType:        "test node",
		CoreVersion:     "0.7.0",
	}

	err = cdb.UpdNodeBaseInfo(nbi001)
	if err != nil {
		t.Fatalf("TestCoreDB UpdNodeBaseInfo err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// UpdNodeBaseInfo 002

	privKey002 := jarviscrypto.GenerateKey()

	nbi002 := &jarviscorepb.NodeBaseInfo{
		ServAddr:        "127.0.0.1:9002",
		Addr:            privKey002.ToAddress(),
		Name:            "node002",
		NodeTypeVersion: "0.1.0",
		NodeType:        "test node",
		CoreVersion:     "0.7.0",
	}

	err = cdb.UpdNodeBaseInfo(nbi002)
	if err != nil {
		t.Fatalf("TestCoreDB UpdNodeBaseInfo err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// UpdNodeInfo 001

	err = cdb.UpdNodeInfo(privKey001.ToAddress())
	if err != nil {
		t.Fatalf("TestCoreDB UpdNodeInfo err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// UpdNodeInfo 002

	err = cdb.UpdNodeInfo(privKey002.ToAddress())
	if err != nil {
		t.Fatalf("TestCoreDB UpdNodeInfo err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// UpdNodeInfo non-node

	err = cdb.UpdNodeInfo("1234567890")
	if err == nil || err != ErrCoreDBHasNotNode {
		t.Fatalf("TestCoreDB UpdNodeInfo err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// TrustNode 001

	err = cdb.TrustNode(privKey001.ToAddress())
	if err != nil {
		t.Fatalf("TestCoreDB TrustNode err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// TrustNode 002

	err = cdb.TrustNode(privKey002.ToAddress())
	if err != nil {
		t.Fatalf("TestCoreDB TrustNode err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// TrustNode non-node

	err = cdb.TrustNode("1234567890")
	if err != nil {
		t.Fatalf("TestCoreDB TrustNode err! %v", err)

		return
	}

	//------------------------------------------------------------------------
	// GetNodes

	lst, err := cdb.GetNodes(100)
	if err != nil {
		t.Fatalf("TestCoreDB GetNodes err! %v", err)

		return
	}

	if len(lst.Nodes) != 2 {
		t.Fatalf("TestCoreDB GetNodes node number is %v", len(lst.Nodes))

		return
	}

	//------------------------------------------------------------------------
	// HasNode non-node

	isok := cdb.hasNode("1234567890")
	if isok {
		t.Fatalf("TestCoreDB hasNode non-node err!")

		return
	}

	//------------------------------------------------------------------------
	// HasNode 001

	isok = cdb.hasNode(privKey001.ToAddress())
	if !isok {
		t.Fatalf("TestCoreDB hasNode node001 err!")

		return
	}

	//------------------------------------------------------------------------
	// HasNode 002

	isok = cdb.hasNode(privKey002.ToAddress())
	if !isok {
		t.Fatalf("TestCoreDB hasNode node002 err!")

		return
	}

	//------------------------------------------------------------------------
	// getNode non-node

	cn := cdb.GetNode("1234567890")
	if cn != nil {
		t.Fatalf("TestCoreDB GetNode non-node err!")

		return
	}

	//------------------------------------------------------------------------
	// getNode 001

	cn = cdb.GetNode(privKey001.ToAddress())
	if cn == nil || cn.ServAddr != nbi001.ServAddr {
		t.Fatalf("TestCoreDB GetNode node001 err!")

		return
	}

	//------------------------------------------------------------------------
	// getNode 002

	cn = cdb.GetNode(privKey002.ToAddress())
	if cn == nil || cn.ServAddr != nbi002.ServAddr {
		t.Fatalf("TestCoreDB getNode node002 err!")

		return
	}

	//------------------------------------------------------------------------
	// foreachNodeEx

	cdb.foreachNodeEx(func(addr string, cni *coredbpb.NodeInfo) {
		if addr == privKey001.ToAddress() {
			if cni.ServAddr != nbi001.ServAddr {
				t.Fatalf("TestCoreDB foreachNodeEx node001 err!")
			}
		} else if addr == privKey002.ToAddress() {
			if cni.ServAddr != nbi002.ServAddr {
				t.Fatalf("TestCoreDB foreachNodeEx node002 err!")
			}
		} else {
			t.Fatalf("TestCoreDB foreachNodeEx err!")
		}
	})

	t.Logf("TestCoreDB is OK")
}
