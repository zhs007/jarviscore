package coredb

import (
	"github.com/graphql-go/graphql"
	"github.com/zhs007/ankadb/graphqlext"
)

// inputTypeNodeInfo - NodeInfo
//		you can see coredb.graphql
var inputTypeNodeInfo = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "NodeInfoInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"addr": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"servAddr": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"connectNums": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"connectedNums": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"ctrlID": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Int64,
			},
			"lstClientAddr": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.ID),
			},
			"addTime": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Int64,
			},
			"connectMe": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"nodeTypeVersion": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"nodeType": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"coreVersion": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"lastSendMsgID": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Int64,
			},
			"lastConnectTime": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Timestamp,
			},
			"lastConnectedTime": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Timestamp,
			},
			"lastConnectMeTime": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Timestamp,
			},
			"lstGroups": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
			"lastRecvMsgID": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Int64,
			},
			"deprecated": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"connType": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"validConnNodes": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
			"timestampDeprecated": &graphql.InputObjectFieldConfig{
				Type: graphqlext.Int64,
			},
			"numsConnectFail": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	},
)
