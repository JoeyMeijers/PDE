import sys
import importlib


def main():
    if len(sys.argv) < 2:
        print("Function name required", file=sys.stderr)
        sys.exit(1)

    func_name = sys.argv[1]
    input_data = sys.stdin.read()

    module = importlib.import_module(f"funcs.{func_name}")

    result = module.run(input_data)

    if result is None:
        return

    if isinstance(result, (dict, list)):
        import json

        sys.stdout.write(json.dumps(result))
    else:
        sys.stdout.write(str(result))


if __name__ == "__main__":
    main()
