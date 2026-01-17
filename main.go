package main

import (
	"context"
	"encoding/json"
	"os"

	"ise/engine"
	"ise/strategy"
)

func main() {
	raw, _ := os.ReadFile("strategy.json")
	var s strategy.Strategy
	json.Unmarshal(raw, &s)

	engine.Run(context.Background(), s)
}
