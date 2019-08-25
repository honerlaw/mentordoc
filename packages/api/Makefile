.PHONY: test build

test:
	go clean -testcache
	go test ./server/... -it

build:
	go build -o bin/server ./src/server/main.go
