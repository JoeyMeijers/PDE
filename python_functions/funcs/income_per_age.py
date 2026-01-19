import sys
import json


def transform(records):
    for record in records:
        age = record.get("age")
        income = record.get("income")
        if age and income and age != 0:
            record["income_per_age"] = income / age
        else:
            record["income_per_age"] = None
    return records


def main():
    input_data = sys.stdin.read()
    records = json.loads(input_data)
    transformed_records = transform(records)
    json.dump(transformed_records, sys.stdout)


if __name__ == "__main__":
    main()
