package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"ise/logger"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() *DockerExecutor {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	return &DockerExecutor{cli: cli}
}

func (d *DockerExecutor) Execute(ctx context.Context, image, function string, input []byte) ([]byte, error) {
	logger.Info("docker execute image=%s function=%s", image, function)

	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:     image,
		Cmd:       []string{"python", "main.py", function},
		OpenStdin: true,
		StdinOnce: true,
		Tty:       false,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("container create: %w", err)
	}

	defer func() {
		_ = d.cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
	}()

	attach, err := d.cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("container attach: %w", err)
	}
	defer attach.Close()

	go func() {
		defer attach.CloseWrite()
		_, _ = io.Copy(attach.Conn, bytes.NewReader(input))
	}()

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("container start: %w", err)
	}

	statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	var exitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("container wait: %w", err)
		}
	case status := <-statusCh:
		exitCode = status.StatusCode
	}

	var stdout, stderr bytes.Buffer
	_, _ = stdcopy.StdCopy(&stdout, &stderr, attach.Reader)

	logger.Debugf("docker exit code=%d", exitCode)
	if stderr.Len() > 0 {
		logger.Debugf("docker stderr: %s", stderr.String())
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("container failed (exit=%d): %s", exitCode, stderr.String())
	}

	return stdout.Bytes(), nil
}
