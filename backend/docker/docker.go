package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerAPI interface {
	ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error)
}

type Client struct {
	// client *client.Client
	cli DockerAPI
}

func NewDockerClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	return &Client{
		cli: cli,
	}, nil
}

func (d *Client) ListContainers() ([]types.Container, error) {
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}
