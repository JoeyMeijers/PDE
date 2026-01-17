package engine

import (
	"context"
	"fmt"
	"os"

	"ise/logger"
)

func Run(ctx context.Context, s Strategy) error {
	reg := NewRegistry()

	// Read raw CSV file content
	if s.Source.Path == "" {
		logger.Error("no source path specified")
		return fmt.Errorf("no source path specified")
	}

	csvContent, err := os.ReadFile(s.Source.Path)
	if err != nil {
		logger.Error("failed to read CSV file: %v", err)
		return err
	}
	if len(csvContent) == 0 {
		logger.Error("CSV file is empty")
		return fmt.Errorf("CSV file is empty")
	}

	// Store raw CSV content for pipeline to process
	reg.Set("data.source", csvContent)

	for i, step := range s.Pipeline {
		logger.Info("step %d start: %s", i, step.ID)

		input := reg.Get(step.Input)

		executor := NewExecutor(step)
		out, err := executor.Execute(ctx, step, input, nil) // cmd=nil, executor bepaalt zelf
		if err != nil {
			logger.Error("step %s failed: %v", step.ID, err)
			return err
		}
		reg.Set(step.Output, out)
		logger.Debugf("output key=%s size=%d", step.Output, len(out))
		logger.Info("step %d finished: %s", i, step.ID)
	}

	result := reg.Get(s.Sink.Input)
	if err := os.WriteFile(s.Sink.Path, result, 0644); err != nil {
		logger.Error("write sink failed: %v", err)
		return err
	}
	logger.Info("writing result to %s", s.Sink.Path)
	return nil
}
