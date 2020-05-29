build:
	go build -o snappy ./cmd/snappy/main.go

build_docker:
	docker build -t czeslavo/snappy .

up:
	docker-compose up -d --build

test:
	go test ./...

test_integration: up
	go test -count=1 -tags=integration ./...

wire:
	wire ./internal/service

lint:
	golangci-lint run -E goimports

TOOLS += github.com/google/wire/cmd/wire@v0.4.0
$(TOOLS): %:
	cd /tmp && GOBIN=$(GOBIN) GO111MODULE=on go get -u $*

tools: $(TOOLS) golangci_lint

GOLANGCI_VERSION=v1.27.0
golangci_lint:
	cd /tmp && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_VERSION}

.PHONY: tools build test test_integration up wire golanci_lint
