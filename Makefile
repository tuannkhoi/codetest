.PHONY: check

check:
	@echo "Running golangci-lint..."
	@golangci-lint run
