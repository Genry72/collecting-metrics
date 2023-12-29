port = 8080
buildFlag = -ldflags="-X 'main.buildVersion=`git describe --tags --abbrev=0`' -X 'main.buildDate=`date`' -X 'main.buildCommit=`git rev-parse HEAD`'"
.PHONY: test
test: build
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
runServer: build
	./cmd/server/server -a ":$(port)" -f "./tests.txt" -d "postgres://postgres:pass@localhost:5432/metrics?sslmode=disable" -k "superKey" -crypto-key "./internal/usecases/cryptor/private.key" -t "192.168.31.1/24"

.PHONY: runAgent
runAgent: build
	./cmd/agent/agent -a ":$(port)" -r 10 -p 2 -k "superKey" -l 1 -crypto-key "./internal/usecases/cryptor/public.key" -ag ":3200"

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: build
build:
	go build $(buildFlag) -o ./cmd/server/ ./cmd/server/
	chmod +x ./cmd/server
	go build $(buildFlag) -o ./cmd/agent/ ./cmd/agent/
	chmod +x ./cmd/agent

.PHONY: mytest
mytest:
	go test -v -count=1 ./...


# Профилирование
.PHONY: getBasePprof
getBasePprof: build
	cd ./internal/handlers && go test -bench ./... -memprofile=../../cmd/server/profiles/base.pprof
	cd ./cmd/agent && go test -bench ./... -memprofile=./profiles/base.pprof

.PHONY: getResultPprof
getResultPprof: build
	cd ./internal/handlers && go test -bench ./... -memprofile=../../cmd/server/profiles/result.pprof
	cd ./cmd/agent && go test -bench ./... -memprofile=./profiles/result.pprof

.PHONY: showServerBaseProfile
showServerProfile:
	go tool pprof -http=":9090" ./cmd/server/profiles/base.pprof ./cmd/server/server

.PHONY: showServerResultProfile
showServerResultProfile:
	go tool pprof -http=":9090" ./cmd/server/profiles/result.pprof ./cmd/server/server

.PHONY: showAgentBaseProfile
showAgentBaseProfile:
	go tool pprof -http=":9090" ./cmd/agent/profiles/base.pprof ./cmd/agent/agent

.PHONY: showAgentResultProfile
showAgentResultProfile:
	go tool pprof -http=":9090" ./cmd/agent/profiles/result.pprof ./cmd/agent/agent

.PHONY: genProto
genProto:
	protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      proto/server.proto