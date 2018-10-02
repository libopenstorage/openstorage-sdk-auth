FROM alpine
MAINTAINER luis@portworx.com

ADD ./_tmp/openstorage-sdk-auth /
ENTRYPOINT ["/openstorage-sdk-auth"]
