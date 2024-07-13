package main

import (
	"os"
	"reflect"
	"testing"
)

func TestSaveDeviceInfo(t *testing.T) {
	d := device{
		UUID:     "1234",
		Name:     "Device1",
		Hostname: "device1.local",
		MAC:      "00:11:22:33:44:55",
		IP:       "192.168.1.2",
	}

	filename := "test_device_info.toml"
	defer os.Remove(filename)

	if err := saveDeviceInfo(filename, d); err != nil {
		t.Fatalf("Failed to save device info: %v", err)
	}

	loadedDevice, err := loadDeviceInfo(filename)
	if err != nil {
		t.Fatalf("Failed to load device info: %v", err)
	}

	if !reflect.DeepEqual(d, loadedDevice) {
		t.Fatalf("Loaded device info does not match saved device info. Got %+v, expected %+v", loadedDevice, d)
	}
}

func TestLoadDeviceInfo_FileNotFound(t *testing.T) {
	_, err := loadDeviceInfo("non_existent_file.toml")
	if err == nil {
		t.Fatalf("Expected error when loading non-existent file, got nil")
	}
}
