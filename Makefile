BINARY_NAME=justjump

# Detect the operating system
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	INSTALL_DIR=/usr/local/bin
endif
ifeq ($(UNAME_S),Darwin)
	INSTALL_DIR=/usr/local/bin
endif

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

build-release:
	go build -o $(BINARY_NAME) -ldflags="-s -w" ./cmd/main.go

build-release-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -ldflags="-s -w" ./cmd/main.go

build-release-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) -ldflags="-s -w" ./cmd/main.go

install:
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

remove:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

.PHONY: build clean install remove build-release-linux build-release-darwin
