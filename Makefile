
.PHONY: help
help:
	@echo "Available commands"
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | sort

.PHONY: tests
.PHONY: test
tests test: # Runs unit tests
	go test

.PHONY: lint
lint: # Run the linter and auto-fix issues where possible
	golangci-lint run --fix

