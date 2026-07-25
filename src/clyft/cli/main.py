import argparse

from clyft.service.artifact import artifact
from clyft.service.init import init


def _build_artifact_parser(subparser) -> None:
    artifact_parser = subparser.add_parser("artifact", help="Create local OCI layout artifact.")
    artifact_parser.add_argument("-t", "--tag", help="OCI artifact tag", required=True)
    artifact_parser.add_argument(
        "-p",
        "--path",
        action="append",
        help="Directory or file path. \
            When file is supplied \
            it's packaged in the root of OCI artifact.",
        required=True,
    )
    artifact_parser.set_defaults(func=_artifact_callback, parser=artifact_parser)


def _artifact_callback(args: argparse.Namespace) -> None:
    artifact(args.tag, args.path)


def _build_init_parser(subparser) -> None:
    init_parser = subparser.add_parser(
        "init", help="Initialise clyft directory in /var/lib/clyft/ or ~/.local/share/clyft/"
    )
    init_parser.set_defaults(func=_init_callback, parser=init_parser)


def _init_callback(args: argparse.Namespace) -> None:  # noqa: ARG001
    print(f"initialized clyft storage at {init()}")


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("clyft")
    subparser = parser.add_subparsers(dest="command")
    _build_artifact_parser(subparser)
    _build_init_parser(subparser)
    return parser
