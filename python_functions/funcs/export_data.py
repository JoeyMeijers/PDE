import sys
import os
import json
import csv


def transform(records):
    """
    Optional transformations before export.
    By default, returns the records unchanged.
    """
    return records


def run(input_data: str) -> None:
    data = input_data.strip()
    if not data:
        return

    try:
        records = json.loads(data)
    except json.JSONDecodeError:
        return

    if not isinstance(records, list):
        return

    records = transform(records)

    export_path = os.environ.get("EXPORT_PATH", "output.json")

    if export_path.endswith(".csv"):
        if not records:
            # No data to write
            return
        with open(export_path, mode="w", newline="", encoding="utf-8") as f:
            writer = csv.DictWriter(f, fieldnames=records[0].keys())
            writer.writeheader()
            writer.writerows(records)
    else:
        with open(export_path, "w", encoding="utf-8") as f:
            json.dump(records, f, ensure_ascii=False, indent=2)


def main():
    input_data = sys.stdin.read()
    run(input_data)


if __name__ == "__main__":
    main()
