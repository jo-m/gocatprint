.PHONY: print test

print: cmd/print/main.go
	go build -o print ./cmd/print/main.go
	sudo setcap 'cap_net_raw,cap_net_admin+eip' print
	./print --log-level=trace

test:
	go test -v -race ./...
