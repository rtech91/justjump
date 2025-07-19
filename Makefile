BINARY_NAME=justjump
INSTALL_DIR=/usr/local/bin

build:
	go build -o $(BINARY_NAME) main.go

build-release:
	go build -o $(BINARY_NAME) -ldflags="-s -w" main.go

install:
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

remove:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

.PHONY: build build-release clean install remove
