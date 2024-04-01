build:
	@go build -o xkcd ./cmd/xkcd
deps:
	@go get ./pkg/xkcd
	@go get ./pkg/database
	@go get ./pkg/words
test:
	@echo "Running Tests"
	@go test -v ./...