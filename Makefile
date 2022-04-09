PROJECT=bullet-journal
BINARY=bj

VERSION=${RELEASE_TAG}

ifndef RELEASE_TAG
	VERSION="master"
endif

TARBALL=${BINARY}-${VERSION}.tar.gz
BUILD_DIR=".build"
OUTPUT_DIR="./bin/"

FILES		?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

build:
	@go build -o bin/bj .

install:
	@rm -rf /opt/ep/ep
	@cp -Rf  ./bin/. /opt/ep

test:
	@go test ./...
	@gofmt -l .
	[ "`gofmt -l $(FILES)`" = "" ]

fmt: ## format the go source files
	@go fmt ./...
	@goimports -w $(FILES)

print-version:
	@echo ${VERSION}

print-files:
	@echo ${FILES}
