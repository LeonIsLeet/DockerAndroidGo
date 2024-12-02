package main

import (
	"log"
	"net"
	"net/url"
	"os/exec"
	"time"

	"androidauto/dockermanager"

	"github.com/electricbubble/guia2"
)

func main() {
	// Load device configurations
	configs, err := dockermanager.LoadDeviceConfigs("config.json")
	if err != nil {
		log.Fatalf("Failed to load device configurations: %v", err)
	}

	// Initialize DockerManager
	dm, err := dockermanager.NewDockerManager()
	if err != nil {
		log.Fatalf("Failed to create DockerManager: %v", err)
	}
	defer dm.Close()

	// Iterate over device configurations
	for _, cfg := range configs {
		// Create and start the Docker container
		containerID, err := dm.CreateAndStartContainer(cfg)
		if err != nil {
			log.Printf("Failed to create and start container for device %s: %v", cfg.CapabilityDeviceName, err)
			continue
		}
		log.Printf("Started container %s for device %s", containerID, cfg.CapabilityDeviceName)

		// Start Appium server on specified port
		appiumPort := extractPort(cfg.AppiumURL)
		err = startAppiumServer(appiumPort)
		if err != nil {
			log.Printf("Failed to start Appium server on port %s: %v", appiumPort, err)
			continue
		}
		log.Printf("Started Appium server on port %s", appiumPort)

		// Wait for emulator to boot up
		time.Sleep(60 * time.Second) // Adjust as needed

		// Connect to the device using guia2
		driver, err := guia2.NewWiFiDriver("localhost")
		if err != nil {
			log.Printf("Failed to create guia2 driver for device %s: %v", cfg.CapabilityDeviceName, err)
			continue
		}
		defer driver.Dispose()

		// Use the driver as needed
		// Example: Launch an app
		err = driver.AppLaunch("com.example.app")
		if err != nil {
			log.Printf("Failed to launch app on device %s: %v", cfg.CapabilityDeviceName, err)
		}
	}
}

// startAppiumServer starts the Appium server on the specified port
func startAppiumServer(port string) error {
	cmd := exec.Command("appium", "-p", port)
	return cmd.Start()
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
