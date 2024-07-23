package main

import (
	"github.com/BurntSushi/toml"
	"os"
)

func saveDevices(devices deviceList) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(devices); err != nil {
		return err
	}

	return nil
}

func loadDevices() (deviceList, error) {
	var devices deviceList
	if _, err := toml.DecodeFile(filename, &devices); err != nil {
		if os.IsNotExist(err) {
			return devices, nil
		}
		return devices, err
	}
	return devices, nil
}
