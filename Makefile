deps:
	@go get ./...
xkcd: deps
	@protoc \
    		--proto_path internal-xkcd/adapters/grpc/proto\
    		 --go_out=./internal-xkcd/adapters/grpc/proto/gen\
    		  --go_opt=paths=source_relative\
    		   --go-grpc_out=./internal-xkcd/adapters/grpc/proto/gen\
    			--go-grpc_opt=paths=source_relative \
    			internal-xkcd/adapters/grpc/proto/*.proto
	@go build -o xkcd-server ./cmd/xkcd-server

auth: deps
	@protoc \
		--proto_path internal-auth/adapters/grpc/proto\
		 --go_out=./internal-auth/adapters/grpc/proto/gen\
		  --go_opt=paths=source_relative\
		   --go-grpc_out=./internal-auth/adapters/grpc/proto/gen\
			--go-grpc_opt=paths=source_relative \
			internal-auth/adapters/grpc/proto/*.proto
	@go build -o auth-server ./cmd/auth-server

web: deps
	@go build -o web-server ./cmd/web-server
servers_start:xkcd web auth
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
	@sh golangci-lint_install.sh
	@./bin/golangci-lint run -v
sec: deps
	@echo "Running Security Checks"
	@trivy fs . --scanners vuln
	@govulncheck ./...
e2e: auth xkcd
	@sh ./e2e_test.sh