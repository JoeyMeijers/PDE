import sys
import json
import csv
from io import StringIO


def run_streaming(input_data):
    """
    Converts CSV data to JSONL format, streaming output directly.
    Input: Raw CSV data as bytes/string
    Output: JSONL data written directly to stdout
    """
    try:
        # Decode if bytes
        if isinstance(input_data, bytes):
            input_data = input_data.decode("utf-8")

        # Stream CSV and output JSONL line by line immediately
        reader = csv.DictReader(StringIO(input_data))

        for row in reader:
            sys.stdout.write(json.dumps(row))
            sys.stdout.write("\n")
            sys.stdout.flush()

    except Exception as e:
        print(f"Error loading data: {e}", file=sys.stderr)


def main():
    input_data = sys.stdin.read()
    run_streaming(input_data)


if __name__ == "__main__":
    main()
