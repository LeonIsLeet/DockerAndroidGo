package dockermanager

import (
	"fmt"
	"net"
	"net/url"

	"github.com/docker/go-connections/nat"
)

// configurePorts configures the necessary port mappings for the container
func configurePorts(appiumURL string) (nat.PortSet, nat.PortMap, error) {
	appiumPort := extractPort(appiumURL)

	// Expose necessary ports for ADB and emulator
	ports := []string{"5555", "5554", appiumPort}

	portSet := nat.PortSet{}
	portMap := nat.PortMap{}

	for _, port := range ports {
		containerPort := nat.Port(fmt.Sprintf("%s/tcp", port))
		portSet[containerPort] = struct{}{}
		portMap[containerPort] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: port,
			},
		}
	}

	return portSet, portMap, nil
}

// extractPort extracts the port from a URL string
func extractPort(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "4723" // default port
	}
	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "4723"
	}
	return port
}
