import sys
import json


def run_streaming(input_data: json):
    """
    Filters JSONL records with income > 50000.
    Input: JSONL formatted data (one JSON per line)
    Output: Filtered JSONL data
    """
    try:
        # Decode if bytes
        if isinstance(input_data, bytes):
            input_data = input_data.decode("utf-8")

        filtered_count = 0

        for line in input_data.strip().split("\n"):
            if line:
                try:
                    row = json.loads(line)
                    income = float(row.get("income", 0))
                    if income > 50000:
                        sys.stdout.write(json.dumps(row))
                        sys.stdout.write("\n")
                        sys.stdout.flush()
                        filtered_count += 1
                except (json.JSONDecodeError, ValueError):
                    pass

        print(
            f"Filtered to {filtered_count} rows with income > 50000",
            file=sys.stderr,
            flush=True,
        )

    except Exception as e:
        print(f"Error in filter_high_income_jsonl: {e}", file=sys.stderr)


def main():
    input_data = sys.stdin.read()
    run_streaming(input_data)


if __name__ == "__main__":
    main()
