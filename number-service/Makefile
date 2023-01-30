.PHONY:build
build:
		go build -o main ./cmd/apiserver/

test:
		go test -v -race -timeout 30s ./...
up:
	docker-compose up -d

down:
	docker-compose down

