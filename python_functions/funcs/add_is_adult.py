def run(data):
    import csv
    import io

    reader = csv.DictReader(io.StringIO(data), delimiter=";")
    rows = list(reader)
    for row in rows:
        row["is_adult"] = str(int(row["age"]) >= 18)
    output = io.StringIO()
    writer = csv.DictWriter(
        output, fieldnames=reader.fieldnames + ["is_adult"], delimiter=";"
    )
    writer.writeheader()
    writer.writerows(rows)
    return output.getvalue()
