import sys
import json


def run(input_data: str) -> str:
    data = input_data.strip()
    if not data:
        return "[]"
    try:
        people = json.loads(data)
    except json.JSONDecodeError:
        return "[]"
    if not isinstance(people, list):
        return "[]"
    adults = [p for p in people if int(p.get("age", 0)) >= 18]
    return json.dumps(adults)


def main():
    input_data = sys.stdin.read()
    output = run(input_data)
    sys.stdout.write(output)


if __name__ == "__main__":
    main()
