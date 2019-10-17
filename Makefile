.PHONY: all openstorage-sdk-auth clean test

REV=$(shell git describe --long --tags --match='v*' --dirty)

all: openstorage-sdk-auth

install:
	go install github.com/libopenstorage/openstorage-sdk-auth/cmd/openstorage-sdk-auth

openstorage-sdk-auth:
	CGO_ENABLED=0 GOOS=linux go build \
		-a -ldflags '-X main.version=$(REV) -extldflags "-static"' \
		-o openstorage-sdk-auth ./cmd/openstorage-sdk-auth

clean:
	rm -rf openstorage-sdk-auth

