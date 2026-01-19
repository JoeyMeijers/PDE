import sys
import json


def transform(records):
    for record in records:
        try:
            age = int(record.get("age", 0))
            if age < 30:
                record["age_group"] = "young"
            elif age < 60:
                record["age_group"] = "adult"
            else:
                record["age_group"] = "senior"
        except:
            record["age_group"] = "unknown"
    return records


def main():
    input_data = sys.stdin.read().strip()
    try:
        records = json.loads(input_data) if input_data else []
    except json.JSONDecodeError:
        records = []

    records = transform(records)

    sys.stdout.write(json.dumps(records))
    sys.stdout.write("\n")
    sys.stdout.flush()


if __name__ == "__main__":
    main()
