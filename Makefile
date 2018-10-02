.PHONY: all openstorage-sdk-auth clean test

REGISTRY_NAME=openstorage
IMAGE_NAME=openstorage-sdk-auth
IMAGE_VERSION=latest
IMAGE_TAG=$(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)

#REV=$(shell git describe --long --tags --match='v*' --dirty)
REV=0.1

all: openstorage-sdk-auth

install:
	go install github.com/libopenstorage/openstorage-sdk-auth/cmd/openstorage-sdk-auth

openstorage-sdk-auth:
	rm -rf _tmp
	mkdir -p _tmp
	CGO_ENABLED=0 GOOS=linux go build \
		-a -ldflags '-X main.version=$(REV) -extldflags "-static"' \
		-o ./_tmp/openstorage-sdk-auth ./cmd/openstorage-sdk-auth

clean:
	rm -rf _tmp

container: openstorage-sdk-auth
	docker build -t $(IMAGE_TAG) .

push: container
	docker push $(IMAGE_TAG)

