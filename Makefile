test:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent
	./metricstest -test.v -test.run=^TestIteration1$$  -binary-path=cmd/server/server
	./metricstest -test.v -test.run=^TestIteration2  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.
	./metricstest -test.v -test.run=^TestIteration3  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.
	./metricstest -test.v -test.run=^TestIteration4  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"

runServer:
	go run ./cmd/server -a :8080
runAgent:
	go run ./cmd/agent -a :8080 -r 10 -p 2

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

build:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent

mytest:
	go test -v -count=1 ./...