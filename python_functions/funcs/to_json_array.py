import sys
import json


def run_streaming(input_data: str):
    """
    Converts JSONL (one JSON per line) to a JSON array.
    Input: JSONL formatted data (one JSON per line)
    Output: Single JSON array
    """
    try:
        # Decode if bytes
        if isinstance(input_data, bytes):
            input_data = input_data.decode("utf-8")

        data = []

        for line in input_data.strip().split("\n"):
            if line:
                data.append(json.loads(line))

        # Output as JSON array
        print(json.dumps(data), flush=True)

    except Exception as e:
        print(f"Error in to_json_array: {e}", file=sys.stderr)


def main():
    input_data = sys.stdin.read()
    run_streaming(input_data)


if __name__ == "__main__":
    main()
