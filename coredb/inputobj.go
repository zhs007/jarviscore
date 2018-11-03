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
			"connectNode": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
	},
)
