export GO111MODULE=on
TARGET := fluxcloud-filebeat

.PHONY: all build clean fmt

all: build

build:
	@echo "Building $(SRC) to ./bin"
	GOBIN=./bin/$(TARGET) go build -o bin/$(TARGET) cmd/$(TARGET)/main.go 

clean:
	@rm -rf bin/$(TARGET)

fmt:
	@gofmt -l -w $(SRC)
