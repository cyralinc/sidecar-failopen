-include env
export
test:
	go test ./...

build:
	go build ./cmd/failopen

build-local:
	go build ./cmd/local

run:
	go run ./cmd/local
