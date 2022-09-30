ifeq ($(OS),Windows_NT)
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=if exist bin rd /s /q bin
else
	GO_CMD_WRAPPER=./gow
	RM_CMD=rm -rf bin
endif

UNAME_S := $(shell uname -s)

ifeq ($(OS),Windows_NT)
	INSTALL_DEPS=windows
	INSTALL_SRC=bin\windows\amd64\lcectl.exe
	INSTALL_CMD=copy $(INSTALL_SRC) "C:\Windows\" /V
else ifeq ($(UNAME_S),Darwin)
	INSTALL_DEPS=mac
	INSTALL_SRC=bin/darwin/amd64/lcectl
	INSTALL_CMD=cp $(INSTALL_SRC) /usr/local/bin/lcectl
else ifeq ($(UNAME_S),Linux)
	INSTALL_DEPS=linux
	INSTALL_SRC=bin/linux/amd64/lcectl
	INSTALL_CMD=cp $(INSTALL_SRC) /usr/local/bin/lcectl
endif

all: clean build

clean:
	$(RM_CMD)

linux:
	$(GO_CMD_WRAPPER) build -o bin/linux/amd64/lcectl

mac:
	$(GO_CMD_WRAPPER) build -o bin/darwin/amd64/lcectl

windows:
	$(GO_CMD_WRAPPER) build -o bin/windows/amd64/lcectl.exe

build: linux mac windows

install: $(INSTALL_DEPS)
	$(INSTALL_CMD)
