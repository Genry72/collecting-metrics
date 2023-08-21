.PHONY: test
test:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent
	./metricstest -test.v -test.run=^TestIteration1$$  -binary-path=cmd/server/server
	./metricstest -test.v -test.run=^TestIteration2  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.
	./metricstest -test.v -test.run=^TestIteration3  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.
	./metricstest -test.v -test.run=^TestIteration4  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"
	./metricstest -test.v -test.run=^TestIteration5  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"
	./metricstest -test.v -test.run=^TestIteration6  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"
	./metricstest -test.v -test.run=^TestIteration7  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"

.PHONY: runServer
runServer:
	go run ./cmd/server -a :8080

.PHONY: runAgent
runAgent:
	go run ./cmd/agent -a :8080 -r 10 -p 2

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: build
build:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent

.PHONY: mytest
mytest:
	go test -v -count=1 ./...