FROM golang:1.6.3
MAINTAINER Ben Guillet <beng@remind101.com>

COPY ./ /go/src/github.com/remind101/logspout
RUN cd /go/src/github.com/remind101/logspout && \
      go install ./ && \
      mv /go/bin/logspout /bin/logspout

ENTRYPOINT ["/bin/logspout"]
VOLUME /mnt/routes
