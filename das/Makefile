GOPATH:=$(shell go env GOPATH)

.PHONY:proto
proto:
	@protoc --proto_path=${GOPATH}/src:.  --go_out=plugins=grpc:. proto/userpb/user.proto
