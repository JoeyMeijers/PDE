package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor() *DockerExecutor {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &DockerExecutor{cli: cli}
}

// ExecuteWithCmd: voert een container uit met een specifiek commando
func (d *DockerExecutor) ExecuteWithCmd(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
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
	defer d.cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})

	attach, err := d.cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, err
	}
	defer attach.Close()

	// Start container first
	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	// Create a channel to signal when input is written
	done := make(chan error, 1)

	// Write input in goroutine and close stdin when done
	go func() {
		defer attach.CloseWrite()
		_, err := io.Copy(attach.Conn, bytes.NewReader(input))
		done <- err
	}()

	// Wait for input to be written
	if err := <-done; err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	statusCh, errCh := d.cli.ContainerWait(ctx2, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			d.cli.ContainerKill(context.Background(), resp.ID, "KILL")
			return nil, err
		}
	case <-statusCh:
	case <-ctx2.Done():
		d.cli.ContainerKill(context.Background(), resp.ID, "KILL")
		return nil, fmt.Errorf("container execution timeout")
	}

	var stdout, stderr bytes.Buffer
	stdcopy.StdCopy(&stdout, &stderr, attach.Reader)

	return stdout.Bytes(), nil
}

// Execute implementeert het Executor interface voor generieke containers
func (d *DockerExecutor) Execute(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
	if len(cmd) == 0 {
		// fallback, gebruik step.Function als standaard commando
		cmd = []string{fmt.Sprintf("/app/%s", step.Function)}
	}
	return d.ExecuteWithCmd(ctx, step, input, cmd)
}