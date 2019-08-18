# Go parameters
GOARCH=amd64
GOBUILD=go build -o
GOTEST=go test -v -race ./...
GOLINT=golint ./...
GOCLEAN=go clean ./...
GOINSTALL=go install cmd/*

# Build parameters
RABBIT_MSG_VER=1.0.0
RABBIT_MSG_LINUX=env GOOS=linux GOARCH=$(GOARCH)
RABBIT_MSG_MAC=env GOOS=darwin GOARCH=$(GOARCH)
RABBIT_MSG_WIN=env GOOS=windows GOARCH=$(GOARCH)

# go-yelp parameters
RABBIT_MSG_DIR=./cmd/
RABBIT_MSG_BIN=rabbit_msg

# go-yelp build for macOS and Linux
LINUX=$(RABBIT_MSG_LINUX) $(GOBUILD) $(RABBIT_MSG_BIN)_v$(RABBIT_MSG_VER)-linux
MAC=$(RABBIT_MSG_MAC) $(GOBUILD) $(RABBIT_MSG_BIN)_v$(RABBIT_MSG_VER)-mac
WIN=$(RABBIT_MSG_WIN) $(GOBUILD) $(RABBIT_MSG_BIN)_v$(RABBIT_MSG_VER)-win

RABBIT_MSG_BUILD_LINUX=$(LINUX) $(RABBIT_MSG_DIR)
RABBIT_MSG_BUILD_MAC=$(MAC) $(RABBIT_MSG_DIR)
RABBIT_MSG_BUILD_WIN=$(WIN) $(RABBIT_MSG_DIR)

all: clean test build install
build: build-linux build-mac build-win
build-mac:
	$(RABBIT_MSG_BUILD_MAC)
build-linux:
	$(RABBIT_MSG_BUILD_LINUX)
build-win:
	$(RABBIT_MSG_BUILD_WIN)
install:
	$(GOINSTALL)
test:
	$(GOTEST)
lint:
	$(GOLINT)
clean:
	$(GOCLEAN)
	rm -f rabbit_msg_v*