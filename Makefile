ifeq ($(OS),Windows_NT)
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=if exist bin rd /s /q bin

	VERSION=$(shell git describe --match 'v[0-9\.]*' --dirty=.m --always --tags || echo "unknown-version")
	GO_LDFLAGS="-X 'liferay.com/liferay/cli/cmd.Version=$(VERSION)'"

	GIT_GO_PATCH_1=patches\issues-305.patch --directory=vendor/github.com/go-git/go-git/v5
	GIT_GO_PATCH_2=patches\issues-gitgo2.patch

	INSTALL_DEPS=windows
	INSTALL_SRC=bin\windows\amd64\liferay.exe
	INSTALL_CMD=@echo off & echo ==== Copy $(INSTALL_SRC) onto your %%PATH%%, or into %windir%
else
	GO_CMD_WRAPPER=./gow
	RM_CMD=rm -rf bin

	VERSION=$(shell git describe --match 'v[0-9\.]*' --dirty='.m' --always --tags | sed 's/^v//' 2>/dev/null || echo "unknown-version")
	GO_LDFLAGS="-X 'liferay.com/liferay/cli/cmd.Version=$(VERSION)'"

	GIT_GO_PATCH_1=patches/issues-305.patch --directory=vendor/github.com/go-git/go-git/v5
	GIT_GO_PATCH_2=patches/issues-gitgo2.patch

	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		INSTALL_DEPS=mac
		INSTALL_SRC=bin/darwin/amd64/liferay
		INSTALL_CMD=cp -f $(INSTALL_SRC) $(GOPATH)/bin
	else ifeq ($(UNAME_S),Linux)
		INSTALL_DEPS=linux
		INSTALL_SRC=bin/linux/amd64/liferay
		INSTALL_CMD=cp -f $(INSTALL_SRC) $(GOPATH)/bin
	endif
endif

all: clean build

.PHONY: patches

patches:
	-git apply $(GIT_GO_PATCH_1)
	-git apply $(GIT_GO_PATCH_2)

clean:
	$(RM_CMD)

linux: patches
	GOOS=linux GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/linux/amd64/liferay

mac: patches
	GOOS=darwin GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/amd64/liferay

mac_m1: patches
	GOOS=darwin GOARCH=arm64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/arm64/liferay

windows: patches
	GOOS=windows GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/windows/amd64/liferay.exe

build: linux mac mac_m1 windows

install: $(INSTALL_DEPS)
	$(INSTALL_CMD)
