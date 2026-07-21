import sys
import argparse
from clyft.service.artifact import artifact

def _build_artifact_parser(subparser) -> None:
    artifact_parser = subparser.add_parser('artifact')
    artifact_parser.add_argument('-t', '--tag', help='OCI artifact tag', required=True)
    artifact_parser.add_argument(
        '-p',
        '--path',
        action='append',
        help='Directory or file path. \
            When file is supplied \
            it\'s packaged in the root of OCI artifact.',
        required=True
    )
    artifact_parser.set_defaults(func=_artifact_callback, parser=artifact_parser)


def _artifact_callback(args: argparse.Namespace) -> None:
    artifact(args.tag, args.path)


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser('clyft')
    subparser = parser.add_subparsers(dest='command')
    _build_artifact_parser(subparser)
    return parser
