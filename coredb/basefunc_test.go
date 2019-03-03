package coredb

import (
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
)

func TestBaseFunc(t *testing.T) {
	//------------------------------------------------------------------------
	// initial ankaDB

	cfg := ankadb.NewConfig()

	cfg.PathDBRoot = "../test/coredbbasefunc"
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   "coredb",
		Engine: "leveldb",
		PathDB: "coredb",
	})

	dblogic, err := ankadb.NewBaseDBLogic(graphql.SchemaConfig{
		Query:    typeQuery,
		Mutation: typeMutation,
	})
	if err != nil {
		t.Fatalf("TestBaseFunc NewBaseDBLogic %v", err)

		return
	}

	ankaDB, err := ankadb.NewAnkaDB(cfg, dblogic)
	if ankaDB == nil {
		t.Fatalf("TestBaseFunc NewAnkaDB %v", err)

		return
	}

	//------------------------------------------------------------------------
	// newPrivateData

	privKey := jarviscrypto.GenerateKey()

	pd, err := newPrivateData(ankaDB, jarviscrypto.Base58Encode(privKey.ToPrivateBytes()), jarviscrypto.Base58Encode(privKey.ToPublicBytes()), privKey.ToAddress(), time.Now().Unix())
	if err != nil {
		t.Fatalf("TestBaseFunc newPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc getPrivateKey return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc newPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// getPrivateKey

	pd, err = getPrivateKey(ankaDB)
	if err != nil {
		t.Fatalf("TestBaseFunc getPrivateKey %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil {
		t.Fatalf("TestBaseFunc getPrivateKey return PriKey or PubKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc getPrivateKey return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	if pd.StrPriKey != jarviscrypto.Base58Encode(privKey.ToPrivateBytes()) {
		t.Fatalf("TestBaseFunc getPrivateKey return StrPriKey fail! (%v - %v)", pd.StrPriKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// getPrivateData

	pd, err = getPrivateData(ankaDB)
	if err != nil {
		t.Fatalf("TestBaseFunc getPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc getPrivateData return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc getPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// updPrivateData

	pd, err = updPrivateData(ankaDB, int64(100))
	if err != nil {
		t.Fatalf("TestBaseFunc updPrivateData %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc updPrivateData return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc updPrivateData return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// updNodeInfo test001

	test001PriKey := jarviscrypto.GenerateKey()

	cni1 := &coredbpb.NodeInfo{
		ServAddr:        "127.0.0.1",
		Addr:            test001PriKey.ToAddress(),
		Name:            "test001",
		ConnectNums:     0,
		ConnectedNums:   0,
		CtrlID:          0,
		LstClientAddr:   nil,
		AddTime:         time.Now().Unix(),
		ConnectMe:       false,
		NodeTypeVersion: "0.1.0",
		NodeType:        "normal",
		CoreVersion:     "0.1.0",
		LastRecvMsgID:   1,
		ConnType:        coredbpb.CONNECTTYPE_UNKNOWN_CONN,
	}

	cn1, err := updNodeInfo(ankaDB, cni1)
	if err != nil {
		t.Fatalf("TestBaseFunc updNodeInfo %v", err)

		return
	}

	if cni1.ServAddr != cn1.ServAddr {
		t.Fatalf("TestBaseFunc updNodeInfo return ServAddr fail! (%v - %v)", cni1.ServAddr, cn1.ServAddr)

		return
	}

	if cni1.Addr != cn1.Addr {
		t.Fatalf("TestBaseFunc updNodeInfo return Addr fail! (%v - %v)", cni1.Addr, cn1.Addr)

		return
	}

	if cni1.Name != cn1.Name {
		t.Fatalf("TestBaseFunc updNodeInfo return Name fail! (%v - %v)", cni1.Name, cn1.Name)

		return
	}

	if cni1.ConnectNums != cn1.ConnectNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectNums fail! (%v - %v)", cni1.ConnectNums, cn1.ConnectNums)

		return
	}

	if cni1.ConnectedNums != cn1.ConnectedNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

		return
	}

	if cni1.CtrlID != cn1.CtrlID {
		t.Fatalf("TestBaseFunc updNodeInfo return CtrlID fail! (%v - %v)", cni1.CtrlID, cn1.CtrlID)

		return
	}

	if cni1.ConnectedNums != cn1.ConnectedNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

		return
	}

	if cni1.AddTime != cn1.AddTime {
		t.Fatalf("TestBaseFunc updNodeInfo return AddTime fail! (%v - %v)", cni1.AddTime, cn1.AddTime)

		return
	}

	if cni1.ConnectMe != cn1.ConnectMe {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectMe fail! (%v - %v)", cni1.ConnectMe, cn1.ConnectMe)

		return
	}

	if cni1.ConnType != cn1.ConnType {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnType fail! (%v - %v)", cni1.ConnType, cn1.ConnType)

		return
	}

	if cni1.NodeTypeVersion != cn1.NodeTypeVersion {
		t.Fatalf("TestBaseFunc updNodeInfo return NodeTypeVersion fail! (%v - %v)", cni1.NodeTypeVersion, cn1.NodeTypeVersion)

		return
	}

	if cni1.NodeType != cn1.NodeType {
		t.Fatalf("TestBaseFunc updNodeInfo return NodeType fail! (%v - %v)", cni1.NodeType, cn1.NodeType)

		return
	}

	if cni1.CoreVersion != cn1.CoreVersion {
		t.Fatalf("TestBaseFunc updNodeInfo return CoreVersion fail! (%v - %v)", cni1.CoreVersion, cn1.CoreVersion)

		return
	}

	if cni1.LastRecvMsgID != cn1.LastRecvMsgID {
		t.Fatalf("TestBaseFunc updNodeInfo return LastRecvMsgID fail! (%v - %v)", cni1.LastRecvMsgID, cn1.LastRecvMsgID)

		return
	}

	//------------------------------------------------------------------------
	// updNodeInfo test002

	test002PriKey := jarviscrypto.GenerateKey()

	cni2 := &coredbpb.NodeInfo{
		ServAddr:        "127.0.0.1",
		Addr:            test002PriKey.ToAddress(),
		Name:            "test002",
		ConnectNums:     0,
		ConnectedNums:   0,
		CtrlID:          0,
		LstClientAddr:   nil,
		AddTime:         time.Now().Unix(),
		ConnectMe:       false,
		NodeTypeVersion: "0.1.0",
		NodeType:        "normal",
		CoreVersion:     "0.1.0",
		LastRecvMsgID:   1,
		ConnType:        coredbpb.CONNECTTYPE_UNKNOWN_CONN,
	}

	cn2, err := updNodeInfo(ankaDB, cni2)
	if err != nil {
		t.Fatalf("TestBaseFunc updNodeInfo %v", err)

		return
	}

	if cni2.ServAddr != cn2.ServAddr {
		t.Fatalf("TestBaseFunc updNodeInfo return ServAddr fail! (%v - %v)", cni2.ServAddr, cn2.ServAddr)

		return
	}

	if cni2.Addr != cn2.Addr {
		t.Fatalf("TestBaseFunc updNodeInfo return Addr fail! (%v - %v)", cni2.Addr, cn2.Addr)

		return
	}

	if cni2.Name != cn2.Name {
		t.Fatalf("TestBaseFunc updNodeInfo return Name fail! (%v - %v)", cni2.Name, cn2.Name)

		return
	}

	if cni2.ConnectNums != cn2.ConnectNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectNums fail! (%v - %v)", cni2.ConnectNums, cn2.ConnectNums)

		return
	}

	if cni2.ConnectedNums != cn2.ConnectedNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

		return
	}

	if cni2.CtrlID != cn2.CtrlID {
		t.Fatalf("TestBaseFunc updNodeInfo return CtrlID fail! (%v - %v)", cni2.CtrlID, cn2.CtrlID)

		return
	}

	if cni2.ConnectedNums != cn2.ConnectedNums {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

		return
	}

	if cni2.AddTime != cn2.AddTime {
		t.Fatalf("TestBaseFunc updNodeInfo return AddTime fail! (%v - %v)", cni2.AddTime, cn2.AddTime)

		return
	}

	if cni2.ConnectMe != cn2.ConnectMe {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnectMe fail! (%v - %v)", cni2.ConnectMe, cn2.ConnectMe)

		return
	}

	if cni2.ConnType != cn2.ConnType {
		t.Fatalf("TestBaseFunc updNodeInfo return ConnType fail! (%v - %v)", cni2.ConnType, cn2.ConnType)

		return
	}

	if cni2.NodeTypeVersion != cn2.NodeTypeVersion {
		t.Fatalf("TestBaseFunc updNodeInfo return NodeTypeVersion fail! (%v - %v)", cni2.NodeTypeVersion, cn2.NodeTypeVersion)

		return
	}

	if cni2.NodeType != cn2.NodeType {
		t.Fatalf("TestBaseFunc updNodeInfo return NodeType fail! (%v - %v)", cni2.NodeType, cn2.NodeType)

		return
	}

	if cni2.CoreVersion != cn2.CoreVersion {
		t.Fatalf("TestBaseFunc updNodeInfo return CoreVersion fail! (%v - %v)", cni2.CoreVersion, cn2.CoreVersion)

		return
	}

	if cni2.LastRecvMsgID != cn2.LastRecvMsgID {
		t.Fatalf("TestBaseFunc updNodeInfo return LastRecvMsgID fail! (%v - %v)", cni2.LastRecvMsgID, cn2.LastRecvMsgID)

		return
	}

	//------------------------------------------------------------------------
	// trustNode

	pd, err = trustNode(ankaDB, test001PriKey.ToAddress())
	if err != nil {
		t.Fatalf("TestBaseFunc trustNode %v", err)

		return
	}

	if pd.PriKey != nil || pd.PubKey != nil || pd.StrPriKey != "" {
		t.Fatalf("TestBaseFunc trustNode return PriKey or PubKey or StrPriKey is non-nil")

		return
	}

	if pd.StrPubKey != jarviscrypto.Base58Encode(privKey.ToPublicBytes()) {
		t.Fatalf("TestBaseFunc trustNode return StrPubKey fail! (%v - %v)", pd.StrPubKey, jarviscrypto.Base58Encode(privKey.ToPublicBytes()))

		return
	}

	//------------------------------------------------------------------------
	// getNodeInfo test001

	cn1, err = getNodeInfo(ankaDB, test001PriKey.ToAddress())
	if err != nil {
		t.Fatalf("TestBaseFunc getNodeInfo %v", err)

		return
	}

	if cni1.ServAddr != cn1.ServAddr {
		t.Fatalf("TestBaseFunc getNodeInfo return ServAddr fail! (%v - %v)", cni1.ServAddr, cn1.ServAddr)

		return
	}

	if cni1.Addr != cn1.Addr {
		t.Fatalf("TestBaseFunc getNodeInfo return Addr fail! (%v - %v)", cni1.Addr, cn1.Addr)

		return
	}

	if cni1.Name != cn1.Name {
		t.Fatalf("TestBaseFunc getNodeInfo return Name fail! (%v - %v)", cni1.Name, cn1.Name)

		return
	}

	if cni1.ConnectNums != cn1.ConnectNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectNums fail! (%v - %v)", cni1.ConnectNums, cn1.ConnectNums)

		return
	}

	if cni1.ConnectedNums != cn1.ConnectedNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

		return
	}

	if cni1.CtrlID != cn1.CtrlID {
		t.Fatalf("TestBaseFunc getNodeInfo return CtrlID fail! (%v - %v)", cni1.CtrlID, cn1.CtrlID)

		return
	}

	if cni1.ConnectedNums != cn1.ConnectedNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

		return
	}

	if cni1.AddTime != cn1.AddTime {
		t.Fatalf("TestBaseFunc getNodeInfo return AddTime fail! (%v - %v)", cni1.AddTime, cn1.AddTime)

		return
	}

	if cni1.ConnectMe != cn1.ConnectMe {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectMe fail! (%v - %v)", cni1.ConnectMe, cn1.ConnectMe)

		return
	}

	if cni1.ConnType != cn1.ConnType {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnType fail! (%v - %v)", cni1.ConnType, cn1.ConnType)

		return
	}

	if cni1.NodeTypeVersion != cn1.NodeTypeVersion {
		t.Fatalf("TestBaseFunc getNodeInfo return NodeTypeVersion fail! (%v - %v)", cni1.NodeTypeVersion, cn1.NodeTypeVersion)

		return
	}

	if cni1.NodeType != cn1.NodeType {
		t.Fatalf("TestBaseFunc getNodeInfo return NodeType fail! (%v - %v)", cni1.NodeType, cn1.NodeType)

		return
	}

	if cni1.CoreVersion != cn1.CoreVersion {
		t.Fatalf("TestBaseFunc getNodeInfo return CoreVersion fail! (%v - %v)", cni1.CoreVersion, cn1.CoreVersion)

		return
	}

	if cni1.LastRecvMsgID != cn1.LastRecvMsgID {
		t.Fatalf("TestBaseFunc getNodeInfo return LastRecvMsgID fail! (%v - %v)", cni1.LastRecvMsgID, cn1.LastRecvMsgID)

		return
	}

	//------------------------------------------------------------------------
	// getNodeInfo test002

	cn2, err = getNodeInfo(ankaDB, test002PriKey.ToAddress())
	if err != nil {
		t.Fatalf("TestBaseFunc getNodeInfo %v", err)

		return
	}

	if cni2.ServAddr != cn2.ServAddr {
		t.Fatalf("TestBaseFunc getNodeInfo return ServAddr fail! (%v - %v)", cni2.ServAddr, cn2.ServAddr)

		return
	}

	if cni2.Addr != cn2.Addr {
		t.Fatalf("TestBaseFunc getNodeInfo return Addr fail! (%v - %v)", cni2.Addr, cn2.Addr)

		return
	}

	if cni2.Name != cn2.Name {
		t.Fatalf("TestBaseFunc getNodeInfo return Name fail! (%v - %v)", cni2.Name, cn2.Name)

		return
	}

	if cni2.ConnectNums != cn2.ConnectNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectNums fail! (%v - %v)", cni2.ConnectNums, cn2.ConnectNums)

		return
	}

	if cni2.ConnectedNums != cn2.ConnectedNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

		return
	}

	if cni2.CtrlID != cn2.CtrlID {
		t.Fatalf("TestBaseFunc getNodeInfo return CtrlID fail! (%v - %v)", cni2.CtrlID, cn2.CtrlID)

		return
	}

	if cni2.ConnectedNums != cn2.ConnectedNums {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

		return
	}

	if cni2.AddTime != cn2.AddTime {
		t.Fatalf("TestBaseFunc getNodeInfo return AddTime fail! (%v - %v)", cni2.AddTime, cn2.AddTime)

		return
	}

	if cni2.ConnectMe != cn2.ConnectMe {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnectMe fail! (%v - %v)", cni2.ConnectMe, cn2.ConnectMe)

		return
	}

	if cni2.ConnType != cn2.ConnType {
		t.Fatalf("TestBaseFunc getNodeInfo return ConnType fail! (%v - %v)", cni2.ConnType, cn2.ConnType)

		return
	}

	if cni2.NodeTypeVersion != cn2.NodeTypeVersion {
		t.Fatalf("TestBaseFunc getNodeInfo return NodeTypeVersion fail! (%v - %v)", cni2.NodeTypeVersion, cn2.NodeTypeVersion)

		return
	}

	if cni2.NodeType != cn2.NodeType {
		t.Fatalf("TestBaseFunc getNodeInfo return NodeType fail! (%v - %v)", cni2.NodeType, cn2.NodeType)

		return
	}

	if cni2.CoreVersion != cn2.CoreVersion {
		t.Fatalf("TestBaseFunc getNodeInfo return CoreVersion fail! (%v - %v)", cni2.CoreVersion, cn2.CoreVersion)

		return
	}

	if cni2.LastRecvMsgID != cn2.LastRecvMsgID {
		t.Fatalf("TestBaseFunc getNodeInfo return LastRecvMsgID fail! (%v - %v)", cni2.LastRecvMsgID, cn2.LastRecvMsgID)

		return
	}

	//------------------------------------------------------------------------
	// getNodeInfos test002

	lstnode, err := getNodeInfos(ankaDB, -1, 0, 100)
	if err != nil {
		t.Fatalf("TestBaseFunc getNodeInfos %v", err)

		return
	}

	for i := 0; i < len(lstnode.Nodes); i++ {
		if lstnode.Nodes[i].Addr == cni1.Addr {
			cn1 = lstnode.Nodes[i]

			if cni1.ServAddr != cn1.ServAddr {
				t.Fatalf("TestBaseFunc getNodeInfos return ServAddr fail! (%v - %v)", cni1.ServAddr, cn1.ServAddr)

				return
			}

			if cni1.Addr != cn1.Addr {
				t.Fatalf("TestBaseFunc getNodeInfos return Addr fail! (%v - %v)", cni1.Addr, cn1.Addr)

				return
			}

			if cni1.Name != cn1.Name {
				t.Fatalf("TestBaseFunc getNodeInfos return Name fail! (%v - %v)", cni1.Name, cn1.Name)

				return
			}

			if cni1.ConnectNums != cn1.ConnectNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectNums fail! (%v - %v)", cni1.ConnectNums, cn1.ConnectNums)

				return
			}

			if cni1.ConnectedNums != cn1.ConnectedNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

				return
			}

			if cni1.CtrlID != cn1.CtrlID {
				t.Fatalf("TestBaseFunc getNodeInfos return CtrlID fail! (%v - %v)", cni1.CtrlID, cn1.CtrlID)

				return
			}

			if cni1.ConnectedNums != cn1.ConnectedNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectedNums fail! (%v - %v)", cni1.ConnectedNums, cn1.ConnectedNums)

				return
			}

			if cni1.AddTime != cn1.AddTime {
				t.Fatalf("TestBaseFunc getNodeInfos return AddTime fail! (%v - %v)", cni1.AddTime, cn1.AddTime)

				return
			}

			if cni1.ConnectMe != cn1.ConnectMe {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectMe fail! (%v - %v)", cni1.ConnectMe, cn1.ConnectMe)

				return
			}

			if cni1.ConnType != cn1.ConnType {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnType fail! (%v - %v)", cni1.ConnType, cn1.ConnType)

				return
			}

			if cni1.NodeTypeVersion != cn1.NodeTypeVersion {
				t.Fatalf("TestBaseFunc getNodeInfos return NodeTypeVersion fail! (%v - %v)", cni1.NodeTypeVersion, cn1.NodeTypeVersion)

				return
			}

			if cni1.NodeType != cn1.NodeType {
				t.Fatalf("TestBaseFunc getNodeInfos return NodeType fail! (%v - %v)", cni1.NodeType, cn1.NodeType)

				return
			}

			if cni1.CoreVersion != cn1.CoreVersion {
				t.Fatalf("TestBaseFunc getNodeInfos return CoreVersion fail! (%v - %v)", cni1.CoreVersion, cn1.CoreVersion)

				return
			}

			if cni1.LastRecvMsgID != cn1.LastRecvMsgID {
				t.Fatalf("TestBaseFunc getNodeInfos return LastRecvMsgID fail! (%v - %v)", cni1.LastRecvMsgID, cn1.LastRecvMsgID)

				return
			}
		} else if lstnode.Nodes[i].Addr == cni2.Addr {
			cn2 = lstnode.Nodes[i]

			if cni2.ServAddr != cn2.ServAddr {
				t.Fatalf("TestBaseFunc getNodeInfos return ServAddr fail! (%v - %v)", cni2.ServAddr, cn2.ServAddr)

				return
			}

			if cni2.Addr != cn2.Addr {
				t.Fatalf("TestBaseFunc getNodeInfos return Addr fail! (%v - %v)", cni2.Addr, cn2.Addr)

				return
			}

			if cni2.Name != cn2.Name {
				t.Fatalf("TestBaseFunc getNodeInfos return Name fail! (%v - %v)", cni2.Name, cn2.Name)

				return
			}

			if cni2.ConnectNums != cn2.ConnectNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectNums fail! (%v - %v)", cni2.ConnectNums, cn2.ConnectNums)

				return
			}

			if cni2.ConnectedNums != cn2.ConnectedNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

				return
			}

			if cni2.CtrlID != cn2.CtrlID {
				t.Fatalf("TestBaseFunc getNodeInfos return CtrlID fail! (%v - %v)", cni2.CtrlID, cn2.CtrlID)

				return
			}

			if cni2.ConnectedNums != cn2.ConnectedNums {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectedNums fail! (%v - %v)", cni2.ConnectedNums, cn2.ConnectedNums)

				return
			}

			if cni2.AddTime != cn2.AddTime {
				t.Fatalf("TestBaseFunc getNodeInfos return AddTime fail! (%v - %v)", cni2.AddTime, cn2.AddTime)

				return
			}

			if cni2.ConnectMe != cn2.ConnectMe {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnectMe fail! (%v - %v)", cni2.ConnectMe, cn2.ConnectMe)

				return
			}

			if cni2.ConnType != cn2.ConnType {
				t.Fatalf("TestBaseFunc getNodeInfos return ConnType fail! (%v - %v)", cni2.ConnType, cn2.ConnType)

				return
			}

			if cni2.NodeTypeVersion != cn2.NodeTypeVersion {
				t.Fatalf("TestBaseFunc getNodeInfos return NodeTypeVersion fail! (%v - %v)", cni2.NodeTypeVersion, cn2.NodeTypeVersion)

				return
			}

			if cni2.NodeType != cn2.NodeType {
				t.Fatalf("TestBaseFunc getNodeInfos return NodeType fail! (%v - %v)", cni2.NodeType, cn2.NodeType)

				return
			}

			if cni2.CoreVersion != cn2.CoreVersion {
				t.Fatalf("TestBaseFunc getNodeInfos return CoreVersion fail! (%v - %v)", cni2.CoreVersion, cn2.CoreVersion)

				return
			}

			if cni2.LastRecvMsgID != cn2.LastRecvMsgID {
				t.Fatalf("TestBaseFunc getNodeInfos return LastRecvMsgID fail! (%v - %v)", cni2.LastRecvMsgID, cn2.LastRecvMsgID)

				return
			}
		}
	}

	t.Logf("TestBaseFunc is OK")
}
