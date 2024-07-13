package main

import (
	"fmt"
	"reflect"
)

type device struct {
	UUID     string
	Name     string
	Hostname string
	MAC      string
	IP       string
}

// main function which creates a new device and saves it to a file

func main() {

	d := device{
		UUID:     "1234",
		Name:     "Device1",
		Hostname: "device1.local",
		MAC:      "00:11:22:33:44:55",
		IP:       "192.168.0.101",
	}

	filename := "test_device_info.toml"

	if err := saveDeviceInfo(filename, d); err != nil {
		fmt.Printf("Failed to save device info: %v", err)
	}

	loadedDevice, err := loadDeviceInfo(filename)
	if err != nil {
		fmt.Printf("Failed to load device info: %v", err)

	}

	if !reflect.DeepEqual(d, loadedDevice) {
		fmt.Printf("Loaded device info does not match saved device info. Got %+v, expected %+v", loadedDevice, d)
	}
}
