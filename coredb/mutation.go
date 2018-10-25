package coredb

import (
	"encoding/base64"

	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb"
	"github.com/zhs007/ankadb/err"
	"github.com/zhs007/ankadb/graphqlext"
	"github.com/zhs007/ankadb/proto"
	pb "github.com/zhs007/jarviscore/coredb/proto"
)

var typeMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"newPrivateData": &graphql.Field{
			Type:        typePrivateData,
			Description: "new private data",
			Args: graphql.FieldConfigArgument{
				"priKey": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"pubKey": &graphql.ArgumentConfig{
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
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_ANKADB_ERR)
				}

				curdb := anka.MgrDB.GetDB("coredb")
				if curdb == nil {
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_CURDB_ERR)
				}

				priKey := params.Args["priKey"].(string)
				pubKey := params.Args["pubKey"].(string)
				addr := params.Args["addr"].(string)
				createTime := params.Args["createTime"].(int64)

				priKeyBytes, err := base64.StdEncoding.DecodeString(priKey)
				if err != nil {
					return nil, ankadberr.NewError(ankadbpb.CODE_QUERY_PARAM_ERR)
				}

				pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
				if err != nil {
					return nil, ankadberr.NewError(ankadbpb.CODE_QUERY_PARAM_ERR)
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

				return pd, nil
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
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_ANKADB_ERR)
				}

				curdb := anka.MgrDB.GetDB("coredb")
				if curdb == nil {
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_CURDB_ERR)
				}

				curOnlineTime := params.Args["curOnlineTime"].(int64)

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
				pd.PriKey = nil

				return pd, nil
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
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_ANKADB_ERR)
				}

				curdb := anka.MgrDB.GetDB("coredb")
				if curdb == nil {
					return nil, ankadberr.NewError(ankadbpb.CODE_CTX_CURDB_ERR)
				}

				ni := &pb.NodeInfo{}
				err := ankadb.GetMsgFromParam(params, "nodeInfo", ni)
				if err != nil {
					return nil, err
				}

				err = ankadb.PutMsg2DB(curdb, []byte(makeNodeInfoKeyID(ni.Addr)), ni)
				if err != nil {
					return nil, err
				}

				return ni, nil
			},
		},
	},
})
