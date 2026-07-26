import json
from pathlib import Path

import httpx

from clyft.utils import CLyftRegistryError, CLyftStoragePathNotFoundError, get_clyft_storage_path, validate_tag

USERNAME = ""
PASSWORD = ""

REGISTRY_HOST = ""
REPO_NAME = ""


def _get_manifest_digest_from_index(clyft_storage_path: Path, tag: str) -> str:
    if not clyft_storage_path.exists() or not clyft_storage_path.is_dir():
        raise CLyftStoragePathNotFoundError(clyft_storage_path)

    index_file = clyft_storage_path / tag / "index.json"

    if not index_file.exists():
        raise FileNotFoundError(f"Missing index.json in {tag}")

    with open(index_file, encoding="utf-8") as f:
        index_data = json.load(f)

    for manifest_desc in index_data.get("manifests", []):
        annotations = manifest_desc.get("annotations", {})
        if annotations.get("org.opencontainers.image.ref.name") == tag:
            return manifest_desc["digest"].split(":", 1)[1]

    raise ValueError(f"Tag '{tag}' not found in {index_file}")


def _is_blob_in_registry(client: httpx.Client, registry_url: str, repo: str, blob_hash: str) -> bool:
    blob_digest = f"sha256:{blob_hash}"
    check_url = f"{registry_url}/v2/{repo}/blobs/{blob_digest}"
    head_res = client.head(check_url)
    if head_res.status_code == 404:
        return False
    head_res.raise_for_status()
    return True


def _get_blob_upload_location(client: httpx.Client, registry_url: str, repo: str) -> str:
    upload_initiate_url = f"{registry_url}/v2/{repo}/blobs/uploads/"
    upload_initiate_response = client.post(upload_initiate_url)
    upload_initiate_response.raise_for_status()
    return upload_initiate_response.headers["Location"]


def _upload_blob(
    client: httpx.Client, registry_url: str, upload_location: str, blob_path: Path, blob_hash: str
) -> None:
    print(f"[Uploading Blob] sha256:{blob_hash[:12]}...")
    blob_digest = f"sha256:{blob_hash}"
    if not upload_location.startswith(("http://", "https://")):
        upload_location = f"{registry_url}{upload_location}"
    separator = "&" if "?" in upload_location else "?"
    final_upload_url = f"{upload_location}{separator}digest={blob_digest}"

    with open(blob_path, "rb") as blob_file:
        headers = {"Content-Type": "application/octet-stream"}
        put_res = client.put(final_upload_url, content=blob_file, headers=headers)
        put_res.raise_for_status()


def _upload_blob_if_missing(
    client: httpx.Client, registry_url: str, repo: str, blob_hash: str, blob_path: Path
) -> None:
    if _is_blob_in_registry(client, registry_url, repo, blob_hash):
        return
    upload_location = _get_blob_upload_location(client, registry_url, repo)
    _upload_blob(client, registry_url, upload_location, blob_path, blob_hash)


def push_oci_layout_to_registry(tag: str) -> None:
    validate_tag(tag)
    clyft_storage_path = get_clyft_storage_path()
    manifest_digest = _get_manifest_digest_from_index(clyft_storage_path, tag)
    manifest_blob_path = clyft_storage_path / tag / "blobs" / "sha256" / manifest_digest

    with open(manifest_blob_path, encoding="utf-8") as f:
        manifest_data = json.load(f)

    auth = httpx.BasicAuth(username=USERNAME, password=PASSWORD)

    try:
        with httpx.Client(auth=auth, follow_redirects=True, timeout=30.0) as client:
            print(f"Connected to {REGISTRY_HOST} as {USERNAME}")

            # 1. Upload Config Blob
            config_hash = manifest_data["config"]["digest"].split(":", 1)[1]
            config_path = clyft_storage_path / tag / "blobs" / "sha256" / config_hash
            print("\nProcessing Artifact Config...")
            _upload_blob_if_missing(client, REGISTRY_HOST, REPO_NAME, config_hash, config_path)

            # 2. Upload Layer Blobs (raw YAMLs or binary data)
            print("\nProcessing Layers...")
            for layer in manifest_data.get("layers", []):
                layer_hash = layer["digest"].split(":", 1)[1]
                layer_path = clyft_storage_path / tag / "blobs" / "sha256" / layer_hash
                filename = layer.get("annotations", {}).get("org.opencontainers.image.title", layer_hash[:12])

                print(f"Layer: {filename}")
                _upload_blob_if_missing(client, REGISTRY_HOST, REPO_NAME, layer_hash, layer_path)

            # 3. Upload Manifest JSON (The final transaction step)
            print(f"\nPublishing Manifest for tag '{tag}'...")
            manifest_url = f"{REGISTRY_HOST}/v2/{REPO_NAME}/manifests/{tag}"
            manifest_media_type = manifest_data.get("mediaType", "application/vnd.oci.image.manifest.v1+json")

            headers = {"Content-Type": manifest_media_type}
            manifest_bytes = json.dumps(manifest_data).encode("utf-8")

            res = client.put(manifest_url, content=manifest_bytes, headers=headers)
            res.raise_for_status()

            print(f"SUCCESS! Tag '{tag}' pushed to {REGISTRY_HOST}/{REPO_NAME}:{tag}")
    except httpx.HTTPError as exc:
        raise CLyftRegistryError(f"failed to push tag '{tag}' to {REGISTRY_HOST}: {exc}") from exc
