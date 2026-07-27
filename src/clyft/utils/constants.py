from pathlib import Path
from typing import Final

APP_NAME: Final[str] = "clyft"

CLYFT_USER_STORAGE_PATH: Final[Path] = Path.home() / ".local" / "share" / "clyft"
CLYFT_ROOT_STORAGE_PATH: Final[Path] = Path("/var/lib/clyft")
