deps:
	@go get ./...

server: deps
	@go build -o xkcd-server ./cmd/xkcd-server

web: deps
	@go build -o web-server ./cmd/web-server

test: deps
	@echo "Running Tests"
	@mkdir -p coverage
	@go test -race -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@go test -v ./...
lint: deps
	@echo "Running Linting and Vetting"
	@gofmt -l .
	@go vet -v ./...
	@sh golangci-lint_install.sh
	@./bin/golangci-lint run -v
sec: deps
	@echo "Running Security Checks"
	@trivy fs . --scanners vuln
	@govulncheck ./...

e2e: server
	@sudo sh ./e2e_test.sh