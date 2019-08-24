.PHONY: test run build

test:
	go clean -testcache
	go test ./server/... -it

run:
	docker kill mentordoc-mysql; docker rm mentordoc-mysql; docker run --name mentordoc-mysql -p 3306:3306 -e MYSQL_USER=userlocal -e MYSQL_ROOT_PASSWORD=password -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=mentor_doc -d mysql:5.7
	DATABASE_HOST=0.0.0.0 DATABASE_PORT=3306 DATABASE_USERNAME=userlocal DATABASE_NAME=mentor_doc DATABASE_PASSWORD=password go run ./src/server/main.go

build:
	go build -o bin/server ./src/server/main.go
