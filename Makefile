ARCH ?= amd64
VRP_BIN := ./build/vrp
BUILD_CLI_FOR_SYSTEM := build-cli-linux-amd
ifeq ($(OS),Windows_NT)
	VRP_BIN := $(addsuffix .exe,$(VRP_BIN))
	BUILD_CLI_FOR_SYSTEM := build-cli-windows-amd
else
	UNAME_S := $(shell uname -s)
	UNAME_P := $(shell uname -p)
	ifneq ($(UNAME_S),Linux)
		ifeq ($(UNAME_S),Darwin)
			VRP_BIN := $(addsuffix -mac,$(VRP_BIN))
		endif
		ifeq ($(UNAME_P),i386)
			VRP_BIN := $(addsuffix -intel,$(VRP_BIN))
			BUILD_CLI_FOR_SYSTEM = build-cli-mac-intel
		endif
		ifeq ($(UNAME_P),arm)
			VRP_BIN := $(addsuffix -apple,$(VRP_BIN))
			BUILD_CLI_FOR_SYSTEM = build-cli-mac-apple
		endif
	endif
endif

.PHONY: build
build:
	$(MAKE) $(BUILD_CLI_FOR_SYSTEM)

clean:
	rm -rf $(VRP_BIN)

build-cli-linux-amd: ## Build for Linux on AMD64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/vrp .

build-cli-linux-arm: ## Build for Linux on ARM
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/vrp-arm .

build-cli-mac-intel: ## Build for macOS on AMD64
	GOOS=darwin GOARCH=amd64 go build -o build/vrp-mac-intel .

build-cli-mac-apple: ## Build macOS on ARM
	GOOS=darwin GOARCH=arm64 go build -o build/vrp-mac-apple .

build-cli-windows-amd: ## Build for Windows on AMD64
	GOOS=windows GOARCH=amd64 go build -o build/vrp.exe . ## Build the vrp CLI for Windows on AMD64

build-cli-windows-arm: ## Build for Windows on ARM
	GOOS=windows GOARCH=arm64 go build -o build/vrp-arm.exe . ## Build the vrp CLI for Windows on ARM

build-cli-linux: build-cli-linux-amd build-cli-linux-arm ## Build for Linux on AMD64 and ARM

build-cli: build-cli-linux-amd build-cli-linux-arm build-cli-mac-intel build-cli-mac-apple build-cli-windows-amd build-cli-windows-arm ## Build the CLI

run-tests: build
	python3 evaluateShared.py --cmd "./build/vrp" --problemDir testdata/