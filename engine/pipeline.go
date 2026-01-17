package engine

import (
	"context"
	"fmt"
	"os"
	"time"

	"ise/logger"
)

func Run(ctx context.Context, s Strategy) error {
	reg := NewRegistry()
	docker := NewDockerExecutor()

	// Source
	data, err := os.ReadFile(s.Source.Path)
	if err != nil {
		return fmt.Errorf("read source failed: %w", err)
	}
	reg.Set("source", data)
	logger.Info("loading source from %s", s.Source.Path)

	for i, step := range s.Pipeline {
		logger.Info("step %d start: %s", i, step.ID)
		if step.Type == "sleep" {
			logger.Info("sleeping %d ms", step.DurationMs)
			time.Sleep(time.Millisecond * time.Duration(step.DurationMs))
			logger.Info("step %d finished sleep", i)
			continue
		}

		input := reg.Get(step.Input)
		out, err := docker.Execute(ctx, step, input)
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
