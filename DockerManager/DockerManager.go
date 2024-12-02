package dockermanager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// DockerManager manages Docker containers
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

// CreateAndStartContainer creates and starts a Docker container with specific configurations
func (dm *DockerManager) CreateAndStartContainer(cfg DeviceConfig) (string, error) {
	ctx := context.Background()

	imageName := fmt.Sprintf("budtmo2/docker-android-pro:emulator_%s", cfg.CapabilityVersion)
	containerName := cfg.CapabilityDeviceName

	// Pull the image
	imageSummary, err := dm.client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: imageName,
		}),
	})

	if err != nil {
		return "", fmt.Errorf("failed to list images: %w", err)
	}

	// Pull the image only if it doesn't exist
	if len(imageSummary) == 0 {
		reader, err := dm.client.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to pull image: %w", err)
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}
	// io.Copy(os.Stdout, reader)

	// Expose necessary ports
	exposedPorts, portBindings, err := configurePorts(cfg.AppiumURL)
	if err != nil {
		return "", err
	}

	// Container configuration
	config := &container.Config{
		Image: imageName,
		Env: []string{
			fmt.Sprintf("DEVICE=%s", cfg.CapabilityDeviceName),
			fmt.Sprintf("EMULATOR_LANGUAGE=%s", "en"),
			fmt.Sprintf("EMULATOR_COUNTRY=%s", "US"),
		},
		ExposedPorts: exposedPorts,
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Privileged:   true, // Needed for Android emulators
	}

	// Create the container
	resp, err := dm.client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	err = dm.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// Close cleans up resources used by the DockerManager
func (dm *DockerManager) Close() error {
	// Implement cleanup if necessary
	return dm.client.Close()
}
