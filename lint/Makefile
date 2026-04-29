.PHONY: lint

lint:
	GOCACHE=/tmp/gocache GOLANGCI_LINT_CACHE=/tmp/golangci-lint golangci-lint run --config ./.golangci.yaml ./...
