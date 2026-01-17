package engine

import "context"

type PythonExecutor struct {
	docker *DockerExecutor
}

func (p *PythonExecutor) Execute(ctx context.Context, step Step, input []byte) ([]byte, error) {
	// DockerExecutor bepaalt CMD op basis van step.Function of step.Language
	return p.docker.Execute(ctx, step, input)
}
