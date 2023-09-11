test:
	go clean -testcache && go test -race -cover ./tests

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run --out-format=github-actions ./...

.PHONY: install-lint-deps test lint