package engine

import (
	"context"
	"os"
	"time"

	"ise/strategy"
)

func Run(ctx context.Context, s strategy.Strategy) error {
	reg := NewRegistry()

	// source
	data, _ := os.ReadFile(s.Source.Path)
	reg.Set("source", data)

	docker := NewDockerExecutor()

	for _, step := range s.Pipeline {
		if step.Type == "sleep" {
			time.Sleep(time.Millisecond * time.Duration(step.DurationMs))
			continue
		}

		input := reg.Get(step.Input)
		out, err := docker.Execute(ctx, step.Image, input)
		if err != nil {
			return err
		}
		reg.Set(step.Output, out)
	}

	result := reg.Get(s.Sink.Input)
	return os.WriteFile(s.Sink.Path, result, 0644)
}
