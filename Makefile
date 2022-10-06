ifeq ($(OS),Windows_NT)
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=if exist bin rd /s /q bin

	VERSION=$(shell git describe --match 'v[0-9]*' --dirty=.m --always --tags || echo "unknown-version")
	GO_LDFLAGS="-X 'liferay.com/lcectl/cmd.Version=$(VERSION)'"

	GIT_GO_PATCH=patches\issues-305.patch

	INSTALL_DEPS=windows
	INSTALL_SRC=bin\windows\amd64\lcectl.exe
	INSTALL_CMD=@echo off & echo ==== Copy $(INSTALL_SRC) onto your %%PATH%%, or into %windir%
else
	GO_CMD_WRAPPER=./gow
	RM_CMD=rm -rf bin

	VERSION=$(shell git describe --match 'v[0-9]*' --dirty='.m' --always --tags | sed 's/^v//' 2>/dev/null || echo "unknown-version")
	GO_LDFLAGS="-X 'liferay.com/lcectl/cmd.Version=$(VERSION)'"

	GIT_GO_PATCH=patches\issues-305.patch

	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		INSTALL_DEPS=mac
		INSTALL_SRC=bin/darwin/amd64/lcectl
		INSTALL_CMD=cp $(INSTALL_SRC) /usr/local/bin/lcectl
	else ifeq ($(UNAME_S),Linux)
		INSTALL_DEPS=linux
		INSTALL_SRC=bin/linux/amd64/lcectl
		INSTALL_CMD=cp $(INSTALL_SRC) /usr/local/bin/lcectl
	endif
endif

all: clean build

patches:
	git apply $(GIT_GO_PATCH) --directory=vendor/github.com/go-git/go-git/v5

clean:
	$(RM_CMD)

linux: patches
	GOOS=linux GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/linux/amd64/lcectl

mac: patches
	GOOS=darwin GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/amd64/lcectl

mac_m1: patches
	GOOS=darwin GOARCH=arm64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/darwin/arm64/lcectl

windows: patches
	GOOS=windows GOARCH=amd64 $(GO_CMD_WRAPPER) build -ldflags=$(GO_LDFLAGS) -o bin/windows/amd64/lcectl.exe

build: linux mac mac_m1 windows

install: $(INSTALL_DEPS)
	$(INSTALL_CMD)
