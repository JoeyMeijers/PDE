use std::io::{self, BufRead};

fn main() {
    let stdin = io::stdin();
    let mut lines = stdin.lock().lines();
    let mut output = Vec::new();

    if let Some(Ok(header)) = lines.next() {
        output.push(format!("{};age_plus_one", header));
    }

    for line in lines {
        if let Ok(line) = line {
            let mut parts: Vec<String> = line.split(';').map(|s| s.to_string()).collect();
            let age: i32 = parts[1].parse().unwrap_or(0);
            let age_plus_one = (age + 1).to_string(); // tijdelijke String in variabele
            parts.push(age_plus_one); // nu veilig
            output.push(parts.join(";"));
        }
    }

    for line in output {
        println!("{}", line);
    }
}
