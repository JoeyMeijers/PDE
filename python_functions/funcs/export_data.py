import sys
import os
import json
import csv


def run_streaming() -> None:
    """
    Streams JSONL input line-by-line without buffering the entire dataset.
    Supports CSV and JSON output formats.
    """
    export_path = os.environ.get("EXPORT_PATH", "output.json")

    records = []
    first_record = True
    csv_writer = None
    csv_file = None

    try:
        if export_path.endswith(".csv"):
            # For CSV, we need to know fieldnames from first record
            csv_file = open(export_path, mode="w", newline="", encoding="utf-8")
        elif export_path.endswith(".jsonl"):
            # For JSONL, just stream directly
            output_file = open(export_path, mode="w", encoding="utf-8")
        else:
            # For JSON, open in write mode for array
            output_file = open(export_path, mode="w", encoding="utf-8")
            output_file.write("[\n")

        for line in sys.stdin:
            line = line.strip()
            if not line:
                continue

            try:
                record = json.loads(line)
            except json.JSONDecodeError:
                continue

            if export_path.endswith(".csv"):
                # Initialize CSV writer on first record
                if first_record and csv_file:
                    csv_writer = csv.DictWriter(csv_file, fieldnames=record.keys())
                    csv_writer.writeheader()

                if csv_writer:
                    csv_writer.writerow(record)

            elif export_path.endswith(".jsonl"):
                # Stream JSONL directly
                output_file.write(json.dumps(record, ensure_ascii=False) + "\n")
                output_file.flush()

            else:
                # Stream JSON array
                if not first_record:
                    output_file.write(",\n")
                output_file.write("  " + json.dumps(record, ensure_ascii=False))
                output_file.flush()

            first_record = False

        # Close files properly
        if export_path.endswith(".csv") and csv_file:
            csv_file.close()
        elif export_path.endswith(".jsonl"):
            output_file.close()
        else:
            output_file.write("\n]\n")
            output_file.close()

    except Exception as e:
        if csv_file and not csv_file.closed:
            csv_file.close()
        rapde e


def main():
    run_streaming()


if __name__ == "__main__":
    main()
