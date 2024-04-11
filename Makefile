build: deps
	@go build -o xkcd ./cmd/xkcd

test: deps
	@echo "Running Tests"
	@go test -v ./...

deps:
	@go get ./...