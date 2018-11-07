package jarviscore

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/zhs007/ankadb"
	"github.com/zhs007/jarviscore/base"
	"github.com/zhs007/jarviscore/coredb"
	"github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
)

const (
	coredbMyPrivKey        = "myprivkey"
	coredbMyNodeInfoPrefix = "ni:"
)

const queryNewPrivateData = `mutation NewPrivateData($strPriKey: ID!, $strPubKey: ID!, $addr: ID!, $createTime: Timestamp!) {
	newPrivateData(strPriKey: $strPriKey, strPubKey: $strPubKey, addr: $addr, createTime: $createTime) {
		strPubKey, addr, createTime
	}
}`

const queryPrivateKey = `{
	privateKey {
	  strPriKey
	  strPubKey
	  addr
	  createTime
	  lstTrustNode
	}
}`

const queryPrivateData = `{
	privateData {
	  strPubKey
	  addr
	  createTime
	  lstTrustNode
	}
}`

const queryNodeInfos = `query NodeInfos($snapshotID: Int64!, $beginIndex: Int!, $nums: Int!) {
	nodeInfos(snapshotID: $snapshotID, beginIndex: $beginIndex, nums: $nums) {
		snapshotID, endIndex, maxIndex, 
		nodes {
			addr, servAddr, name, connectNums, connectedNums, ctrlID, lstClientAddr, addTime
		}
	}
}`

const queryUpdNodeInfo = `mutation UpdNodeInfo($nodeInfo: NodeInfoInput!) {
	updNodeInfo(nodeInfo: $nodeInfo) {
		addr, servAddr, name, connectNums, connectedNums, ctrlID, lstClientAddr, addTime, connectMe, connectNode
	}
}`

const queryTrustNode = `mutation TrustNode($addr: ID!) {
	trustNode(addr: $addr) {
		strPubKey, addr, createTime, lstTrustNode
	}
}`

// CoreDB - jarvisnode core database
type CoreDB struct {
	ankaDB       *ankadb.AnkaDB
	privKey      *jarviscrypto.PrivateKey
	lstTrustNode []string
}

func newCoreDB() (*CoreDB, error) {
	ankaDB, err := coredb.NewCoreDB(config.AnkaDB.DBPath, config.AnkaDB.HTTPServ, config.AnkaDB.Engine)
	if err != nil {
		jarvisbase.Error("newCoreDB:NewAnkaLDB", zap.Error(err))

		return nil, err
	}

	return &CoreDB{
		ankaDB: ankaDB,
	}, nil
}

func (db *CoreDB) savePrivateKey() error {
	if db.privKey == nil {
		jarvisbase.Error("savePrivateKey", zap.Error(ErrNoPrivateKey))

		return ErrNoPrivateKey
	}

	params := make(map[string]interface{})
	params["strPriKey"] = jarviscrypto.Base58Encode(db.privKey.ToPrivateBytes())
	params["strPubKey"] = jarviscrypto.Base58Encode(db.privKey.ToPublicBytes())
	params["addr"] = db.privKey.ToAddress()
	params["createTime"] = time.Now().Unix()

	ret, err := db.ankaDB.LocalQuery(context.Background(), queryNewPrivateData, params)
	if err != nil {
		jarvisbase.Error("savePrivateKey", zap.Error(err))

		return err
	}

	jarvisbase.Info("savePrivateKey",
		jarvisbase.JSON("result", ret))

	return nil
}

func (db *CoreDB) loadPrivateKeyEx() error {
	err := db._loadPrivateKey()
	if err != nil {
		jarvisbase.Info("loadPrivateKeyEx:_loadPrivateKey",
			zap.Error(err))

		db.privKey = jarviscrypto.GenerateKey()
		jarvisbase.Info("loadPrivateKeyEx:GenerateKey",
			zap.String("privkey", db.privKey.ToAddress()))

		return db.savePrivateKey()
	}

	jarvisbase.Info("loadPrivateKeyEx:OK",
		zap.String("privkey", db.privKey.ToAddress()))

	if len(config.LstTrustNode) > 0 {
		for i := range config.LstTrustNode {
			if !db.IsTrustNode(config.LstTrustNode[i]) {
				db.TrustNode(config.LstTrustNode[i])
			}
		}
	}

	return nil
}

func (db *CoreDB) _loadPrivateKey() error {
	result, err := db.ankaDB.LocalQuery(context.Background(), queryPrivateKey, nil)
	if err != nil {
		return err
	}

	jarvisbase.Info("_loadPrivateKey",
		jarvisbase.JSON("result", result))

	if result.HasErrors() {
		return result.Errors[0]
	}

	rpd := &coredb.ResultPrivateKey{}
	err = ankadb.MakeObjFromResult(result, rpd)
	if err != nil {
		return err
	}

	jarvisbase.Info("_loadPrivateKey",
		jarvisbase.JSON("rpd", rpd))

	bytesPrikey, err := jarviscrypto.Base58Decode(rpd.PrivateKey.StrPriKey)
	if err != nil {
		return err
	}

	// tmp := jarviscrypto.Base58Encode(bytesPrikey)
	// jarvisbase.Info("_loadPrivateKey",
	// 	zap.String("recheck base58", tmp))

	privkey := jarviscrypto.NewPrivateKey()
	err = privkey.FromBytes(bytesPrikey)
	if err != nil {
		return err
		// db.privKey = jarviscrypto.GenerateKey()

		// return db.savePrivateKey()
	}

	db.privKey = privkey
	db.lstTrustNode = rpd.PrivateKey.LstTrustNode

	return nil
}

func (db *CoreDB) _foreachNode(oneach func(string, *coredbpb.NodeInfo), snapshotID int64, beginIndex int, nums int) (*coredbpb.NodeInfoList, error) {
	params := make(map[string]interface{})
	params["snapshotID"] = snapshotID
	params["beginIndex"] = beginIndex
	params["nums"] = nums
	params["createTime"] = time.Now().Unix()

	result, err := db.ankaDB.LocalQuery(context.Background(), queryNodeInfos, params)
	rnis := &coredbpb.NodeInfoList{}
	err = ankadb.MakeMsgFromResult(result, rnis)
	if err != nil {
		return nil, err
	}

	for _, v := range rnis.Nodes {
		oneach(v.Addr, v)
	}

	return rnis, nil
}

func (db *CoreDB) foreachNodeEx(oneach func(string, *coredbpb.NodeInfo)) error {
	rnis, err := db._foreachNode(oneach, 0, 0, 128)
	if err != nil {
		return err
	}

	for bi := rnis.EndIndex; bi < rnis.MaxIndex; {
		rnis, err = db._foreachNode(oneach, rnis.SnapshotID, int(bi), 128)
		if err != nil {
			return err
		}

		bi = rnis.EndIndex
	}

	return nil
}

// func (db *coreDB) saveNode(cni *NodeInfo) error {
// 	ni := &coredbpb.NodeInfo{
// 		ServAddr:      cni.baseinfo.ServAddr,
// 		Addr:          cni.baseinfo.Addr,
// 		Name:          cni.baseinfo.Name,
// 		ConnectNums:   int32(cni.connectNums),
// 		ConnectedNums: int32(cni.connectedNums),
// 	}

// 	params := make(map[string]interface{})

// 	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", ni)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := db.ankaDB.LocalQuery(context.Background(), queryUpdNodeInfo, params)
// 	if err != nil {
// 		return err
// 	}

// 	jarvisbase.Info("saveNode", jarvisbase.JSON("result", result))

// 	return nil
// }

func (db *CoreDB) saveNodeEx(cni *coredbpb.NodeInfo) error {
	// ni := &coredbpb.NodeInfo{
	// 	ServAddr:      cni.baseinfo.ServAddr,
	// 	Addr:          cni.baseinfo.Addr,
	// 	Name:          cni.baseinfo.Name,
	// 	ConnectNums:   int32(cni.connectNums),
	// 	ConnectedNums: int32(cni.connectedNums),
	// }

	params := make(map[string]interface{})

	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", cni)
	if err != nil {
		return err
	}

	result, err := db.ankaDB.LocalQuery(context.Background(), queryUpdNodeInfo, params)
	if err != nil {
		return err
	}

	jarvisbase.Info("saveNode", jarvisbase.JSON("result", result))

	return nil
}

// TrustNode - trust node with addr
func (db *CoreDB) TrustNode(addr string) (string, error) {
	params := make(map[string]interface{})
	params["addr"] = addr

	ret, err := db.ankaDB.LocalQuery(context.Background(), queryTrustNode, params)
	if err != nil {
		jarvisbase.Error("trustNode", zap.Error(err))

		return err.Error(), err
	}

	s, err := json.Marshal(ret)
	if err != nil {
		jarvisbase.Error("CoreDB.GetMyState", zap.Error(err))

		return err.Error(), err
	}

	jarvisbase.Info("trustNode",
		zap.String("result", string(s)))

	rpd := &coredb.ResultPrivateKey{}
	err = ankadb.MakeObjFromResult(ret, rpd)
	if err != nil {
		return err.Error(), err
	}

	db.lstTrustNode = rpd.PrivateKey.LstTrustNode

	return string(s), nil
}

// IsTrustNode - is trust node with addr
func (db *CoreDB) IsTrustNode(addr string) bool {
	if len(db.lstTrustNode) <= 0 {
		return false
	}

	for i := range db.lstTrustNode {
		if db.lstTrustNode[i] == addr {
			return true
		}
	}

	return false
}

// GetMyState - get my state
func (db *CoreDB) GetMyState() (string, error) {
	ret, err := db.ankaDB.LocalQuery(context.Background(), queryPrivateData, nil)
	if err != nil {
		jarvisbase.Error("CoreDB.GetMyState", zap.Error(err))

		return err.Error(), err
	}

	s, err := json.Marshal(ret)
	if err != nil {
		jarvisbase.Error("CoreDB.GetMyState", zap.Error(err))

		return err.Error(), err
	}

	jarvisbase.Info("GetMyState",
		zap.String("result", string(s)))

	rpd := &coredb.ResultPrivateKey{}
	err = ankadb.MakeObjFromResult(ret, rpd)
	if err != nil {
		return err.Error(), err
	}

	db.lstTrustNode = rpd.PrivateKey.LstTrustNode

	return string(s), nil
}

// GetNodes - get jarvis nodes
func (db *CoreDB) GetNodes(nums int) (string, error) {
	params := make(map[string]interface{})
	params["snapshotID"] = 0
	params["beginIndex"] = 0
	params["nums"] = nums
	params["createTime"] = time.Now().Unix()

	ret, err := db.ankaDB.LocalQuery(context.Background(), queryNodeInfos, params)

	s, err := json.Marshal(ret)
	if err != nil {
		jarvisbase.Error("CoreDB.GetNodes", zap.Error(err))

		return err.Error(), err
	}

	jarvisbase.Info("GetNodes",
		zap.String("result", string(s)))

	return string(s), nil
}
