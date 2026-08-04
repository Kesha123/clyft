#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
context="$(cd "${script_dir}/../.." && pwd)"
sha="$(git -C "${context}" rev-parse HEAD)"

_main() {
	local version="${1:?usage: build.sh <version> [platforms]}"
	local platforms="${2:-linux/amd64,linux/arm64}"
	local container_runtime="${2:-podman}"

	"${container_runtime}" buildx build \
		--platform "${platforms}" \
		--tag "clyft:${version}" \
		--tag "clyft:sha256-${sha}" \
		-f "${script_dir}/Containerfile" "${context}"
}

_main "$@"
