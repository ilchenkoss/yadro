build: deps
	@go build -o xkcd ./cmd/xkcd

test: deps
	@echo "Running Tests"
	@go test -v ./...

bench:
	@go test -bench= ./pkg/indexing/.
deps:
	@go get ./...