# ISE — Polyglot Data Pipeline Engine

A language-agnostic, container-based data pipeline engine that executes structured data transformations via a flexible pipeline of functions, regardless of whether they're written in Python, Rust, or other languages.

## Key Features

- **Pipeline-based**: Define a sequence of steps that transform, filter, or enrich data
- **Polyglot & Containerized**: Each step runs in a Docker container, fully isolated, in any language (Python, Rust, etc.)
- **Streaming Architecture**: Data flows line-by-line (JSONL format) between steps without buffering the entire dataset in memory
- **Extensible**: Easily add new functions in any language without modifying the core engine
- **Logging & Observability**: Every step is logged with start time, duration, output size, and errors

## Use Cases

- ETL processes for large datasets
- Multi-language function orchestration
- Data transformation prototyping and statistical calculations
- Stress and performance testing of pipelines

---

## How It Works

### 1. Architecture Overview

ISE consists of three main components:

1. **Go Engine** (`engine/`) — Orchestrates the pipeline execution, manages I/O streaming
2. **Strategy Definition** (`strategy.json`) — Declares the pipeline steps, their inputs/outputs, and Docker images
3. **Functions** (`python_functions/`, `rust_functions/`) — Containerized transformation logic in various languages

### 2. Data Flow

Data flows through the pipeline in **JSONL format** (JSON Lines: one JSON object per line), streaming from step to step:

```
CSV File (sources/big_test.csv)
    ↓
[load_data - Python] → converts CSV to JSONL
    ↓ (data.loaded)
[filter_high_income - Python] → filters records by income > 50000
    ↓ (data.filtered)
[is_adult - Rust] → adds "is_adult" boolean field
    ↓ (data.final)
[export_data - Python] → exports to output format (JSON/CSV/JSONL)
    ↓
Output File (output/result.json)
```

### 3. Registry: In-Memory Stream Storage

The **Registry** ([engine/registry.go](engine/registry.go)) stores `io.Reader` streams with key-value pairs:

- `"data.source"` — The source CSV file
- `"data.loaded"` — CSV converted to JSONL
- `"data.filtered"` — Filtered JSONL
- `"data.final"` — Final transformed data
- etc.

Each step reads from its **input** key and writes to its **output** key. No temporary files are created — everything streams in memory.

### 4. Container Execution

Each pipeline step:
1. Reads input from stdin (JSONL)
2. Processes it line-by-line or in batches
3. Writes output to stdout (JSONL)
4. Gets captured and stored in the Registry

This architecture enables processing of datasets **larger than available RAM** because data never accumulates in memory — it streams through.

---

## Strategy Configuration

Define your pipeline in `strategy.json`:

```json
{
  "source": {
    "path": "sources/big_test.csv"
  },
  "pipeline": [
    {
      "id": "load_data",
      "language": "python",
      "image": "strategy-python:latest",
      "function": "load_data",
      "input": "data.source",
      "output": "data.loaded"
    },
    {
      "id": "filter_high_income",
      "language": "python",
      "image": "strategy-python:latest",
      "function": "filter_high_income_jsonl",
      "input": "data.loaded",
      "output": "data.filtered"
    },
    {
      "id": "is_adult",
      "language": "rust",
      "image": "strategy-rust:latest",
      "function": "is_adult",
      "input": "data.filtered",
      "output": "data.final"
    },
    {
      "id": "export_data",
      "language": "python",
      "image": "strategy-python:latest",
      "function": "export_data",
      "input": "data.final",
      "output": "data.exported"
    }
  ],
  "sink": {
    "path": "output/result.json",
    "input": "data.final"
  }
}
```

### Configuration Fields

- **`id`** — Unique identifier for the step
- **`language`** — `"python"`, `"rust"`, or other supported language
- **`image`** — Docker image name to run the function in
- **`function`** — Name of the function to execute (maps to filename)
- **`input`** — Registry key for input data
- **`output`** — Registry key where output is stored
- **`source.path`** — Input file path (relative to workspace root)
- **`sink.path`** — Output file path; `sink.input` specifies which registry key to write

---

## How to Add Functions

### Adding a Python Function

1. **Create the function file** in `python_functions/funcs/`:

```python
# python_functions/funcs/my_transform.py

import sys
import json

def main():
    """
    Streams JSONL from stdin, processes each line, outputs JSONL to stdout.
    """
    for line in sys.stdin:
        line = line.strip()
        if not line:
            continue

        try:
            record = json.loads(line)

            # Your transformation logic here
            record["my_field"] = "my_value"

            # Output immediately
            sys.stdout.write(json.dumps(record, ensure_ascii=False) + "\n")
            sys.stdout.flush()

        except json.JSONDecodeError:
            continue

if __name__ == "__main__":
    main()
```

2. **Add to `strategy.json`**:

```json
{
  "id": "my_transform",
  "language": "python",
  "image": "strategy-python:latest",
  "function": "my_transform",
  "input": "data.loaded",
  "output": "data.transformed"
}
```

3. **Build the Docker image** (if not already built):

```bash
docker build -f dockerfiles/python/Dockerfile -t strategy-python:latest .
```

### Adding a Rust Function

1. **Create the binary** in `rust_functions/src/bin/`:

```rust
// rust_functions/src/bin/my_transform.rs

use serde_json::{json, Value};
use std::io::{self, BufRead, Write};

fn main() {
    let stdin = io::stdin();
    let reader = stdin.lock();

    for line in reader.lines() {
        if let Ok(line) = line {
            if line.trim().is_empty() {
                continue;
            }

            match serde_json::from_str::<Value>(&line) {
                Ok(mut entry) => {
                    // Your transformation logic here
                    entry["my_field"] = json!("my_value");

                    if let Ok(output) = serde_json::to_string(&entry) {
                        let _ = io::stdout().write_all(output.as_bytes());
                        let _ = io::stdout().write_all(b"\n");
                        let _ = io::stdout().flush();
                    }
                }
                Err(_) => continue,
            }
        }
    }
}
```

2. **Build the Rust project**:

```bash
cd rust_functions && cargo build --release --bin my_transform
```

3. **Add to Docker image** — Update `dockerfiles/rust/Dockerfile` to include your binary, then rebuild:

```bash
docker build -f dockerfiles/rust/Dockerfile -t strategy-rust:latest .
```

4. **Add to `strategy.json`**:

```json
{
  "id": "my_transform",
  "language": "rust",
  "image": "strategy-rust:latest",
  "function": "my_transform",
  "input": "data.loaded",
  "output": "data.transformed"
}
```

### Function Best Practices

- **Always stream**: Read line-by-line from stdin, don't buffer the entire input
- **Flush output regularly**: Use `sys.stdout.flush()` in Python or `io::stdout().flush()` in Rust
- **Handle errors gracefully**: Skip malformed records with `continue` instead of crashing
- **Use JSONL format**: One JSON object per line for pipeline compatibility
- **Keep containers small**: Only include necessary dependencies in Docker images

---

## Running the Pipeline

### Build and Run

```bash
# Build Docker images
make build

# Run the pipeline
make run

# View logs (with debug output)
go run main.go -debug
```

### Output

Results are written to the path specified in `strategy.json`'s `sink.path`. By default: `output/result.json`

---

## Performance Considerations

### Memory Usage

ISE is designed for **streaming large datasets**:
- Data flows line-by-line through pipes, not in memory
- Each step processes incrementally
- No entire dataset is ever loaded into RAM

**Performance tips:**
- Keep batch sizes reasonable in Python (10K-100K lines)
- Use `.flush()` after each batch to prevent buffering
- Minimize Docker overhead with lightweight base images

### Current Pipeline Performance

The existing pipeline is **fully streaming-capable**:
- `load_data` — Batches CSV rows (10K at a time)
- `filter_high_income` — Line-by-line filtering
- `is_adult` — Line-by-line Rust processing
- `export_data` — Streams output line-by-line (no memory buffering)

---

## Project Structure

```
.
├── main.go                          # Entry point
├── strategy.json                    # Pipeline configuration
├── go.mod                           # Go dependencies
├── Makefile                         # Build commands
│
├── engine/                          # Core pipeline engine
│   ├── pipeline.go                  # Main execution logic
│   ├── registry.go                  # Stream storage
│   ├── executor.go                  # Language executors
│   ├── docker.go                    # Docker container management
│   └── types.go                     # Data structures
│
├── logger/                          # Logging utilities
│   └── logger.go
│
├── python_functions/                # Python functions
│   ├── requirements.txt
│   ├── funcs/
│   │   ├── load_data.py             # CSV → JSONL
│   │   ├── filter_high_income_jsonl.py
│   │   ├── export_data.py           # JSONL → JSON/CSV/JSONL
│   │   └── ... (other functions)
│   └── main.py
│
├── rust_functions/                  # Rust functions
│   ├── Cargo.toml
│   ├── src/bin/
│   │   ├── is_adult.rs              # Adds is_adult field
│   │   └── ... (other binaries)
│   └── target/                      # Build artifacts
│
├── dockerfiles/                     # Docker images
│   ├── python/Dockerfile
│   └── rust/Dockerfile
│
├── sources/                         # Input data
│   └── big_test.csv
│
└── output/                          # Output data
    └── result.json
```

---

## Troubleshooting

### Pipeline hangs or runs out of memory

- Check if functions are streaming (not buffering entire input)
- Verify `.flush()` calls after writing output
- Monitor Docker memory usage: `docker stats`

### Docker image not found

- Ensure images are built: `make build`
- Check image names in `strategy.json` match built images

### Wrong output format

- Verify `export_data.py` output format matches `sink.path` extension
- Check JSONL format: one JSON object per line, no extra characters
