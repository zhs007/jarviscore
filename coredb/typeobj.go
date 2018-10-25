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
			"priKey": &graphql.Field{
				Type: graphql.ID,
			},
			"pubKey": &graphql.Field{
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
			"beginIndex": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"nodes": &graphql.Field{
				Type: graphql.NewList(typeNodeInfo),
			},
		},
	},
)
