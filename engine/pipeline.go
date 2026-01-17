package engine

import (
	"context"
	"fmt"
	"io"
	"os"

	"ise/logger"
)

func Run(ctx context.Context, s Strategy) error {
	reg := NewRegistry()

	// Check source path
	if s.Source.Path == "" {
		logger.Error("no source path specified")
		return fmt.Errorf("no source path specified")
	}

	// Open CSV file for streaming
	srcFile, err := os.Open(s.Source.Path)
	if err != nil {
		logger.Error("failed to open CSV file: %v", err)
		return err
	}
	defer srcFile.Close()

	// Store source reader in registry
	reg.SetReader("data.source", srcFile)

	for i, step := range s.Pipeline {
		logger.Info("step %d start: %s", i, step.ID)

		inputReader := reg.GetReader(step.Input)

		executor := NewExecutor(step)
		outReader, err := executor.Stream(ctx, step, inputReader) // Stream returns io.Reader
		if err != nil {
			logger.Error("step %s failed: %v", step.ID, err)
			return err
		}

		reg.SetReader(step.Output, outReader)
		logger.Info("step %d finished: %s", i, step.ID)
	}

	// Write final sink directly to file
	outReader := reg.GetReader(s.Sink.Input)
	sinkFile, err := os.Create(s.Sink.Path)
	if err != nil {
		logger.Error("failed to create sink file: %v", err)
		return err
	}
	defer sinkFile.Close()

	_, err = io.Copy(sinkFile, outReader)
	if err != nil {
		logger.Error("failed to write sink: %v", err)
		return err
	}

	logger.Info("writing result to %s", s.Sink.Path)
	return nil
}
