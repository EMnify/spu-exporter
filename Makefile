SHELL := /usr/bin/env bash
NAME := spu-exporter
IMPORT := github.com/EMnify/$(NAME)
BIN := bin
DIST := dist
GO := go
EXECUTABLE := $(NAME)

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)
GENERATE ?= $(PACKAGES)

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

ifndef VERSION
	VERSION ?= $(shell git rev-parse --short HEAD)
endif

ifndef REVISION
	REVISION ?= $(shell git rev-parse --short HEAD)
endif

LDFLAGS += -s -w
LDFLAGS += -X "main.Version=$(VERSION)"
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -X "main.Revision=$(REVISION)"

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf $(BIN)/
	rm -rf $(DIST)/

.PHONY: sync
sync:
	$(GO) mod download

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golangci-lint > /dev/null; if [ $$? -ne 0 ]; then \
		(echo "please install golangci-lint"; exit 1) \
	fi
	golangci-lint run -v

.PHONY: test
test:
	GOBIN="$(PWD)" $(GO) install github.com/haya14busa/goverage@latest
	./goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: build
build: $(BIN)/$(EXECUTABLE)

$(BIN)/$(EXECUTABLE): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ .

.PHONY: release
release: release-dirs release-build release-checksums

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)

.PHONY: release-build
release-build:
	GOBIN="$(PWD)" $(GO) install github.com/mitchellh/gox@latest
	./gox  -os="linux darwin" -arch="amd64" -verbose -ldflags '-w $(LDFLAGS)' -output="$(DIST)/$(EXECUTABLE)-{{.OS}}-{{.Arch}}" .

.PHONY: release-checksums
release-checksums:
	cd $(DIST); $(foreach file, $(wildcard $(DIST)/$(EXECUTABLE)-*), sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)