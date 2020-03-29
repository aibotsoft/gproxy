gen:
	protoc -I $$GOPATH/ -I=.  --gogo_out=plugins=grpc,paths=source_relative:. *.proto
