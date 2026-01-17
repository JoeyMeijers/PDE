import sys
import csv
import io

# lees CSV van stdin
input_csv = sys.stdin.read()
reader = csv.DictReader(io.StringIO(input_csv), delimiter=";")
rows = list(reader)

# voeg kolom is_adult toe
for row in rows:
    row["is_adult"] = str(int(row["age"]) >= 18)

# schrijf CSV naar stdout
output = io.StringIO()
fieldnames = reader.fieldnames + ["is_adult"]
writer = csv.DictWriter(output, fieldnames=fieldnames, delimiter=";")
writer.writeheader()
writer.writerows(rows)

print(output.getvalue(), end="")
