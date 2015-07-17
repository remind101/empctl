.PHONY: cmd test build

test:
	godep go test ./...

cmd:
	godep go build -o build/empctl github.com/remind101/empctl/cmd/empctl

build:
	docker build -t remind101/empctl .
