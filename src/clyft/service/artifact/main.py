import hashlib
import json
import shutil
from dataclasses import dataclass, field
from datetime import UTC, datetime
from pathlib import Path

_IMAGE_LAYOUT_VERSION = "1.0.0"
_LAYER_MEDIA_TYPE = "application/octet-stream"
_FILE_DIGEST_ALGORITHM = "sha256"


def artifact(tag: str, paths: list[str]) -> None:
    files = _expand_paths(paths)
    index = _build_layout(tag=tag, paths=files)
    _write_layout(index=index, dest=Path(tag))


def _expand_paths(paths: list[str]) -> list[str]:
    files = []
    for path in paths:
        p = Path(path)
        if not p.exists():
            raise FileNotFoundError(f"path does not exist: {path}")
        if p.is_file():
            files.append(path)
            continue
        for root, _dirs, dirs_files in p.walk(follow_symlinks=False):
            for file in dirs_files:
                full_path = root / file
                files.append(str(full_path))
    return files


def _build_layout(tag: str, paths: list[str]) -> OCIIndex:
    layers = [OCILayer(path=p) for p in paths]
    config = OCIConfig()
    manifest = OCIManifest(config=config, layers=layers)
    return OCIIndex(manifests=[manifest], ref_name=tag)


def _write_layout(index: OCIIndex, dest: Path) -> None:
    if dest.exists() and dest.is_dir():
        shutil.rmtree(dest)

    dest.mkdir(parents=True, exist_ok=True)
    blobs_dir = dest / "blobs" / "sha256"
    blobs_dir.mkdir(parents=True, exist_ok=True)

    (dest / "oci-layout").write_text(json.dumps(OCILayout().to_dict(), indent=2))

    for manifest in index.manifests:
        config = manifest.config
        config_digest = config.digest.split(":", 1)[1]
        (blobs_dir / config_digest).write_bytes(config.content)

        for layer in manifest.layers:
            layer_digest = layer.digest.split(":", 1)[1]
            shutil.copy(layer.path, blobs_dir / layer_digest)

        manifest_bytes = manifest.to_json_bytes()
        manifest_digest = manifest.digest.split(":", 1)[1]
        (blobs_dir / manifest_digest).write_bytes(manifest_bytes)

    (dest / "index.json").write_text(json.dumps(index.to_dict(), indent=2))


@dataclass
class OCILayout:
    image_layout_version: str = _IMAGE_LAYOUT_VERSION

    def to_dict(self) -> dict:
        return {"imageLayoutVersion": self.image_layout_version}


@dataclass
class OCILayer:
    path: str
    media_type: str = field(init=False)
    annotations: dict = field(init=False)
    size: int = field(init=False)
    digest: str = field(init=False)

    def __post_init__(self) -> None:
        self.media_type = _LAYER_MEDIA_TYPE
        self.annotations = {
            "org.opencontainers.image.title": Path(self.path).name,
            "org.opencontainers.image.filepath": self.path,
        }
        self.size = Path(self.path).stat().st_size
        with open(self.path, "rb") as file:
            self.digest = f"{_FILE_DIGEST_ALGORITHM}:{hashlib.file_digest(file, _FILE_DIGEST_ALGORITHM).hexdigest()}"

    def to_dict(self) -> dict:
        return {
            "mediaType": self.media_type,
            "digest": self.digest,
            "size": self.size,
            "annotations": self.annotations,
        }


@dataclass
class OCIConfig:
    media_type: str = "application/vnd.oci.empty.v1+json"
    content: bytes = b"{}"
    digest: str = field(init=False)
    size: int = field(init=False)

    def __post_init__(self) -> None:
        self.digest = f"sha256:{hashlib.sha256(self.content).hexdigest()}"
        self.size = len(self.content)

    def to_dict(self) -> dict:
        return {
            "mediaType": self.media_type,
            "digest": self.digest,
            "size": self.size,
        }


@dataclass
class OCIManifest:
    config: OCIConfig
    layers: list[OCILayer]
    schema_version: int = 2
    media_type: str = "application/vnd.oci.image.manifest.v1+json"
    digest: str = field(init=False)
    size: int = field(init=False)

    def __post_init__(self) -> None:
        content = self.to_json_bytes()
        self.digest = f"sha256:{hashlib.sha256(content).hexdigest()}"
        self.size = len(content)

    def to_dict(self) -> dict:
        return {
            "schemaVersion": self.schema_version,
            "mediaType": self.media_type,
            "config": self.config.to_dict(),
            "layers": [layer.to_dict() for layer in self.layers],
        }

    def to_json_bytes(self) -> bytes:
        return json.dumps(self.to_dict()).encode("utf-8")


@dataclass
class OCIIndex:
    manifests: list[OCIManifest]
    ref_name: str
    schema_version: int = 2
    annotations: dict = field(init=False)

    def __post_init__(self) -> None:
        self.annotations = {
            "org.opencontainers.image.ref.name": self.ref_name,
            "org.opencontainers.image.created": datetime.now(UTC).isoformat().replace("+00:00", "Z"),
        }

    def to_dict(self) -> dict:
        return {
            "schemaVersion": self.schema_version,
            "manifests": [
                {
                    "mediaType": manifest.media_type,
                    "digest": manifest.digest,
                    "size": manifest.size,
                    "annotations": self.annotations,
                }
                for manifest in self.manifests
            ],
        }
