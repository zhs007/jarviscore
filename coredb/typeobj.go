package coredb

import (
	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb/graphqlext"
)

// typePrivateData - PrivateData
//		you can see coredb.graphql
var typePrivateData = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PrivateData",
		Fields: graphql.Fields{
			"strPriKey": &graphql.Field{
				Type: graphql.ID,
			},
			"strPubKey": &graphql.Field{
				Type: graphql.ID,
			},
			"createTime": &graphql.Field{
				Type: graphqlext.Timestamp,
			},
			"onlineTime": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"addr": &graphql.Field{
				Type: graphql.ID,
			},
			"lstTrustNode": &graphql.Field{
				Type: graphql.NewList(graphql.ID),
			},
		},
	},
)

// typeNodeInfo - NodeInfo
//		you can see coredb.graphql
var typeNodeInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "NodeInfo",
		Fields: graphql.Fields{
			"addr": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"servAddr": &graphql.Field{
				Type: graphql.ID,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"connectNums": &graphql.Field{
				Type: graphql.Int,
			},
			"connectedNums": &graphql.Field{
				Type: graphql.Int,
			},
			"ctrlID": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"lstClientAddr": &graphql.Field{
				Type: graphql.NewList(graphql.ID),
			},
			"addTime": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"connectMe": &graphql.Field{
				Type: graphql.Boolean,
			},
			"nodeTypeVersion": &graphql.Field{
				Type: graphql.String,
			},
			"nodeType": &graphql.Field{
				Type: graphql.String,
			},
			"coreVersion": &graphql.Field{
				Type: graphql.String,
			},
			"lastSendMsgID": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"lastConnectTime": &graphql.Field{
				Type: graphqlext.Timestamp,
			},
			"lastConnectedTime": &graphql.Field{
				Type: graphqlext.Timestamp,
			},
			"lastConnectMeTime": &graphql.Field{
				Type: graphqlext.Timestamp,
			},
			"lstGroups": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"lastRecvMsgID": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"deprecated": &graphql.Field{
				Type: graphql.Boolean,
			},
			"connType": &graphql.Field{
				Type: graphql.Int,
			},
			"validConnNodes": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"timestampDeprecated": &graphql.Field{
				Type: graphqlext.Int64,
			},
			"numsConnectFail": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// typeNodeInfoList - NodeInfoList
//		you can see coredb.graphql
var typeNodeInfoList = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "NodeInfoList",
		Fields: graphql.Fields{
			"snapshotID": &graphql.Field{
				Type: graphql.NewNonNull(graphqlext.Int64),
			},
			"endIndex": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"maxIndex": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nodes": &graphql.Field{
				Type: graphql.NewList(typeNodeInfo),
			},
		},
	},
)
