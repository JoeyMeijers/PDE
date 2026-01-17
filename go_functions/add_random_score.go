package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		// Only add random_score if score field exists
		if _, hasScore := entry["score"]; hasScore {
			entry["random_score"] = rand.Float64() * 100
		}

		output, err := json.Marshal(entry)
		if err != nil {
			continue
		}

		writer.Write(output)
		writer.WriteString("\n")
		rowCount++

		// Flush every 50000 rows
		if rowCount%50000 == 0 {
			writer.Flush()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	writer.Flush()
}
