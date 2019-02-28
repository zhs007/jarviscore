package coredb

// // newdb - new db logic
// func newdb(dbpath string, httpAddr string, engine string) (*ankadb.AnkaDB, error) {
// 	cfg := ankadb.NewConfig()

// 	cfg.AddrHTTP = httpAddr
// 	cfg.PathDBRoot = dbpath
// 	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
// 		Name:   "coredb",
// 		Engine: engine,
// 		PathDB: path.Join(dbpath, "coredb"),
// 	})

// 	ankaDB, err := ankadb.NewAnkaDB(cfg, newDBLogic())
// 	if ankaDB == nil {
// 		jarvisbase.Error("newdb", zap.Error(err))

// 		return nil, err
// 	}

// 	jarvisbase.Info("newdb", zap.String("dbpath", dbpath),
// 		zap.String("httpAddr", httpAddr), zap.String("engine", engine))

// 	return ankaDB, err
// }

// // dblogic -
// type dblogic struct {
// 	schema graphql.Schema
// }

// // newDBLogic -
// func newDBLogic() ankadb.DBLogic {
// 	return ankadb.NewBaseDBLogic()
// 	var schema, _ = graphql.NewSchema(
// 		graphql.SchemaConfig{
// 			Query:    typeQuery,
// 			Mutation: typeMutation,
// 			// Types:    curTypes,
// 		},
// 	)

// 	return &dblogic{
// 		schema: schema,
// 	}
// }

// // OnQuery -
// func (logic *dblogic) OnQuery(ctx context.Context, request string, values map[string]interface{}) (*graphql.Result, error) {
// 	result := graphql.Do(graphql.Params{
// 		Schema:         logic.schema,
// 		RequestString:  request,
// 		VariableValues: values,
// 		Context:        ctx,
// 	})
// 	// if len(result.Errors) > 0 {
// 	// 	fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
// 	// }

// 	return result, nil
// }

// // OnQueryStream -
// func (logic *dblogic) OnQueryStream(ctx context.Context, request string, values map[string]interface{}, funcOnQueryStream ankadb.FuncOnQueryStream) error {
// 	return nil
// }
