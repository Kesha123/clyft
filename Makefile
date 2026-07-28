.PHONY: build vet tidy vendor

build:
	go build -o bin/

vet:
	go vet

tidy:
	go mod tidy

vendor:
	go mod vendor
