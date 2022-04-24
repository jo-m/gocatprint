.PHONY: build run test format lint check

build:
	go build -o print ./cmd/print/main.go

run: build
	sudo setcap 'cap_net_raw,cap_net_admin+eip' print
	sudo ./print \
		--log-level=info \
		--hci-device 2 \
		--printer-name GB03 \
		pkg/printer/testdata/test.png

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
