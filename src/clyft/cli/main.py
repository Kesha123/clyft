import argparse
from clyft.service.artifact.main import artifact

def _build_artifact_parser(subparser):
    artifact_parser = subparser.add_parser('artifact')
    artifact_parser.add_argument('-t', '--tag', help='OCI artifact tag')
    artifact_parser.add_argument(
        '-ap',
        '--add-path',
        action='append',
        help='Directory or file path. \
            When file is supplied \
            it\'s packaged in the root of OCI artifact.'
    )
    artifact_parser.set_defaults(func=artifact, parser=artifact_parser)

def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser('clyft')
    subparser = parser.add_subparsers(dest='command')
    _build_artifact_parser(subparser)
    return parser
