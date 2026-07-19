import sys
import argparse
from clyft.cli.main import build_parser

def main() -> None:
    parser = build_parser()
    args: argparse.Namespace = parser.parse_args()
    if len(sys.argv) == 1:
        parser.print_help(sys.stderr)
        sys.exit(1)
    try:
        args.func(args)
    except Exception as exception:
        print(exception, file=sys.stderr)
        sys.exit(1)

if __name__ == '__main__':
    main()
