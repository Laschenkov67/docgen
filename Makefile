.PHONY: test lint vet fmt cover bench tidy ci

GO ?= go

test:
	$(GO) test ./... -race -count=1

cover:
	$(GO) test ./... -race -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -func=coverage.out | tail -1

bench:
	$(GO) test ./... -bench=. -benchmem -run=^$$

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -s -w .

tidy:
	$(GO) mod tidy

ci: tidy fmt vet lint test