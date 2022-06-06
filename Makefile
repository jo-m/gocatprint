.PHONY: build run help test format lint check

build:
	go build -o catprint ./cmd/catprint/
	go build -o catprintsimple ./cmd/catprintsimple/
	go build -o catprintcam ./cmd/catprintcam/

run: build
	sudo ./catprint \
		--log-level=info \
		--timeout 10s \
		--threshold \
		pkg/printer/testdata/test.png
	
	sudo ./catprint \
		--log-level=info \
		--timeout 10s \
		pkg/printer/testdata/swan.jpg

help: build
	./catprint --help

test:
	go test -v -race ./...

format:
	gofmt -w .
	go mod tidy

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"

	# go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...
	
	go vet ./...

check: lint test
