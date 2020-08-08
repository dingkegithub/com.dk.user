GOPATH:=$(shell go env GOPATH)

.PHONY:logic
logic:pre
	go build -o bin/user-logic-service main.go

pre:
	mkdir -p bin

clean:
	rm bin/user-logic-service
