import os
from enum import Enum
from pathlib import Path

from clyft.utils.constants import ROOT_USER_CLYFT_PATH, USER_CLYFT_PATH


class ContainerRuntime(Enum):
    DOCKER = "docker"
    PODMAN = "podman"

    def __str__(self) -> str:
        return self.value


def get_clyft_storage_path() -> Path:
    if os.geteuid() == 0:
        return ROOT_USER_CLYFT_PATH
    return Path.home() / USER_CLYFT_PATH


def validate_tag(tag: str) -> None:
    if not tag or tag in (".", "..") or Path(tag).name != tag:
        raise ValueError(f"invalid tag: {tag!r}")
