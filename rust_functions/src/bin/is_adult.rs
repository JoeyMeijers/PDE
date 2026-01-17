use serde_json::{json, Value};
use std::io::{self, BufRead, Write};

fn main() {
    // Stream JSONL from stdin
    let stdin = io::stdin();
    let reader = stdin.lock();

    for line in reader.lines() {
        if let Ok(line) = line {
            if line.trim().is_empty() {
                continue;
            }

            match serde_json::from_str::<Value>(&line) {
                Ok(mut entry) => {
                    if let Some(age_str) = entry.get("age").and_then(|v| v.as_str()) {
                        if let Ok(age) = age_str.parse::<i32>() {
                            entry["is_adult"] = json!(age >= 18);
                        }
                    }

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
