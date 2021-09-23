export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=protos/ --go_out=./jarviscore2pb --go_opt=paths=source_relative protos/*.proto
protoc --proto_path=protos/ --go-grpc_out=./jarviscore2pb --go-grpc_opt=paths=source_relative protos/*.proto