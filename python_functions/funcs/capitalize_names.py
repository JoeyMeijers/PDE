import csv
import io


def run(input_data: str) -> str:
    reader = csv.DictReader(io.StringIO(input_data), delimiter=";")
    output = io.StringIO()
    fieldnames = reader.fieldnames or []
    writer = csv.DictWriter(output, fieldnames=fieldnames, delimiter=";")
    writer.writeheader()

    # Process each row
    for row in reader:
        if "name" in row:
            row["name"] = _capitalize_name(row["name"])
        writer.writerow(row)

    return output.getvalue()


def _capitalize_name(name: str) -> str:
    return name.upper()
