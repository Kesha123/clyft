#!/usr/bin/env sh
set -eu

script_dir="$(cd "$(dirname "$0")" && pwd)"
context="$(cd "${script_dir}/../.." && pwd)"
sha="$(git -C "${context}" rev-parse HEAD)"
container_runtime="${3:-podman}"

_main() {
	version="${1:?usage: build.sh <version> [platforms] [container_runtime]}"
	platforms="${2:-linux/amd64,linux/arm64}"

	"${container_runtime}" buildx build \
		--platform "${platforms}" \
		--tag "clyft:${version}" \
		--tag "clyft:sha256-${sha}" \
		-f "${script_dir}/Containerfile" "${context}"
}

_main "$@"
