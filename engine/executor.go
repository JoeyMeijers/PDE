package engine

import "context"

type Executor interface {
	Execute(ctx context.Context, step Step, input []byte) ([]byte, error)
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
