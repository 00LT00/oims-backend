package docker

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Run(timestamp string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	ContainerID := ""
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/oims" {
				ContainerID = container.ID
			}
		}
	}

	resp, err := cli.ContainerExecCreate(ctx, ContainerID, types.ExecConfig{
		WorkingDir:   "/oims/Efficientdet",
		Cmd:          strings.Split("python model_inspect.py --output_image_dir results/ --input_image historys/"+timestamp+"/*.jpg", " "),
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		panic(err)
	}

	resp2, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		panic(err)
	}
	defer resp2.Close()

	err = cli.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		panic(err)
	}

	_, _ = io.Copy(os.Stdout, resp2.Reader)
}
