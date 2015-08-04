.PHONY: build run

build:
	docker build -t remind101/logspout .

run:
	docker run --rm \
		-e LOGSPOUT=ignore \
		--name="logspout" \
		--volume=/var/run/docker.sock:/var/run/docker.sock \
		remind101/logspout kinesis://

