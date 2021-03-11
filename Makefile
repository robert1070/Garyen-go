# Makefile
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "GitNotFound")
BUILD_DATE=$(shell date '+%Y-%m-%d %H:%M:%S')

default: build

build: clean
	go build -ldflags "-X \"main.BuildVersion=${COMMIT_HASH}\" -X \"main.BuildDate=$(BUILD_DATE)\"" -race -o ./bin/garen .

build-linux: clean
	GOOS=linux GOARCH=amd64 go build -o ./bin/garen .

clean:
	@rm -rf bin/garyen