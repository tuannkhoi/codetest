.PHONY: lint test test-integration

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

test:
	@go test ./... -short -cover

test-integration:
	@docker compose up -d
	@go test -tags=integration ./... -cover
