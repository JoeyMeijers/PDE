def run(data):
    import csv
    import io

    reader = csv.DictReader(io.StringIO(data), delimiter=";")
    rows = [
        row
        for row in reader
        if row.get("is_adult") == "True" or row.get("is_adult") == "1"
    ]
    output = io.StringIO()
    writer = csv.DictWriter(output, fieldnames=reader.fieldnames, delimiter=";")
    writer.writeheader()
    writer.writerows(rows)
    return output.getvalue()
