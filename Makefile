.PHONY: build run test format lint check

build:
	go build -o print ./cmd/print/main.go
	sudo setcap 'cap_net_raw,cap_net_admin+eip' print

run: build
	./print --log-level=trace pkg/printer/testdata/test.png

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
