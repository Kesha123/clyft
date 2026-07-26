from pathlib import Path

from clyft.utils import get_clyft_storage_path


def _create_clyft_storage_path() -> None:
    get_clyft_storage_path().mkdir(parents=True, exist_ok=True)


def init() -> Path:
    _create_clyft_storage_path()
    return get_clyft_storage_path()
