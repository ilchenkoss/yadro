deps:
	@go get ./...
grpc_gen:
	@protoc \
    		--proto_path pkg/proto\
    		 --go_out=./pkg/proto/gen\
    		  --go_opt=paths=source_relative\
    		   --go-grpc_out=./pkg/proto/gen\
    			--go-grpc_opt=paths=source_relative \
    			pkg/proto/*.proto
xkcd: deps
	@go build -o xkcd-server ./cmd/xkcd-server
auth: deps
	@go build -o auth-server ./cmd/auth-server
web: deps
	@go build -o web-server ./cmd/web-server
servers_start:grpc_gen xkcd web auth
	@sh servers_start.sh

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
	@./bin/golangci-lint run -v
sec: deps
	@echo "Running Security Checks"
	@trivy fs . --scanners vuln
	@govulncheck ./...

e2e: auth xkcd
	@sh ./e2e_test.sh

make_all: test lint sec servers_start