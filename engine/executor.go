package engine

import (
	"context"
	"fmt"
	"io"
)

// Executor interface voor verschillende taal executors
type Executor interface {
	Execute(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error)
	Stream(ctx context.Context, step Step, input io.Reader) (io.Reader, error) // new method
}

func NewExecutor(step Step) Executor {
	switch step.Language {
	case "python":
		return &PythonExecutor{docker: NewDockerExecutor()}
	case "rust":
		return &RustExecutor{docker: NewDockerExecutor()}
	default:
		return &DockerExecutor{}
	}
}

// PythonExecutor runt Python scripts in een container
// PythonExecutor runt Python scripts in een container via Docker
type PythonExecutor struct {
	docker *DockerExecutor
}

func (p *PythonExecutor) Execute(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
	// Als cmd leeg is, stel automatisch het Python commando in
	if len(cmd) == 0 {
		cmd = []string{"python", fmt.Sprintf("/app/funcs/%s.py", step.Function)}
	}
	return p.docker.Execute(ctx, step, input, cmd)
}

func (p *PythonExecutor) Stream(ctx context.Context, step Step, input io.Reader) (io.Reader, error) {
	return p.docker.Stream(ctx, step, input)
}

// RustExecutor runt Rust binaries in een container
type RustExecutor struct {
	docker *DockerExecutor
}

func (r *RustExecutor) Execute(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
	if len(cmd) == 0 {
		cmd = []string{fmt.Sprintf("/app/%s", step.Function)}
	}
	return r.docker.Execute(ctx, step, input, cmd)
}

func (r *RustExecutor) Stream(ctx context.Context, step Step, input io.Reader) (io.Reader, error) {
	return r.docker.Stream(ctx, step, input)
}