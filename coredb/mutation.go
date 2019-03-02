package coredb

import (
	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/ankadb/graphqlext"
	pb "github.com/zhs007/jarviscore/coredb/proto"
	"github.com/zhs007/jarviscore/crypto"
)

// newPrivateData - new private data
func newPrivateData(anka ankadb.AnkaDB, priKey string, pubKey string, addr string,
	createTime int64) (*pb.PrivateData, error) {

	curdb := anka.GetDBMgr().GetDB("coredb")
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	priKeyBytes, err := jarviscrypto.Base58Decode(priKey)
	if err != nil {
		return nil, ankadb.ErrQuertParams
	}

	pubKeyBytes, err := jarviscrypto.Base58Decode(pubKey)
	if err != nil {
		return nil, ankadb.ErrQuertParams
	}

	pd := &pb.PrivateData{
		PriKey:     priKeyBytes,
		PubKey:     pubKeyBytes,
		CreateTime: createTime,
		Addr:       addr,
		OnlineTime: 0,
	}

	err = ankadb.PutMsg2DB(curdb, []byte(keyMyPrivateData), pd)
	if err != nil {
		return nil, err
	}

	// private key not allow query
	pd.PriKey = nil
	pd.PubKey = nil
	pd.StrPubKey = pubKey

	return pd, nil
}

// updPrivateData - update private data
func updPrivateData(anka ankadb.AnkaDB, curOnlineTime int64) (*pb.PrivateData, error) {
	curdb := anka.GetDBMgr().GetDB("coredb")
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	pd := &pb.PrivateData{}
	err := ankadb.GetMsgFromDB(curdb, []byte(keyMyPrivateData), pd)
	if err != nil {
		return nil, err
	}

	pd.OnlineTime += curOnlineTime

	err = ankadb.PutMsg2DB(curdb, []byte(keyMyPrivateData), pd)
	if err != nil {
		return nil, err
	}

	// private key not allow query
	pd.StrPubKey = jarviscrypto.Base58Encode(pd.PubKey)
	pd.PriKey = nil
	pd.PubKey = nil

	return pd, nil
}

// trustNode - trust a node
func trustNode(anka ankadb.AnkaDB, addr string) (*pb.PrivateData, error) {
	curdb := anka.GetDBMgr().GetDB("coredb")
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	pd := &pb.PrivateData{}
	err := ankadb.GetMsgFromDB(curdb, []byte(keyMyPrivateData), pd)
	if err != nil {
		return nil, err
	}

	hasaddr := false
	if len(pd.LstTrustNode) > 0 {
		for i := range pd.LstTrustNode {
			if pd.LstTrustNode[i] == addr {
				hasaddr = true

				break
			}
		}
	}

	if !hasaddr {
		pd.LstTrustNode = append(pd.LstTrustNode, addr)
	}

	err = ankadb.PutMsg2DB(curdb, []byte(keyMyPrivateData), pd)
	if err != nil {
		return nil, err
	}

	// private key not allow query
	pd.StrPubKey = jarviscrypto.Base58Encode(pd.PubKey)
	pd.PriKey = nil
	pd.PubKey = nil

	return pd, nil
}

// updNodeInfo - update node info
func updNodeInfo(anka ankadb.AnkaDB, ni *pb.NodeInfo) (*pb.NodeInfo, error) {
	curdb := anka.GetDBMgr().GetDB("coredb")
	if curdb == nil {
		return nil, ankadb.ErrCtxCurDB
	}

	err := ankadb.PutMsg2DB(curdb, []byte(makeNodeInfoKeyID(ni.Addr)), ni)
	if err != nil {
		return nil, err
	}

	return ni, nil
}

// typeMutation - define mutation for graphql
var typeMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"newPrivateData": &graphql.Field{
			Type:        typePrivateData,
			Description: "new private data",
			Args: graphql.FieldConfigArgument{
				"strPriKey": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"strPubKey": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"addr": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"createTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphqlext.Timestamp),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
				if anka == nil {
					return nil, ankadb.ErrCtxAnkaDB
				}

				priKey := params.Args["strPriKey"].(string)
				pubKey := params.Args["strPubKey"].(string)
				addr := params.Args["addr"].(string)
				createTime := params.Args["createTime"].(int64)

				return newPrivateData(anka, priKey, pubKey, addr, createTime)
			},
		},
		"updPrivateData": &graphql.Field{
			Type:        typePrivateData,
			Description: "upd private data",
			Args: graphql.FieldConfigArgument{
				"curOnlineTime": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphqlext.Int64),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
				if anka == nil {
					return nil, ankadb.ErrCtxAnkaDB
				}

				curOnlineTime := params.Args["curOnlineTime"].(int64)

				return updPrivateData(anka, curOnlineTime)
			},
		},
		"trustNode": &graphql.Field{
			Type:        typePrivateData,
			Description: "trust node",
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

				addr := params.Args["addr"].(string)

				return trustNode(anka, addr)
			},
		},
		"updNodeInfo": &graphql.Field{
			Type:        typeNodeInfo,
			Description: "update node info",
			Args: graphql.FieldConfigArgument{
				"nodeInfo": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(inputTypeNodeInfo),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				anka := ankadb.GetContextValueAnkaDB(params.Context, interface{}("ankadb"))
				if anka == nil {
					return nil, ankadb.ErrCtxAnkaDB
				}

				ni := &pb.NodeInfo{}
				err := ankadb.GetMsgFromParam(params, "nodeInfo", ni)
				if err != nil {
					return nil, err
				}

				return updNodeInfo(anka, ni)
			},
		},
	},
})
