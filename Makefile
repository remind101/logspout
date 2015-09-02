NAME=logspout
VERSION=$(shell cat VERSION)

debug:
	docker build -f Dockerfile.debug --no-cache -t remind101/logspout:debug .

run:
	docker run --rm \
	  -it \
		-p 8000:80 \
		-e LOGSPOUT=ignore \
		-e DEBUG=true \
	  --env-file=.env \
		--name="logspout" \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		--volume=$(GOPATH)/src/github.com/remind101/logspout/routes:/mnt/routes \
		remind101/logspout:debug

dev:
	@docker history $(NAME):dev &> /dev/null \
		|| docker build -f Dockerfile.dev -t $(NAME):dev .
	@docker run --rm \
		-e DEBUG=true \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/go/src/github.com/gliderlabs/logspout \
		-p 8000:80 \
		-e ROUTE_URIS=$(ROUTE) \
		$(NAME):dev

build:
	mkdir -p build
	docker build -t $(NAME):$(VERSION) .
	docker save $(NAME):$(VERSION) | gzip -9 > build/$(NAME)_$(VERSION).tgz

release:
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create gliderlabs/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

circleci:
	rm ~/.gitconfig
ifneq ($(CIRCLE_BRANCH), release)
	echo build-$$CIRCLE_BUILD_NUM > VERSION
endif

.PHONY: release debug run build
