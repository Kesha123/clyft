from pathlib import Path


class CLyftStoragePathNotFoundError(FileNotFoundError):
    def __init__(self, storage_path: str | Path):
        super().__init__(f"clyft storage path not found: {storage_path}. run 'clyft init' first.")


class CLyftRegistryError(Exception):
    """Raised when an OCI registry operation fails."""
