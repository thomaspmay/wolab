package main

import (
	"github.com/BurntSushi/toml"
	"os"
)

func saveDeviceInfo(filename string, d device) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(d); err != nil {
		return err
	}

	return nil
}

func loadDeviceInfo(filename string) (device, error) {
	var d device
	if _, err := toml.DecodeFile(filename, &d); err != nil {
		return d, err
	}
	return d, nil
}
