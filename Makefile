test1:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent
	./metricstest -test.v -test.run=^TestIteration1$$  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port="8080"

runServer:
	go run ./cmd/server
