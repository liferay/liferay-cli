ifeq ($(OS),Windows_NT)
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=if exist bin rd /s /q bin

	VERSION=$(shell git describe --match 'v[0-9\.]*' --dirty=.m --always --tags || echo "unknown-version")
	GO_LDFLAGS="-X 'liferay.com/liferay/cli/cmd.Version=$(VERSION)'"

	GIT_GO_PATCH_1=patches\issues-305.patch --directory=vendor/github.com/go-git/go-git/v5
	GIT_GO_PATCH_2=patches\issues-gitgo2.patch

	INSTALL_DEPS=windows
	INSTALL_SRC=bin\windows\amd64\liferay.exe
	INSTALL_CMD=copy /B /Y $(INSTALL_SRC) $(GOPATH)\bin
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

.PHONY: all clean patches

all: clean build

patches:
	-git apply $(GIT_GO_PATCH_1)
	-git apply $(GIT_GO_PATCH_2)

clean:
	$(RM_CMD)

test:
	$(GO_CMD_WRAPPER) test ./docker

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

install: $(INSTALL_DEPS)
	$(INSTALL_CMD)
