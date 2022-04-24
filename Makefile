.PHONY: build run test format lint check

build:
	go build -o catprint ./cmd/catprint/main.go

run: build
	sudo ./catprint \
		--log-level=trace \
		--hci-device 0 \
		--timeout 10s \
		--no-scale \
		--no-dither \
		pkg/printer/testdata/test.png
	
	sudo ./catprint \
		--log-level=trace \
		--hci-device 0 \
		--timeout 10s \
		pkg/printer/testdata/swan.jpg

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
