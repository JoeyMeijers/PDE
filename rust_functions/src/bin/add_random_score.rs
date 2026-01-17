use rand::Rng;
use serde_json::{json, Value};
use std::io::{self, BufRead, BufWriter, Write};

fn main() {
    // Stream JSONL from stdin
    let stdin = io::stdin();
    let reader = stdin.lock();
    let mut rng = rand::thread_rng();
    let stdout = io::stdout();
    let mut writer = BufWriter::new(stdout.lock());

    let mut row_count = 0;
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
                        let _ = writer.write_all(output.as_bytes());
                        let _ = writer.write_all(b"\n");
                        row_count += 1;

                        // Flush every 50000 rows for better performance
                        if row_count % 50000 == 0 {
                            let _ = writer.flush();
                        }
                    }
                }
                Err(_) => continue,
            }
        }
    }

    // Final flush
    let _ = writer.flush();
}