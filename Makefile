BINARY_NAME=clyft
BUILD_DIR=bin
PLATFORMS=linux/amd64 linux/arm64

.PHONY: build clean vet fmt-check tidy tidy-check vendor test start-zot stop-zot

build:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os_arch=($${platform//\// }); \
		GOOS=$${os_arch[0]}; \
		GOARCH=$${os_arch[1]}; \
		extension=""; \
		output_name="$(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH$$extension"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -mod=vendor -ldflags="-w -s" -o $$output_name .; \
	done

clean:
	@rm -rf $(BUILD_DIR)

vet:
	go vet ./...

fmt-check:
	@out=$$(gofmt -l . | grep -v '^vendor/'); \
	if [ -n "$$out" ]; then echo "$$out"; echo "run: gofmt -w ."; exit 1; fi

tidy:
	go mod tidy

tidy-check: tidy
	git diff --exit-code go.mod go.sum

vendor: tidy-check
	go mod vendor
	go install

test:
	go test ./... -race -count=1

start-zot:
	$(MAKE) -C test start-zot

stop-zot:
	$(MAKE) -C test stop-zot
