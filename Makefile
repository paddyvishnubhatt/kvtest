GOARCH := amd64
GOOS := linux

all: local
local: 
	go build -o kvserver main.go
clean: ## Remove temporary files
	rm -f kvserver
	rm -rf snapshots
	rm -f *.dat
	rm -f *.db
	go clean