
.PHONY:das
das:pre
	go build -o bin/user-das-service das/main.go

.PHONY:logic
logic:pre
	go build -o bin/user-logic-service logic/main.go

pre:
	mkdir -p bin

clean:
	rm bin/user-das-service
	rm bin/user-logic-service