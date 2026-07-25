import getpass
import sqlite3
from pathlib import Path

USER_CLYFT_PATH = ".local/share/clyft/"
ROOT_USER_CLYFT_PATH = "/var/lib/clyft/"
CLYFT_DATABASE_NAME = "database.db"


def _get_current_user() -> str:
    return getpass.getuser()


def get_clyft_path() -> str:
    if _get_current_user() == "root":
        return ROOT_USER_CLYFT_PATH
    return str(Path.home() / USER_CLYFT_PATH)


def create_clyft_path() -> None:
    Path(get_clyft_path()).mkdir(parents=True, exist_ok=True)


def create_clyft_database() -> None:
    database_path = Path(get_clyft_path()) / CLYFT_DATABASE_NAME
    connection = sqlite3.connect(database_path)
    try:
        pass
    finally:
        connection.close()


def init() -> str:
    create_clyft_path()
    create_clyft_database()
    return get_clyft_path()
