.PHONY: build test

PACKAGES=`go list ./... | grep -v /vendor/`

build:
	go build -o radagast ./cmd/radagast

build-linux:
	GOOS=linux go build -o radagast-linux ./cmd/radagast

test:
	go test -v ${PACKAGES}
