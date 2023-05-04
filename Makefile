BINARY_NAME=cli

ifeq ($(OS),Windows_NT)
	export GOBIN=$(USERPROFILE)\go\bin
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=if exist bin rd /s /q bin

	VERSION=$(shell git describe --match 'v[0-9\.]*' --dirty=.m --always --tags || echo "unknown-version")

	GIT_GO_PATCH_1=patches\issues-305.patch --directory=vendor/github.com/go-git/go-git/v5
	GIT_GO_PATCH_2=patches\issues-gitgo2.patch
else
	export GOBIN=$(HOME)/go/bin
	GO_CMD_WRAPPER=./gow
	RM_CMD=rm -rf bin

	VERSION=$(shell git describe --match 'v[0-9\.]*' --dirty='.m' --always --tags | sed 's/^v//' 2>/dev/null || echo "unknown-version")

	GIT_GO_PATCH_1=patches/issues-305.patch --directory=vendor/github.com/go-git/go-git/v5
	GIT_GO_PATCH_2=patches/issues-gitgo2.patch
endif

GO_LDFLAGS="-X 'liferay.com/liferay/cli/cmd.Version=$(VERSION)'"

.PHONY: all clean patches

all: clean build

clean:
	$(GO_CMD_WRAPPER) clean
	$(RM_CMD)

goenv:
	$(GO_CMD_WRAPPER) env

patches:
	-git apply $(GIT_GO_PATCH_1)
	-git apply $(GIT_GO_PATCH_2)

test:
	$(GO_CMD_WRAPPER) test ./...

testv:
	$(GO_CMD_WRAPPER) test -v ./...

linux: export GOOS=linux
linux: export GOARCH=amd64
linux: patches
	$(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/linux/amd64/liferay

mac: export GOOS=darwin
mac: export GOARCH=amd64
mac: patches
	$(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/amd64/liferay

mac_m1: export GOOS=darwin
mac_m1: export GOARCH=arm64
mac_m1: patches
	$(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/arm64/liferay

windows: export GOOS=windows
windows: export GOARCH=amd64
windows: patches
	$(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/windows/amd64/liferay.exe

build: linux mac mac_m1 windows

install:
	$(GO_CMD_WRAPPER) install -ldflags=$(GO_LDFLAGS)
	@echo "🤖 $(BINARY_NAME) is installed to $(GOBIN)"
	@echo "🤖 Make sure $(GOBIN) is in your PATH environment variable"
	@echo "🤖 Try running: cli --version"
