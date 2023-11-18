port = 8080
.PHONY: test
test:
	go build -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent
	./metricstest -test.v -test.run=^TestIteration1$$  -binary-path=cmd/server/server   -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration2  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration3  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration4  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration5  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration6  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration7  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration8  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)
	./metricstest -test.v -test.run=^TestIteration9  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"
	./metricstest -test.v -test.run=^TestIteration10  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"
	./metricstest -test.v -test.run=^TestIteration11  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"
	./metricstest -test.v -test.run=^TestIteration12  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"
	./metricstest -test.v -test.run=^TestIteration13  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"
	./metricstest -test.v -test.run=^TestIteration14  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"  -key="superKey"
	./metricstest -test.v -test.run=^TestIteration15  -binary-path=cmd/server/server  -agent-binary-path=cmd/agent/agent  -source-path=.  -server-port=$(port)  -file-storage-path="./tests.txt"  -database-dsn="postgres://postgres:pass@localhost:5432/metrics?sslmode=disable"  -key="superKey"

.PHONY: runServer
runServer:
	go run ./cmd/server -a ":$(port)" -f "./tests.txt" -d "postgres://postgres:pass@localhost:5432/metrics?sslmode=disable" -k "superKey"

.PHONY: runAgent
runAgent:
	go build -o ./cmd/agent/ ./cmd/agent/ && \
	./cmd/agent/agent -a ":$(port)" -r 10 -p 0 -k "superKey" -l 1 -pprofAddress ":8080"

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