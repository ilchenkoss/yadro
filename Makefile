build: deps
	@go build -o xkcd ./cmd/xkcd

test: deps
	@echo "Running Tests"
	@go test -v ./...
server: deps
	@go build -o xkcd-server ./cmd/xkcd-server
bench:
	@go test -bench=. ./pkg/indexing/.
deps:
	@go get ./...