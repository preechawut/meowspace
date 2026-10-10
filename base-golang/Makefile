.PHONY: run migrate test build fmt vet db
run:
	go run ./cmd/api
migrate:
	go run ./cmd/migrate
test:
	go test -race ./...
build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
fmt:
	gofmt -w cmd internal migrations
vet:
	go vet ./...
db:
	docker compose up -d --wait
