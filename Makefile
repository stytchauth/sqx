
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

.PHONY: services
services: # Start local development services (databases, etc.)
	docker-compose up -d

.PHONY: services-down
services-down: # _Destroy_ local development services
	docker-compose down
