APP := proxy
CMD := ./cmd/$(APP)
BIN_DIR := ./bin
GO := go

.PHONY: help deps tidy download build run test fmt vet clean

help:
	@echo "make deps        - tidy + download deps"
	@echo "make build       - build ./cmd/$(APP) -> ./bin/$(APP)"
	@echo "make run         - run ./cmd/$(APP)"
	@echo "make test        - run tests (race)"
	@echo "make fmt         - format code"
	@echo "make vet         - static checks"
	@echo "make clean       - remove ./bin"

deps: tidy download

tidy:
	$(GO) mod tidy

download:
	$(GO) mod download

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -o $(BIN_DIR)/$(APP) $(CMD)

run:
	$(GO) run $(CMD)

test:
	$(GO) test ./... -race -count=1

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	rm -rf $(BIN_DIR)