protoc -I proto/ proto/jarviscore.proto --go_out=plugins=grpc:proto
protoc -I proto/ proto/jarvistask.proto --go_out=plugins=grpc:proto
protoc -I coredb/proto/ coredb/proto/coredb.proto --go_out=plugins=grpc:coredb/proto