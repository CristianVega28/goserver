APP_NAME = goserver
VERSION = 0.0.1

APP-NAME=main
WATCH-FILE=watcher
# OS ?= $(shell go env GOOS)
# ARCH ?= $(shell go env GOARCH)
FILES = $(APP-NAME) $(WATCH-FILE)

UNNAME_S := $(strip $(shell uname -s))
UNNAME_M := $(strip $(shell uname -m))

ifeq ($(UNNAME_S),Darwin)
OS := darwin
endif

ifeq ($(UNNAME_S),Linux)
OS := linux
endif

ifeq ($(UNNAME_S),MINGW64_NT-10.0)
OS := windows
endif

ifeq ($(UNNAME_M),x86_64)
ARCH := amd64
else ifeq ($(UNNAME_M),arm64)
ARCH := arm64
endif
	
print:
	@echo UNNAME_S=$(UNNAME_S)
	@echo UNNAME_M=$(UNNAME_M)
	@echo OS=$(OS)
	@echo ARCH=$(ARCH)

run:
	export GOSERVER_ENV=development && go run watcher.go
build: 
	@echo "Building $(APP_NAME) for $(OS)/$(ARCH)..."

	@for f in $(FILES); do \
		GOOS=$(OS) GOARCH=$(ARCH) go build -o $$f-$(OS)-$(ARCH) $$f.go; \
	done
clear:
	rm -rf ./tmp
	rm -rf ./*.exe
	rm -rf *.db
	rm -rf main watcher
	sudo rm -rf /usr/local/bin/goserver
	sudo rm -rf /usr/local/bin/watcher-$(OS)-$(ARCH)
	sudo rm -rf /usr/local/bin/main-$(OS)-$(ARCH)

prod:
	printf '#!/bin/bash\n' > goserver && \
	printf 'export GOSERVER_ENV=production\n' >> goserver && \
	printf 'cd %s\n' "$(PWD)" >> goserver && \
	printf 'watcher-$(OS)-$(ARCH)\n' >> goserver && \
	chmod +x goserver && \
	sudo mv goserver /usr/local/bin/ && \
	sudo mv ./watcher-$(OS)-$(ARCH) /usr/local/bin/ && \
	sudo mv ./main-$(OS)-$(ARCH) /usr/local/bin/ && \
	echo "$(APP_NAME) version $(VERSION) installed successfully!"


	
	
	
	