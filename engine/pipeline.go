package engine

import (
	"context"
	"fmt"
	"os"
	"time"
)

func Run(ctx context.Context, s Strategy) error {
	reg := NewRegistry()
	docker := NewDockerExecutor()

	// Load source
	data, err := os.ReadFile(s.Source.Path)
	if err != nil {
		return err
	}
	reg.Set("source", data)

	for _, step := range s.Pipeline {
		if step.Type == "sleep" {
			time.Sleep(time.Millisecond * time.Duration(step.DurationMs))
			continue
		}

		input := reg.Get(step.Input)
		out, err := docker.Execute(ctx, step, input)
		if err != nil {
			return fmt.Errorf("step %s failed: %w", step.ID, err)
		}

		reg.Set(step.Output, out)
	}

	// Write sink
	result := reg.Get(s.Sink.Input)
	return os.WriteFile(s.Sink.Path, result, 0644)
}
