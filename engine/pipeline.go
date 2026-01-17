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

	// SOURCE
	logger.Info("loading source from %s", s.Source.Path)
	data, err := os.ReadFile(s.Source.Path)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}
	reg.Set("source", data)

	docker := NewDockerExecutor()

	// PIPELINE
	for i, step := range s.Pipeline {
		logger.Info("step %d start: %s", i, step.ID)

		if step.Type == "sleep" {
			logger.Info("sleeping %d ms", step.DurationMs)
			time.Sleep(time.Millisecond * time.Duration(step.DurationMs))
			logger.Info("step %d finished sleep", i)
			continue
		}

		if step.Executor != "docker" {
			return fmt.Errorf("unsupported executor: %s", step.Executor)
		}

		input := reg.Get(step.Input)
		if input == nil {
			return fmt.Errorf("missing input key: %s", step.Input)
		}

		logger.Info("executing docker image=%s function=%s", step.Image, step.Function)
		logger.Debugf("input key=%s size=%d", step.Input, len(input))

		out, err := docker.Execute(ctx, step.Image, step.Function, input)
		if err != nil {
			return fmt.Errorf("step %s failed: %w", step.ID, err)
		}

		reg.Set(step.Output, out)
		logger.Debugf("output key=%s size=%d", step.Output, len(out))
		logger.Info("step %d finished: %s", i, step.ID)
	}

	// SINK
	result := reg.Get(s.Sink.Input)
	if result == nil {
		return fmt.Errorf("missing sink input: %s", s.Sink.Input)
	}

	// maak directory aan
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
