package engine

import "context"

type RustExecutor struct {
	docker *DockerExecutor
}

func (r *RustExecutor) Execute(ctx context.Context, step Step, input []byte) ([]byte, error) {
	return r.docker.Execute(ctx, step, input)
}
