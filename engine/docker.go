package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() *DockerExecutor {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return &DockerExecutor{cli: cli}
}

func (d *DockerExecutor) Execute(ctx context.Context, step Step, input []byte) ([]byte, error) {
	var cmd []string

	switch step.Language {
	case "python":
		cmd = []string{"python", "main.py", step.Function}
	case "rust":
		cmd = []string{"./add_age_plus_one"} // executable in Dockerfile
	case "node":
		cmd = []string{"node", "funcs/capitalize_names.js"}
	default:
		cmd = []string{} // fallback
	}

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:     step.Image,
		Cmd:       cmd,
		OpenStdin: true,
		StdinOnce: true,
		Tty:       false,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	attach, err := d.cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		defer attach.CloseWrite()
		io.Copy(attach.Conn, bytes.NewReader(input))
	}()

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, attach.Reader)

	if stderr.Len() > 0 {
		return nil, fmt.Errorf("container error: %s", stderr.String())
	}

	return stdout.Bytes(), nil
}
