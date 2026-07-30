.PHONY: build vet fmt-check tidy tidy-check vendor test start-zot stop-zot

build:
	go build -o bin/ -mod=vendor

vet:
	go vet ./...

fmt-check:
	@out=$$(gofmt -l . | grep -v '^vendor/'); \
	if [ -n "$$out" ]; then echo "$$out"; echo "run: gofmt -w ."; exit 1; fi

tidy:
	go mod tidy

tidy-check: tidy
	git diff --exit-code go.mod go.sum

vendor:
	go mod vendor

test:
	go test ./... -race -count=1

start-zot:
	$(MAKE) -C test start-zot

stop-zot:
	$(MAKE) -C test stop-zot
