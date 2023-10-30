.DEFAUTL_GOAL := info
include .env
export

info:
	@echo "To make teminel run make build"

build:
	mkdir -p .build
	go build $(BUILD_FLAGS) -o .build/$(BUILD_BIN_NAME) .

compress: build
	@upx --brute .build/$(BUILD_BIN_NAME)
