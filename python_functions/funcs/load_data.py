import sys
import json
import csv
from io import StringIO


def run_streaming(batch_size=10000):
    """
    Reads CSV from stdin line by line and outputs JSONL in batches to stdout.
    This avoids keeping the entire CSV in memory.
    """
    import sys, csv, json

    reader = csv.DictReader(sys.stdin)
    buffer = []
    count = 0

    for row in reader:
        buffer.append(json.dumps(row))
        count += 1
        if count % batch_size == 0:
            sys.stdout.write("\n".join(buffer) + "\n")
            sys.stdout.flush()
            buffer = []

    if buffer:
        sys.stdout.write("\n".join(buffer) + "\n")
        sys.stdout.flush()

    print(
        f"DEBUG: wrote {count} rows in batches of {batch_size}",
        file=sys.stderr,
        flush=True,
    )


def main():
    run_streaming()


if __name__ == "__main__":
    main()
