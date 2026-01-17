import sys
import importlib


def main():
    if len(sys.argv) < 2:
        print("Usage: main.py <function_name>", file=sys.stderr)
        sys.exit(1)

    func_name = sys.argv[1]

    # stdin lezen
    data = sys.stdin.read()

    # functie dynamisch importeren
    module = importlib.import_module(f"funcs.{func_name}")
    result = module.run(data)

    # stdout
    print(result, end="")


if __name__ == "__main__":
    main()
