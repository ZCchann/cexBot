GOPATH ?= $(shell go env GOPATH)

ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

GO        := GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
GOBUILD   := $(GO) build
GOCLEAN	  := $(GO) clean
GOLINT    := golangci-lint
# LINTCONF  := .golangci.yml
# CONF      := ./conf/config.json
# BIN       := ./bin

$(shell mkdir -p bin)

# 编译时间
COMPILE_TIME = $(shell date +"%Y-%M-%d %H:%M:%S")
# CFLAGS += $(COMPILE_TIME)
# GIT版本号
GIT_REVISION = $(shell git show -s --pretty=format:%h)
CFLAGS += $(GIT_REVISION)

.PHONY: fmt test lint clean
# default: all

fmt:
	$(info ******************** checking formatting ********************)
	gofmt -s -w ./

test:
	$(info ******************** running tests ********************)
	go test $(shell go list ./... | grep -v /test)

lint:
	$(info ******************** running lint tools ********************)
	$(GOLINT) run -v


build:
	@$(GOBUILD) -o ./bin/cexbot ./main.go
	@cp -r ./conf ./bin

pack:
	@rm -f ./cexbot.*.zip
	@zip -qr cexbot.zip bin/*

clean:
	@rm -f cexbot.*.zip
	@rm -rf bin
