# Use ':=' instead of '=' to avoid multiple evaluation of NOW.
# # Substitute problematic characters with underscore using tr,
# #   make doesn't like spaces and ':' in filenames.
#NOW := $(shell date +"%c" | tr ' :' '__')
NOW := $(shell date +"%s" )
UNAME := $(shell uname -s)
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
BUILD_NUMBER_FILE := .buildno
BUILD_NUMBER := $(shell cat $(BUILD_NUMBER_FILE))

default:
	@echo Building native adz binary
	@go vet
	@go build -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"

linux:
	@echo Building Linux binary
	@GOOS="linux" GOARCH="amd64" go build -o adz -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"

osx:
	@echo Building Darwin binary
	@GOOS="darwin" GOARCH="amd64" go build -o adz -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"

release: 
	mkdir -p dist/linux
	mkdir -p dist/darwin
	@echo Building native binary for local testing/use
	@go build -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"
	@echo Building Darwin binary for distribution
	@GOOS="darwin" GOARCH="amd64" go build -o dist/darwin/adz -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"
	@echo Building Linux binary for distribution
	@GOOS="linux" GOARCH="amd64" go build -o dist/linux/adz -ldflags "-X=main.BUILD=$(NOW) -X=main.VERSION=$(BUILD_NUMBER)"
	tar -cvzf adz-$(BUILD_NUMBER).tar.gz dist/*
	ghr  --username postfix-tools --token ${GITHUB_TOKEN} --replace $(BUILD_NUMBER) adz-$(BUILD_NUMBER).tar.gz

# This is useful for when you have a local docker 
# # but no local Go installation
docker-nolocalgo:
	@echo using centurylink/golang-builder to build docker container
	docker pull centurylink/golang-builder 
	docker run --rm -v ${PWD}:/src -v /var/run/docker.sock:/var/run/docker.sock  centurylink/golang-builder

clean:
	@rm -f adz
	@rm -f adz.*.*
	@rm -rf dist/
