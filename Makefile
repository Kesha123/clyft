.PHONY: build vet tidy vendor

build:
	go build -o bin/ -mod=vendor

vet:
	go vet ./...

tidy:
	go mod tidy

vendor:
	go mod vendor

start-zot:
	$(MAKE) -C test start-zot

stop-zot:
	$(MAKE) -C test stop-zot
