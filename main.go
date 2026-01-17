package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"ise/engine"
	"ise/logger"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	logger.InitLogger(*debug)

	raw, err := os.ReadFile("strategy.json")
	if err != nil {
		logger.Error("read strategy: %v", err)
		os.Exit(1)
	}

	var s engine.Strategy
	if err := json.Unmarshal(raw, &s); err != nil {
		logger.Error("parse strategy: %v", err)
		os.Exit(1)
	}

	if err := engine.Run(context.Background(), s); err != nil {
		logger.Error("run failed: %v", err)
		os.Exit(1)
	}

	logger.Info("strategy finished successfully")
}
