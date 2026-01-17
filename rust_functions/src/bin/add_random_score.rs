use rand::Rng;
use serde_json::{json, Value};
use std::io::{self, BufRead, Write};

fn main() {
    // Stream JSONL from stdin
    let stdin = io::stdin();
    let reader = stdin.lock();
    let mut rng = rand::thread_rng();

    for line in reader.lines() {
        if let Ok(line) = line {
            if line.trim().is_empty() {
                continue;
            }

            match serde_json::from_str::<Value>(&line) {
                Ok(mut entry) => {
                    if entry.get("score").is_some() {
                        let random_score: f64 = rng.gen_range(0.0..100.0);
                        entry["random_score"] = json!(random_score);
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