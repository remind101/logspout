FROM gliderlabs/alpine:3.1
ENV GOPATH=/go
ENTRYPOINT ["/bin/logspout"]
VOLUME /mnt/routes

RUN apk add --update go git mercurial
RUN go get -u github.com/tools/godep

COPY ./ /go/src/github.com/remind101/logspout
RUN cd /go/src/github.com/remind101/logspout && \
      /go/bin/godep go install ./ && \
      mv /go/bin/logspout /bin/logspout
MAINTAINER Ben Guillet <beng@remind101.com>
