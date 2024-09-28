.PHONY: build generate

SHELL=/bin/bash -o pipefail
#BUILDCMD=env GOOS=darwin GOARCH=amd64 go build -v
BUILDCMD=env GOOS=linux GOARCH=amd64 go build -v

build: generate
	$(BUILDCMD) -o safe-size-pay cmd/*.go

generate:
	go generate ./...