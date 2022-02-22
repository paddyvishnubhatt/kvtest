GOARCH := amd64
GOOS := linux

all: local
protoc: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative kv/kv.proto
local: protoc
	go build -o kvserver main.go
	go build -o kvclient client/client.go
clean: ## Remove temporary files
	rm -f kvserver
	rm -rf /tmp/raft_dir
	mkdir /tmp/raft_dir
	mkdir /tmp/raft_dir/node{1,2,3}
	go clean