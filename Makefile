all: clean build

clean:
	rm -rf bin

linux:
	go build -o bin/linux/amd64/lcectl

mac:
	go build -o bin/darwin/amd64/lcectl

windows:
	go build -o bin/windows/amd64/lcectl.exe

build: linux mac windows
