BUILD_VERSION 		:= v0.0.0
PLATFORM      		?= linux/amd64,linux/arm64
CONTAINER_RUNTIME 	?= podman

.PHONY: build/%

build/%:
	@./$*/build.sh "$(BUILD_VERSION)" "$(PLATFORM)" "$(CONTAINER_RUNTIME)"
