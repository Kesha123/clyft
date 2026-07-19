import sys
import argparse


def artifact(args: argparse.Namespace) -> None:
    if not args.tag or not args.add_path:
        print("error: at least one of --tag or --add-path is required",
              file=sys.stderr)
        args.parser.print_help(sys.stderr)
        sys.exit(1)
    print(args)
