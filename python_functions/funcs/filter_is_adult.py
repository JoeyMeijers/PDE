import sys
import json
import pandas as pd
# import polars as pl  # optionally use Polars


def transform(df: pd.DataFrame) -> pd.DataFrame:
    """
    Pas hier je echte data-transformaties toe:
    - filteren
    - nieuwe kolommen
    - aggregaties
    """
    # Voorbeeld: alleen volwassenen >=18
    if "age" in df.columns:
        df = df[df["age"] >= 18]
        df["is_adult"] = True
    return df


def run(input_data: str) -> str:
    data = input_data.strip()
    if not data:
        return "[]"
    try:
        records = json.loads(data)
    except json.JSONDecodeError:
        return "[]"
    if not isinstance(records, list):
        return "[]"

    # Converteer naar DataFrame
    df = pd.DataFrame(records)
    # of voor Polars:
    # df = pl.DataFrame(records)

    # Transformeer data
    df = transform(df)

    # Converteer terug naar JSON
    return df.to_json(orient="records")


def main():
    input_data = sys.stdin.read()
    output = run(input_data)
    sys.stdout.write(output)


if __name__ == "__main__":
    main()
