.PHONY: all
all: dependencies fmt linter test
.PHONY: dependencies
dependencies:
	@echo "=> Executing go mod tidy for ensure dependencies..."
	@go mod tidy
.PHONY: fmt
fmt:
	@echo "=> Executing go fmt..."
	@go fmt ./...
.PHONY: linter
linter:
	@echo "=> Executing golangci-lint"
	@golangci-lint run ./...
.PHONY: test
test:
	@echo "=> Running tests"
	@go test ./... -covermode=atomic -coverpkg=./... -count=1 -race;\
	exit_code=$$?;\
 	exit $$exit_code
