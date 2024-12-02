package dockermanager

import (
	"encoding/json"
	"io/ioutil"
)

// DeviceConfig represents the configuration for a device/emulator
type DeviceConfig struct {
	AppiumURL            string `json:"appium_url"`
	CapabilityPlatform   string `json:"capability_platform"`
	CapabilityVersion    string `json:"capability_platform_version"`
	CapabilityDeviceName string `json:"capability_device_name"`
	CapabilityAutomation string `json:"capability_automation_name"`
}

// LoadDeviceConfigs loads device configurations from a JSON file
func LoadDeviceConfigs(filePath string) ([]DeviceConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var configs []DeviceConfig
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
