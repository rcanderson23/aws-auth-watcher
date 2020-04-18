GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
BIN_NAME=aws-auth-watcher

build:
	$(GOBUILD) -o $(BIN_NAME) -v

standalone:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -ldflags '-w -extldflags "-static"' -o $(BIN_NAME) -v
