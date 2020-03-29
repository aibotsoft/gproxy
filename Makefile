gen:
	#(protoc -I=. --gogo_out=. *.proto)
	protoc -I $$GOPATH/ -I=.  --gogo_out=plugins=grpc,paths=source_relative:. *.proto
	#protoc --gogofast_out=plugins=grpc:. *.proto
