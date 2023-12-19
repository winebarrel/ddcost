.PHONY: all
all: vet test build

.PHONY: build
build:
	go build ./cmd/ddcost

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run
