.PHONY: build push up

build:
	@go test ./... | awk 'END{exit(c)}/^ok|^\?/{ next; };/^FAIL/{c=1};1'
	@go run .

up:
	@make -s -C .. up

push:
	@make -s -C .. push
