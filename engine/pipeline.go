package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ise/logger"
	"ise/strategy"
)

func Run(ctx context.Context, s strategy.Strategy) error {
	reg := NewRegistry()

	logger.Info("loading source from %s", s.Source.Path)

	// Load source
	data, err := os.ReadFile(s.Source.Path)
	if err != nil {
		return fmt.Errorf("read source failed: %w", err)
	}
	reg.Set("source", data)

	docker := NewDockerExecutor()

	for i, step := range s.Pipeline {
		logger.Info("step %d start", i)

		// Sleep step
		if step.Type == "sleep" {
			logger.Info("sleeping %d ms", step.DurationMs)
			time.Sleep(time.Millisecond * time.Duration(step.DurationMs))
			continue
		}

		// Validate input
		input := reg.Get(step.Input)
		if input == nil {
			return fmt.Errorf("step %d missing input key '%s'", i, step.Input)
		}

		logger.Info("executing docker image %s", step.Image)
		logger.Debugf("input key=%s size=%d", step.Input, len(input))

		start := time.Now()

		out, err := docker.Execute(ctx, step.Image, input)
		if err != nil {
			return fmt.Errorf("step %d failed: %w", i, err)
		}

		duration := time.Since(start)

		logger.Info("step %d finished in %s", i, duration)
		logger.Debugf("output key=%s size=%d", step.Output, len(out))

		reg.Set(step.Output, out)
	}

	// Write sink
	logger.Info("writing result to %s", s.Sink.Path)

	result := reg.Get(s.Sink.Input)
	if result == nil {
		return fmt.Errorf("sink input key '%s' not found", s.Sink.Input)

	} // maak directory aan indien nodig
	dir := filepath.Dir(s.Sink.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create output dir: %w", err)
	}
	logger.Info("writing result to %s", s.Sink.Path)
	if err := os.WriteFile(s.Sink.Path, result, 0644); err != nil {
		return fmt.Errorf("write sink failed: %w", err)
	}

	logger.Info("pipeline completed successfully")

	return nil
}
