GOPATH:=$(shell go env GOPATH)

.PHONY:proto
proto:
	@protoc --proto_path=${GOPATH}/src:.  --go_out=plugins=grpc:. proto/userpb/user.proto

.PHONY:das
das:pre
	go build -o bin/user-das-service main.go

pre:
	mkdir -p bin

clean:
	rm bin/user-das-service
