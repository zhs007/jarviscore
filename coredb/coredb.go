package coredb

import (
	"context"
	"path"

	"github.com/graphql-go/graphql"

	"github.com/zhs007/ankadb"
)

// NewCoreDB - new core db
func NewCoreDB(dbpath string, httpAddr string, engine string) (*ankadb.AnkaDB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   "coredb",
		Engine: engine,
		PathDB: path.Join(dbpath, "coredb"),
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, newDBLogic())
	if ankaDB == nil {
		return nil, err
	}

	return ankaDB, err
}

// coreDB -
type coreDB struct {
	schema graphql.Schema
}

// newDBLogic -
func newDBLogic() ankadb.DBLogic {
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    typeQuery,
			Mutation: typeMutation,
			// Types:    curTypes,
		},
	)

	return &coreDB{
		schema: schema,
	}
}

// OnQuery -
func (logic *coreDB) OnQuery(ctx context.Context, request string, values map[string]interface{}) (*graphql.Result, error) {
	result := graphql.Do(graphql.Params{
		Schema:         logic.schema,
		RequestString:  request,
		VariableValues: values,
		Context:        ctx,
	})
	// if len(result.Errors) > 0 {
	// 	fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	// }

	return result, nil
}

// OnQueryStream -
func (logic *coreDB) OnQueryStream(ctx context.Context, request string, values map[string]interface{}, funcOnQueryStream ankadb.FuncOnQueryStream) error {
	return nil
}
