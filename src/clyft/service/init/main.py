import getpass
from pathlib import Path

USER_CLYFT_PATH = ".local/share/clyft/"
ROOT_USER_CLYFT_PATH = "/var/lib/clyft/"


def _get_current_user() -> str:
    return getpass.getuser()


def get_clyft_path() -> str:
    if _get_current_user() == "root":
        return ROOT_USER_CLYFT_PATH
    return str(Path.home() / USER_CLYFT_PATH)


def create_clyft_path() -> None:
    Path(get_clyft_path()).mkdir(parents=True, exist_ok=True)


def init() -> str:
    create_clyft_path()
    return get_clyft_path()
