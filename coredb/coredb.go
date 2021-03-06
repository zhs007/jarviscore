package coredb

import (
	"context"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/ankadb/database"
	jarvisbase "github.com/zhs007/jarviscore/base"
	coredbpb "github.com/zhs007/jarviscore/coredb/proto"
	jarviscrypto "github.com/zhs007/jarviscore/crypto"
	jarviscorepb "github.com/zhs007/jarviscore/proto"
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
			addr, servAddr, name, connectNums, connectedNums, ctrlID, lstClientAddr, addTime, 
			nodeType, coreVersion, nodeTypeVersion, lastSendMsgID, lastRecvMsgID, deprecated, timestampDeprecated, numsConnectFail
		}
	}
}`

const queryUpdNodeInfo = `mutation UpdNodeInfo($nodeInfo: NodeInfoInput!) {
	updNodeInfo(nodeInfo: $nodeInfo) {
		addr, servAddr, name, connectNums, connectedNums, ctrlID, lstClientAddr, addTime, 
		connectMe, connType, nodeType, coreVersion, nodeTypeVersion, lastSendMsgID, lastRecvMsgID, deprecated, timestampDeprecated, numsConnectFail
	}
}`

const queryTrustNode = `mutation TrustNode($addr: ID!) {
	trustNode(addr: $addr) {
		strPubKey, addr, createTime, lstTrustNode
	}
}`

// CoreDB - jarvisnode core database
type CoreDB struct {
	ankaDB       ankadb.AnkaDB
	privKey      *jarviscrypto.PrivateKey
	lstTrustNode []string
	mapNodes     sync.Map
	cfgTrustNode []string
}

// NewCoreDB -
func NewCoreDB(dbpath string, httpAddr string, engine string, lstTrust []string) (*CoreDB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   CoreDBName,
		Engine: engine,
		PathDB: CoreDBName,
	})

	dblogic, err := ankadb.NewBaseDBLogic(graphql.SchemaConfig{
		Query:    typeQuery,
		Mutation: typeMutation,
	})
	if err != nil {
		jarvisbase.Error("newdb", zap.Error(err))

		return nil, err
	}

	ankaDB, err := ankadb.NewAnkaDB(cfg, dblogic)
	if ankaDB == nil {
		jarvisbase.Error("newdb", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newdb", zap.String("dbpath", dbpath),
		zap.String("httpAddr", httpAddr), zap.String("engine", engine))

	return &CoreDB{
		ankaDB:       ankaDB,
		cfgTrustNode: lstTrust,
	}, nil
}

// GetPrivateKey -
func (db *CoreDB) GetPrivateKey() *jarviscrypto.PrivateKey {
	return db.privKey
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

	ret, err := db.ankaDB.Query(context.Background(), queryNewPrivateData, params)
	if err != nil {
		jarvisbase.Error("savePrivateKey", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(ret)
	if err != nil {
		jarvisbase.Error("CoreDB.savePrivateKey:GetResultError", zap.Error(err))

		return err
	}

	// jarvisbase.Info("savePrivateKey",
	// 	jarvisbase.JSON("result", ret))

	return nil
}

// Init -
func (db *CoreDB) Init() error {
	err := db.loadPrivateKey()
	if err != nil {
		jarvisbase.Error("CoreDB.Init:loadPrivateKey", zap.Error(err))

		return err
	}

	err = db.TrustNodes(db.cfgTrustNode)
	if err != nil {
		jarvisbase.Error("CoreDB.Init:TrustNodes", zap.Error(err))

		return err
	}

	err = db.loadAllNodes()
	if err != nil {
		jarvisbase.Error("CoreDB.Init:loadAllNodes", zap.Error(err))

		return err
	}

	return nil
}

// loadPrivateKey -
func (db *CoreDB) loadPrivateKey() error {
	err := db._loadPrivateKey()
	if err != nil {
		if err != database.ErrNotFound {
			jarvisbase.Info("loadPrivateKey:_loadPrivateKey",
				zap.Error(err))
		}

		db.privKey = jarviscrypto.GenerateKey()
		jarvisbase.Info("loadPrivateKey:GenerateKey",
			zap.String("privkey", db.privKey.ToAddress()))

		return db.savePrivateKey()
	}

	myaddr := db.privKey.ToAddress()

	jarvisbase.Info("loadPrivateKey:OK",
		zap.String("privkey", myaddr))

	return nil
}

func (db *CoreDB) _loadPrivateKey() error {
	result, err := db.ankaDB.Query(context.Background(), queryPrivateKey, nil)
	if err != nil {
		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB._loadPrivateKey:GetResultError", zap.Error(err))

		return err
	}

	pd := &coredbpb.PrivateData{}
	err = ankadb.MakeMsgFromResultEx(result, "privateKey", pd)
	if err != nil {
		jarvisbase.Warn("CoreDB._loadPrivateKey:MakeMsgFromResultEx", zap.Error(err))

		return err
	}

	if pd.Addr == "" {
		return database.ErrNotFound
	}

	jarvisbase.Info("_loadPrivateKey",
		jarvisbase.JSON("privateKey", pd))

	bytesPrikey, err := jarviscrypto.Base58Decode(pd.StrPriKey)
	if err != nil {
		jarvisbase.Warn("CoreDB._loadPrivateKey:Base58Decode", zap.Error(err))

		return err
	}

	privkey := jarviscrypto.NewPrivateKey()
	err = privkey.FromBytes(bytesPrikey)
	if err != nil {
		jarvisbase.Warn("CoreDB._loadPrivateKey:privkey.FromBytes", zap.Error(err))

		return err
	}

	db.privKey = privkey
	db.lstTrustNode = pd.LstTrustNode

	return nil
}

func (db *CoreDB) _foreachNode(oneach func(string, *coredbpb.NodeInfo), snapshotID int64, beginIndex int, nums int) (*coredbpb.NodeInfoList, error) {
	params := make(map[string]interface{})
	params["snapshotID"] = snapshotID
	params["beginIndex"] = beginIndex
	params["nums"] = nums
	// params["createTime"] = time.Now().Unix()

	result, err := db.ankaDB.Query(context.Background(), queryNodeInfos, params)
	if err != nil {
		jarvisbase.Warn("CoreDB._foreachNode:Query", zap.Error(err))

		return nil, err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB._foreachNode:GetResultError", zap.Error(err))

		return nil, err
	}

	lst := &coredbpb.NodeInfoList{}
	err = ankadb.MakeMsgFromResultEx(result, "nodeInfos", lst)
	if err != nil {
		jarvisbase.Warn("CoreDB._foreachNode:MakeMsgFromResultEx", zap.Error(err))

		return nil, err
	}

	for _, v := range lst.Nodes {
		oneach(v.Addr, v)
	}

	return lst, nil
}

func (db *CoreDB) foreachNodeEx(oneach func(string, *coredbpb.NodeInfo)) error {
	rnis, err := db._foreachNode(oneach, 0, 0, 128)
	if err != nil {
		jarvisbase.Warn("CoreDB.foreachNodeEx:_foreachNode", zap.Error(err))

		return err
	}

	for bi := rnis.EndIndex; bi < rnis.MaxIndex; {
		rnis, err = db._foreachNode(oneach, rnis.SnapshotID, int(bi), 128)
		if err != nil {
			jarvisbase.Warn("CoreDB.foreachNodeEx:for:_foreachNode:", zap.Error(err))

			return err
		}

		bi = rnis.EndIndex
	}

	return nil
}

// loadAllNodes -
func (db *CoreDB) loadAllNodes() error {
	curnodes := 0

	err := db.foreachNodeEx(func(key string, val *coredbpb.NodeInfo) {
		val.ConnectMe = false
		val.ConnType = coredbpb.CONNECTTYPE_UNKNOWN_CONN
		val.Deprecated = false
		val.LastMsgID4RequestNodes = 0
		val.NodesVersion = ""
		val.LastNodesVersion = ""

		if val.LastRecvMsgID < 0 {
			val.LastRecvMsgID = 0
		}

		db.mapNodes.Store(val.Addr, val)
		curnodes++
	})
	if err != nil {
		jarvisbase.Warn("CoreDB.loadAllNodes", zap.Error(err))

		return err
	}

	jarvisbase.Info("CoreDB.loadAllNodes", zap.Int("nodes", curnodes))

	return nil
}

// UpdNodeBaseInfo -
func (db *CoreDB) UpdNodeBaseInfo(ni *jarviscorepb.NodeBaseInfo) error {
	cni := db.GetNode(ni.Addr)
	if cni == nil {
		cni = &coredbpb.NodeInfo{
			ServAddr:      ni.ServAddr,
			Addr:          ni.Addr,
			Name:          ni.Name,
			ConnectNums:   0,
			ConnectedNums: 0,
			// CtrlID:          0,
			LstClientAddr:   nil,
			AddTime:         time.Now().Unix(),
			ConnectMe:       false,
			NodeTypeVersion: ni.NodeTypeVersion,
			NodeType:        ni.NodeType,
			CoreVersion:     ni.CoreVersion,
			LastRecvMsgID:   0,
			ConnType:        coredbpb.CONNECTTYPE_UNKNOWN_CONN,
		}
	} else {
		cni.ServAddr = ni.ServAddr
		cni.Name = ni.Name
		cni.NodeTypeVersion = ni.NodeTypeVersion
		cni.NodeType = ni.NodeType
		cni.CoreVersion = ni.CoreVersion
	}

	params := make(map[string]interface{})

	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", cni)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeBaseInfo:MakeParamsFromMsg", zap.Error(err))

		return err
	}

	result, err := db.ankaDB.Query(context.Background(), queryUpdNodeInfo, params)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeBaseInfo:Query", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeBaseInfo:GetResultError", zap.Error(err))

		return err
	}

	// jarvisbase.Debug("updNodeBaseInfo", jarvisbase.JSON("result", result))

	db.mapNodes.Store(cni.Addr, cni)

	return nil
}

// OnNodeConnected -
func (db *CoreDB) OnNodeConnected(ni *jarviscorepb.NodeBaseInfo) error {
	cni := db.GetNode(ni.Addr)
	if cni == nil {
		cni = &coredbpb.NodeInfo{
			ServAddr:      ni.ServAddr,
			Addr:          ni.Addr,
			Name:          ni.Name,
			ConnectNums:   0,
			ConnectedNums: 0,
			// CtrlID:          0,
			LstClientAddr:   nil,
			AddTime:         time.Now().Unix(),
			ConnectMe:       false,
			NodeTypeVersion: ni.NodeTypeVersion,
			NodeType:        ni.NodeType,
			CoreVersion:     ni.CoreVersion,
			LastRecvMsgID:   0,
			ConnType:        coredbpb.CONNECTTYPE_UNKNOWN_CONN,
		}
	} else {
		cni.ServAddr = ni.ServAddr
		cni.Name = ni.Name
		cni.NodeTypeVersion = ni.NodeTypeVersion
		cni.NodeType = ni.NodeType
		cni.CoreVersion = ni.CoreVersion
		cni.Deprecated = false
	}

	params := make(map[string]interface{})

	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", cni)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnNodeConnected:MakeParamsFromMsg", zap.Error(err))

		return err
	}

	result, err := db.ankaDB.Query(context.Background(), queryUpdNodeInfo, params)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnNodeConnected:Query", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnNodeConnected:GetResultError", zap.Error(err))

		return err
	}

	db.mapNodes.Store(cni.Addr, cni)

	return nil
}

// OnIConnectedNode -
func (db *CoreDB) OnIConnectedNode(ni *jarviscorepb.NodeBaseInfo) error {
	cni := db.GetNode(ni.Addr)
	if cni == nil {
		cni = &coredbpb.NodeInfo{
			ServAddr:      ni.ServAddr,
			Addr:          ni.Addr,
			Name:          ni.Name,
			ConnectNums:   0,
			ConnectedNums: 0,
			// CtrlID:          0,
			LstClientAddr:   nil,
			AddTime:         time.Now().Unix(),
			ConnectMe:       false,
			NodeTypeVersion: ni.NodeTypeVersion,
			NodeType:        ni.NodeType,
			CoreVersion:     ni.CoreVersion,
			LastRecvMsgID:   0,
			ConnType:        coredbpb.CONNECTTYPE_UNKNOWN_CONN,
		}
	} else {
		cni.ServAddr = ni.ServAddr
		cni.Name = ni.Name
		cni.NodeTypeVersion = ni.NodeTypeVersion
		cni.NodeType = ni.NodeType
		cni.CoreVersion = ni.CoreVersion
		cni.Deprecated = false
	}

	params := make(map[string]interface{})

	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", cni)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnIConnectedNode:MakeParamsFromMsg", zap.Error(err))

		return err
	}

	result, err := db.ankaDB.Query(context.Background(), queryUpdNodeInfo, params)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnIConnectedNode:Query", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.OnIConnectedNode:GetResultError", zap.Error(err))

		return err
	}

	db.mapNodes.Store(cni.Addr, cni)

	return nil
}

// UpdNodeInfo -
func (db *CoreDB) UpdNodeInfo(addr string) error {
	cni := db.GetNode(addr)
	if cni == nil {
		jarvisbase.Warn("CoreDB.UpdNodeInfo", zap.Error(ErrCoreDBHasNotNode))

		return ErrCoreDBHasNotNode
	}

	params := make(map[string]interface{})

	err := ankadb.MakeParamsFromMsg(params, "nodeInfo", cni)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeInfo:MakeParamsFromMsg", zap.Error(err))

		return err
	}

	result, err := db.ankaDB.Query(context.Background(), queryUpdNodeInfo, params)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeInfo:Query", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.UpdNodeInfo:GetResultError", zap.Error(err))

		return err
	}

	// jarvisbase.Debug("UpdNodeInfo", jarvisbase.JSON("result", result))

	db.mapNodes.Store(cni.Addr, cni)

	return nil
}

// TrustNode - trust node with addr
func (db *CoreDB) TrustNode(addr string) error {
	if db.IsTrustNode(addr) {
		return nil
	}

	params := make(map[string]interface{})
	params["addr"] = addr

	result, err := db.ankaDB.Query(context.Background(), queryTrustNode, params)
	if err != nil {
		jarvisbase.Error("trustNode", zap.Error(err))

		return err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.TrustNode:GetResultError", zap.Error(err))

		return err
	}

	pd := &coredbpb.PrivateData{}
	err = ankadb.MakeMsgFromResultEx(result, "trustNode", pd)
	if err != nil {
		jarvisbase.Warn("CoreDB.TrustNode:MakeMsgFromResultEx", zap.Error(err))

		return err
	}

	db.lstTrustNode = pd.LstTrustNode

	return nil
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

// GetMyData - get my state
func (db *CoreDB) GetMyData() (*coredbpb.PrivateData, error) {
	result, err := db.ankaDB.Query(context.Background(), queryPrivateData, nil)
	if err != nil {
		jarvisbase.Error("CoreDB.GetMyData", zap.Error(err))

		return nil, err
	}

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.GetMyData:GetResultError", zap.Error(err))

		return nil, err
	}

	pd := &coredbpb.PrivateData{}
	err = ankadb.MakeMsgFromResultEx(result, "privateData", pd)
	if err != nil {
		jarvisbase.Error("CoreDB.GetMyData:MakeMsgFromResultEx err", zap.Error(err))

		return nil, err
	}

	db.lstTrustNode = pd.LstTrustNode

	return pd, nil
}

// GetNodes - get jarvis nodes
func (db *CoreDB) GetNodes(nums int) (*coredbpb.NodeInfoList, error) {
	params := make(map[string]interface{})
	params["snapshotID"] = int64(0)
	params["beginIndex"] = 0
	params["nums"] = nums
	params["createTime"] = time.Now().Unix()

	result, err := db.ankaDB.Query(context.Background(), queryNodeInfos, params)

	err = ankadb.GetResultError(result)
	if err != nil {
		jarvisbase.Warn("CoreDB.GetNodes:GetResultError", zap.Error(err))

		return nil, err
	}

	lst := &coredbpb.NodeInfoList{}
	err = ankadb.MakeMsgFromResultEx(result, "nodeInfos", lst)
	if err != nil {
		jarvisbase.Warn("CoreDB.GetNodes:MakeMsgFromResultEx", zap.Error(err))

		return nil, err
	}

	return lst, nil
}

// has node
func (db *CoreDB) hasNode(addr string) bool {
	return db.GetNode(addr) != nil
}

// GetNode - get node with addr
func (db *CoreDB) GetNode(addr string) *coredbpb.NodeInfo {
	// var cni *coredbpb.NodeInfo
	ifcni, ok := db.mapNodes.Load(addr)
	if ok {
		cni, typeok := ifcni.(*coredbpb.NodeInfo)
		if typeok {
			return cni
		}
	}

	return nil
}

// FindNodeWithServAddr - get node
func (db *CoreDB) FindNodeWithServAddr(servaddr string) *coredbpb.NodeInfo {
	var curnode *coredbpb.NodeInfo

	db.mapNodes.Range(func(key, value interface{}) bool {
		cni, ok := value.(*coredbpb.NodeInfo)
		if ok && cni.ServAddr == servaddr {
			curnode = cni

			return false
		}

		return true
	})

	return curnode
}

// Start - start
func (db *CoreDB) Start(ctx context.Context) error {
	return db.ankaDB.Start(ctx)
}

// ForEachMapNodes - foreach mapNodes
func (db *CoreDB) ForEachMapNodes(oneach func(string, *coredbpb.NodeInfo) error) error {
	var curerr error

	db.mapNodes.Range(func(key, value interface{}) bool {
		curaddr, addrok := key.(string)
		curnode, nodeok := value.(*coredbpb.NodeInfo)
		if addrok && nodeok {
			err := oneach(curaddr, curnode)
			if err != nil {
				curerr = err

				return false
			}
		}

		return true
	})

	return curerr
}

// FindMapNode - find node in mapNodes
func (db *CoreDB) FindMapNode(name string) *coredbpb.NodeInfo {
	var curnode *coredbpb.NodeInfo

	db.mapNodes.Range(func(key, value interface{}) bool {
		cni, ok := value.(*coredbpb.NodeInfo)
		if ok && cni.Name == name {
			curnode = cni

			return false
		}

		return true
	})

	return curnode
}

// Close - close database
func (db *CoreDB) Close() {
	db.ankaDB.GetDBMgr().GetDB(CoreDBName).Close()
}

// GetNewSendMsgID - get msgid
func (db *CoreDB) GetNewSendMsgID(addr string) int64 {
	curnode := db.GetNode(addr)
	if curnode != nil {
		curnode.LastSendMsgID = curnode.LastSendMsgID + 1

		db.UpdNodeInfo(addr)

		return curnode.LastSendMsgID
	}

	return 0
}

// GetCurRecvMsgID - get msgid
func (db *CoreDB) GetCurRecvMsgID(addr string) int64 {
	curnode := db.GetNode(addr)
	if curnode != nil {
		return curnode.LastRecvMsgID
	}

	return 1
}

// UpdSendMsgID - update msgid
func (db *CoreDB) UpdSendMsgID(addr string, msgid int64) {
	curnode := db.GetNode(addr)
	if curnode != nil {
		if curnode.LastSendMsgID < msgid {
			curnode.LastSendMsgID = msgid

			db.UpdNodeInfo(addr)
		}
	}
}

// UpdRecvMsgID - update msgid
func (db *CoreDB) UpdRecvMsgID(addr string, msgid int64) {
	curnode := db.GetNode(addr)
	if curnode != nil {
		curnode.LastRecvMsgID = msgid

		db.UpdNodeInfo(addr)
	}
}

// UpdMsgID - update msgid
func (db *CoreDB) UpdMsgID(addr string, sendmsgid int64, recvmsgid int64) {
	curnode := db.GetNode(addr)
	if curnode != nil {
		if sendmsgid > curnode.LastSendMsgID {
			curnode.LastSendMsgID = sendmsgid
		}

		curnode.LastRecvMsgID = recvmsgid

		db.UpdNodeInfo(addr)
	}
}

// CountNodeNums - count all node nums
func (db *CoreDB) CountNodeNums() int {
	nums := 0
	db.mapNodes.Range(func(key, value interface{}) bool {
		_, ok := value.(*coredbpb.NodeInfo)
		if ok {
			nums++
		}

		return true
	})

	return nums
}

// TrustNodes - trust some nodes
func (db *CoreDB) TrustNodes(lstNode []string) error {
	for _, v := range lstNode {
		err := db.TrustNode(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddTask - add task
func (db *CoreDB) AddTask(ctx context.Context, task *jarviscorepb.JarvisTask) error {
	key := makeTaskKey(task)

	curdb := db.ankaDB.GetDBMgr().GetDB(CoreDBName)
	if curdb == nil {
		return ankadb.ErrCtxCurDB
	}

	err := ankadb.PutMsg2DB(curdb, []byte(key), task)
	if err != nil {
		return err
	}

	return nil
}

// CountMyNodesVersion - count my nodes version
func (db *CoreDB) CountMyNodesVersion() string {
	lst := &coredbpb.NodeInfoList2{}

	db.ForEachMapNodes(func(key string, v *coredbpb.NodeInfo) error {
		//!!! skip the deprecated node
		if v.Deprecated {
			return nil
		}

		curnode := &coredbpb.NodeInfo{
			ServAddr: v.ServAddr,
			Addr:     v.Addr,
			Name:     v.Name,
		}

		lst.Nodes = append(lst.Nodes, curnode)

		return nil
	})

	sort.Slice(lst.Nodes, func(i, j int) bool {
		return lst.Nodes[i].Addr > lst.Nodes[j].Addr
	})

	str, err := jarvisbase.MD5Protobuf(lst)
	if err != nil {
		jarvisbase.Warn("CoreDB.CountMyNodesVersion:MD5Protobuf",
			zap.Error(err))

		return ""
	}

	return str
}

// SetNodesVersion - set nodes version
func (db *CoreDB) SetNodesVersion(addr string, nodesVersion string) error {
	cn := db.GetNode(addr)
	if cn != nil {
		cn.NodesVersion = nodesVersion
	}

	return nil
}

// SetLastNodesVersion - set last nodes version
func (db *CoreDB) SetLastNodesVersion(addr string, nodesVersion string) error {
	cn := db.GetNode(addr)
	if cn != nil {
		cn.LastNodesVersion = nodesVersion
	}

	return nil
}

// AddLogInfo - add loginfo
func (db *CoreDB) AddLogInfo(ctx context.Context, loginfo *jarviscorepb.LogInfo) error {
	key := makeLogInfoKey(loginfo)

	curdb := db.ankaDB.GetDBMgr().GetDB(CoreDBName)
	if curdb == nil {
		return ankadb.ErrCtxCurDB
	}

	err := ankadb.PutMsg2DB(curdb, []byte(key), loginfo)
	if err != nil {
		return err
	}

	return nil
}

// DelLogInfo - del loginfo
func (db *CoreDB) DelLogInfo(ctx context.Context, loginfo *jarviscorepb.LogInfo) error {
	key := makeLogInfoKey(loginfo)

	curdb := db.ankaDB.GetDBMgr().GetDB(CoreDBName)
	if curdb == nil {
		return ankadb.ErrCtxCurDB
	}

	err := curdb.Delete([]byte(key))
	if err != nil {
		return err
	}

	return nil
}

// GetAllLogInfos - get loginfo
func (db *CoreDB) GetAllLogInfos(ctx context.Context) ([]*jarviscorepb.LogInfo, error) {
	var lst []*jarviscorepb.LogInfo

	curdb := db.ankaDB.GetDBMgr().GetDB(CoreDBName)
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	it := curdb.NewIteratorWithPrefix([]byte(prefixLog))
	if it.Error() != nil {
		jarvisbase.Warn("CoreDB.GetAllLogInfos:NewIteratorWithPrefix",
			zap.Error(it.Error()))

		return nil, it.Error()
	}

	for {
		if it.Valid() {
			li := jarviscorepb.LogInfo{}
			err := ankadb.GetMsgFromIterator(it, &li)
			if err == nil {
				lst = append(lst, &li)
			}
		}

		if !it.Next() {
			break
		}
	}

	return lst, nil
}

// GetAppLogInfos - get app loginfo
func (db *CoreDB) GetAppLogInfos(ctx context.Context, appName string) ([]*jarviscorepb.LogInfo, error) {
	var lst []*jarviscorepb.LogInfo

	curdb := db.ankaDB.GetDBMgr().GetDB(CoreDBName)
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	it := curdb.NewIteratorWithPrefix([]byte(prefixLog + appName + ":"))
	if it.Error() != nil {
		jarvisbase.Warn("CoreDB.GetAppLogInfos:NewIteratorWithPrefix",
			zap.Error(it.Error()))

		return nil, it.Error()
	}

	for {
		if it.Valid() {
			li := jarviscorepb.LogInfo{}
			err := ankadb.GetMsgFromIterator(it, &li)
			if err == nil {
				lst = append(lst, &li)
			}
		}

		if !it.Next() {
			break
		}
	}

	return lst, nil
}
