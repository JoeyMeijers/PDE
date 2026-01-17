import sys
import importlib
import csv
import io


def main():
    if len(sys.argv) < 2:
        print("Usage: python main.py <function_name>", file=sys.stderr)
        sys.exit(1)

    func_name = sys.argv[1]

    # Dynamisch module importeren uit funcs/
    try:
        module = importlib.import_module(f"funcs.{func_name}")
    except ModuleNotFoundError:
        print(f"Function {func_name} not found in funcs/", file=sys.stderr)
        sys.exit(1)

    # Functie 'run' uitvoeren, input via stdin, output naar stdout
    if not hasattr(module, "run"):
        print(f"Function {func_name} does not have a 'run' function", file=sys.stderr)
        sys.exit(1)

    input_data = sys.stdin.read()
    output_data = module.run(input_data)
    print(output_data, end="")


if __name__ == "__main__":
    main()
