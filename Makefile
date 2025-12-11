APP_NAME := loc-counter
BUILD_DIR := bin

.PHONY: all clean build-all build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64

all: build-all

clean:
	rm -rf $(BUILD_DIR)

build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
