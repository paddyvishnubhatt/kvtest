GOARCH := amd64
GOOS := linux

all: local
local: 
	go build -o kvserver main.go
clean: ## Remove temporary files
	rm -f kvserver
	rm -rf /tmp/raft_dir
	mkdir /tmp/raft_dir
	mkdir /tmp/raft_dir/node{1,2,3}
	go clean