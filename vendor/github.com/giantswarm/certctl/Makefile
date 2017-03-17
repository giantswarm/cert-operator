PROJECT=certctl

BUILD_PATH := $(shell pwd)/.gobuild

PROJECT_PATH := "$(BUILD_PATH)/src/github.com/giantswarm"

BIN=$(PROJECT)

.PHONY: clean get-deps deps run-tests fmt install

GOPATH := $(BUILD_PATH)

SOURCE=$(shell find . -name '*.go')
GOVERSION=1.6.2
VERSION=$(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif
ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

all: get-deps $(BIN)

ci: clean all run-tests

clean:
	rm -rf $(BUILD_PATH) $(BIN)

get-deps: .gobuild

deps:
	@${MAKE} -B -s .gobuild

.gobuild:
	@mkdir -p $(PROJECT_PATH)
	@rm -f $(PROJECT_PATH)/$(PROJECT) && cd "$(PROJECT_PATH)" && ln -s ../../../.. $(PROJECT)
	#
	# Pin and fetch private dependencies.
	@builder get dep -b v0.6.0 git@github.com:hashicorp/vault.git $(BUILD_PATH)/src/github.com/hashicorp/vault
	#
	# Fetch public dependencies via `go get`
	GOPATH=$(GOPATH) go get -d -v github.com/giantswarm/$(PROJECT)

$(BIN): VERSION $(SOURCE)
	echo Building for $(GOOS)/$(GOARCH)
	docker run \
		--rm \
		-v $(shell pwd):/usr/code \
		-e GOPATH=/usr/code/.gobuild \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-w /usr/code \
		golang:$(GOVERSION) \
		go build -a -ldflags " \
			-X github.com/giantswarm/certctl/cli.version=$(VERSION) \
			-X github.com/giantswarm/certctl/cli.goVersion=$(GOVERSION) \
			-X github.com/giantswarm/certctl/cli.gitCommit=$(COMMIT) \
			-X github.com/giantswarm/certctl/cli.osArch=$(GOOS)/$(GOARCH)\
		" \
		-o $(BIN)


run-tests:
	GOPATH=$(GOPATH) go test -v ./...

fmt:
	gofmt -l -w .

install: $(BIN)
	cp $(BIN) /usr/local/bin/
