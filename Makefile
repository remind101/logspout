.PHONY: build run

build:
	docker build -t remind101/logspout .

test:
	go test -race $(shell go list ./... | grep -v /vendor/)

run:
	docker run --rm \
		-p 8000:80 \
		-e LOGSPOUT=ignore \
		--env-file=.env \
		--name="logspout" \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		--volume=$(GOPATH)/src/github.com/remind101/logspout/routes:/mnt/routes \
		remind101/logspout
