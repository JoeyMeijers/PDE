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
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &DockerExecutor{cli: cli}
}

// ExecuteWithCmd runs a container with input []byte and collects output in memory.
func (d *DockerExecutor) ExecuteWithCmd(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
	reader := bytes.NewReader(input)
	outReader, err := d.Stream(ctx, step, reader)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, outReader)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Execute implementeert het Executor interface voor generieke containers
func (d *DockerExecutor) Execute(ctx context.Context, step Step, input []byte, cmd []string) ([]byte, error) {
	if len(cmd) == 0 {
		// fallback, gebruik step.Function als standaard commando
		cmd = []string{fmt.Sprintf("/app/%s", step.Function)}
	}
	return d.ExecuteWithCmd(ctx, step, input, cmd)
}

// Stream runs a container and streams input/output safely.
func (d *DockerExecutor) Stream(ctx context.Context, step Step, input io.Reader) (io.Reader, error) {
	var cmd []string
	if step.Language == "python" {
		cmd = []string{"python", fmt.Sprintf("/app/funcs/%s.py", step.Function)}
	} else {
		cmd = []string{fmt.Sprintf("/app/%s", step.Function)}
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

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	// Stream input
	go func() {
		defer attach.CloseWrite()
		io.Copy(attach.Conn, input)
	}()

	// Use pipe for stdout
	stdoutPipe, stdoutWriter := io.Pipe()
	go func() {
		defer stdoutWriter.Close()
		var stderrBuf bytes.Buffer
		_, err := stdcopy.StdCopy(stdoutWriter, &stderrBuf, attach.Reader)
		if err != nil {
			fmt.Printf("docker stdcopy error: %v\n", err)
		}
		if stderrBuf.Len() > 0 {
			fmt.Printf("docker container stderr: %s\n", stderrBuf.String())
		}
		d.cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
	}()

	return stdoutPipe, nil
}