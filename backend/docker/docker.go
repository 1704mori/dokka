package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerAPI interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
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
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (d *Client) FindContainer(id string) (types.Container, error) {
	var container types.Container

	containers, err := d.ListContainers()
	if err != nil {
		return container, err
	}

	for _, c := range containers {
		if c.ID == id {
			container = c
			break
		}
	}

	if container.ID == "" {
		return container, errors.New(fmt.Sprintf("unable to find container %s", id))
	}

	return container, nil
}

func (d *Client) StopContainer(id string) (bool, error) {
	err := d.cli.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Client) RestartContainer(id string) (bool, error) {
	err := d.cli.ContainerRestart(context.Background(), id, container.StopOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}
