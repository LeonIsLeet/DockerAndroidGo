package DockerManager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerManager struct {
	client *client.Client
}

// NewDockerManager initializes a new DockerManager instance
func NewDockerManager() (*DockerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerManager{client: cli}, nil
}

// CreateContainer creates a new Docker container
func (dm *DockerManager) CreateContainer(image_name string, containerName string, cmd []string) (string, error) {
	ctx := context.Background()

	// Pull the image if it doesn't exist locally
	_, err := dm.client.ImagePull(ctx, image_name, image.PullOptions{})
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	config := &container.Config{
		Image: image_name,
		Cmd:   cmd,
	}

	resp, err := dm.client.ContainerCreate(ctx, config, nil, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
	return resp.ID, nil
}

// StartContainer starts a Docker container
func (dm *DockerManager) StartContainer(containerID string) error {
	ctx := context.Background()
	err := dm.client.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

// StopContainer stops a Docker container
func (dm *DockerManager) StopContainer(containerID string, timeout *time.Duration) error {
	ctx := context.Background()
	var timeoutSeconds *int
	if timeout != nil {
		seconds := int(timeout.Seconds())
		timeoutSeconds = &seconds
	}

	err := dm.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: timeoutSeconds})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

// GetContainerStats retrieves the stats of a running container
func (dm *DockerManager) GetContainerStats(containerID string) (types.StatsJSON, error) {
	ctx := context.Background()
	resp, err := dm.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return types.StatsJSON{}, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer resp.Body.Close()

	var stats types.StatsJSON
	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		return types.StatsJSON{}, fmt.Errorf("failed to decode stats: %w", err)
	}
	return stats, nil
}

// Close cleans up resources used by the DockerManager
func (dm *DockerManager) Close() error {
	return dm.client.Close()
}
