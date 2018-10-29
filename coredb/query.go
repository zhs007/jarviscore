package coredb

import (
	"github.com/golang/protobuf/proto"
	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/ankadb/graphqlext"
	"github.com/zhs007/ankadb/proto"
	pb "github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
)

var typeQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"privateKey": &graphql.Field{
				Type: typePrivateData,
				Args: graphql.FieldConfigArgument{},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
					if anka == nil {
						return nil, ankadb.ErrCtxAnkaDB
					}

					curdb := anka.MgrDB.GetDB("coredb")
					if curdb == nil {
						return nil, ankadb.ErrCtxCurDB
					}

					pd := &pb.PrivateData{}
					err := ankadb.GetMsgFromDB(curdb, []byte(keyMyPrivateData), pd)
					if err != nil {
						return nil, err
					}

					pd.StrPriKey = jarviscrypto.Base58Encode(pd.PriKey)
					pd.StrPubKey = jarviscrypto.Base58Encode(pd.PubKey)
					pd.PriKey = nil
					pd.PubKey = nil

					return pd, nil
				},
			},
			"privateData": &graphql.Field{
				Type: typePrivateData,
				Args: graphql.FieldConfigArgument{},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
					if anka == nil {
						return nil, ankadb.ErrCtxAnkaDB
					}

					curdb := anka.MgrDB.GetDB("coredb")
					if curdb == nil {
						return nil, ankadb.ErrCtxCurDB
					}

					pd := &pb.PrivateData{}
					err := ankadb.GetMsgFromDB(curdb, []byte(keyMyPrivateData), pd)
					if err != nil {
						return nil, err
					}

					// private key not allow query
					pd.StrPubKey = jarviscrypto.Base58Encode(pd.PubKey)
					pd.PriKey = nil
					pd.PubKey = nil

					return pd, nil
				},
			},
			"nodeInfo": &graphql.Field{
				Type: typeNodeInfo,
				Args: graphql.FieldConfigArgument{
					"addr": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
					if anka == nil {
						return nil, ankadb.ErrCtxAnkaDB
					}

					curdb := anka.MgrDB.GetDB("coredb")
					if curdb == nil {
						return nil, ankadb.ErrCtxCurDB
					}

					addr := params.Args["addr"].(string)

					keyid := makeNodeInfoKeyID(addr)
					buf, err := curdb.Get([]byte(keyid))
					td := &pb.NodeInfo{}

					err = proto.Unmarshal(buf, td)
					if err != nil {
						return nil, ankadb.ErrQuertResultDecode
					}

					return td, nil
				},
			},
			"nodeInfos": &graphql.Field{
				Type: typeNodeInfoList,
				Args: graphql.FieldConfigArgument{
					"snapshotID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphqlext.Int64),
					},
					"beginIndex": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"nums": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
					if anka == nil {
						return nil, ankadb.ErrCtxAnkaDB
					}

					curdb := anka.MgrDB.GetDB("coredb")
					if curdb == nil {
						return nil, ankadb.ErrCtxCurDB
					}

					mgrSnapshot := anka.MgrDB.GetMgrSnapshot("coredb")
					if mgrSnapshot == nil {
						return nil, ankadb.ErrCtxSnapshotMgr
					}

					snapshotID := params.Args["snapshotID"].(int64)
					beginIndex := params.Args["beginIndex"].(int)
					nums := params.Args["nums"].(int)
					if beginIndex < 0 || nums <= 0 {
						return nil, ankadb.ErrQuertParams
					}

					lstNodeInfo := &pb.NodeInfoList{}
					var pSnapshot *ankadbpb.Snapshot

					if snapshotID > 0 {
						pSnapshot = mgrSnapshot.Get(snapshotID)
					} else {
						var err error
						pSnapshot, err = mgrSnapshot.NewSnapshot([]byte(prefixKeyNodeInfo))
						if err != nil {
							return nil, ankadb.ErrCtxSnapshotMgr
						}
					}

					lstNodeInfo.SnapshotID = pSnapshot.SnapshotID
					lstNodeInfo.MaxIndex = int32(len(pSnapshot.Keys))

					curi := beginIndex
					for ; curi < len(pSnapshot.Keys) && len(lstNodeInfo.Nodes) < nums; curi++ {
						ni := pb.NodeInfo{}
						err := ankadb.GetMsgFromDB(curdb, []byte(pSnapshot.Keys[curi]), &ni)
						if err == nil {
							lstNodeInfo.Nodes = append(lstNodeInfo.Nodes)
						}
					}

					lstNodeInfo.EndIndex = int32(curi)

					return lstNodeInfo, nil
				},
			},
		},
	},
)
