.PHONY: proto test

proto:
	buf generate

test:
	go test -v ./...
