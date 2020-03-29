gen:
	protoc -I $$GOPATH/ -I=.  --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,paths=source_relative:. *.proto
