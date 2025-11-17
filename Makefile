APP_NAME = goserver
VERSION = 0.0.1

APP-NAME=main
WATCH-FILE=watcher
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
FILES = $(APP-NAME) $(WATCH-FILE)

run:
	go run watcher.go
build: 
	@echo "Building $(APP_NAME) for $(OS)/$(ARCH)..."

	@for f in $(FILES); do \
		GOOS=$(OS) GOARCH=$(ARCH) go build -o $$f $$f.go; \
	done
clear:
	rm  -rf ./tmp
	rm -rf ./*.exe
	
	
	