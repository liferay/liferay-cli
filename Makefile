ifeq ($(OS),Windows_NT)
	GO_CMD_WRAPPER=gow.cmd
	RM_CMD=rd /s /q
else
	GO_CMD_WRAPPER=./gow
	RM_CMD=rm -rf
endif

all: clean build

clean:
	$(RM_CMD) ./bin

linux:
	$(GO_CMD_WRAPPER) build -o bin/linux/amd64/lcectl

mac:
	$(GO_CMD_WRAPPER) build -o bin/darwin/amd64/lcectl

windows:
	$(GO_CMD_WRAPPER) build -o bin/windows/amd64/lcectl.exe

build: linux mac windows
