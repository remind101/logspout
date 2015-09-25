FROM golang:1.5.1
MAINTAINER Ben Guillet <beng@remind101.com>

ENV GOPATH=/go
ENTRYPOINT ["/bin/logspout"]
VOLUME /mnt/routes

RUN go get -u github.com/tools/godep

COPY ./ /go/src/github.com/remind101/logspout
RUN cd /go/src/github.com/remind101/logspout && \
      /go/bin/godep go install ./ && \
      mv /go/bin/logspout /bin/logspout
