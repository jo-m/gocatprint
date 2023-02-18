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

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"
	find . \( -name '*.c' -or -name '*.h' \) -exec clang-format-10 --style=file --dry-run --Werror {} +
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks all,-ST1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: lint test

format:
	gofmt -w -s .
	find . \( -name '*.c' -or -name '*.h' \) -exec clang-format-10 --style=file -i {} +
