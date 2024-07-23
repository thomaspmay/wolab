package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestAddDevice(t *testing.T) {

	defer os.Remove(filename)

	cmd := &cobra.Command{}
	cmd.Flags().String("uuid", "1234", "")
	cmd.Flags().String("name", "Device1", "")
	cmd.Flags().String("hostname", "device1.local", "")
	cmd.Flags().String("mac", "00:11:22:33:44:55", "")
	cmd.Flags().String("ip", "192.168.1.2", "")

	addDevice(cmd, []string{})

	devices, err := loadDevices()
	if err != nil {
		t.Fatalf("Failed to load devices: %v", err)
	}

	expected := device{
		UUID:     "1234",
		Name:     "Device1",
		Hostname: "device1.local",
		MAC:      "00:11:22:33:44:55",
		IP:       "192.168.1.2",
	}

	if len(devices.Devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(devices.Devices))
	}

	if !reflect.DeepEqual(devices.Devices[0], expected) {
		t.Fatalf("Expected %+v, got %+v", expected, devices.Devices[0])
	}
}

func TestAddSecondDevice(t *testing.T) {
	defer os.Remove(filename)

	// Add the first device
	cmd1 := &cobra.Command{}
	cmd1.Flags().String("uuid", "1234", "")
	cmd1.Flags().String("name", "Device1", "")
	cmd1.Flags().String("hostname", "device1.local", "")
	cmd1.Flags().String("mac", "00:11:22:33:44:55", "")
	cmd1.Flags().String("ip", "192.168.1.2", "")

	addDevice(cmd1, []string{})

	// Add the second device
	cmd2 := &cobra.Command{}
	cmd2.Flags().String("uuid", "5678", "")
	cmd2.Flags().String("name", "Device2", "")
	cmd2.Flags().String("hostname", "device2.local", "")
	cmd2.Flags().String("mac", "66:77:88:99:AA:BB", "")
	cmd2.Flags().String("ip", "192.168.1.3", "")

	addDevice(cmd2, []string{})

	devices, err := loadDevices()
	if err != nil {
		t.Fatalf("Failed to load devices: %v", err)
	}

	expectedDevices := []device{
		{
			UUID:     "1234",
			Name:     "Device1",
			Hostname: "device1.local",
			MAC:      "00:11:22:33:44:55",
			IP:       "192.168.1.2",
		},
		{
			UUID:     "5678",
			Name:     "Device2",
			Hostname: "device2.local",
			MAC:      "66:77:88:99:AA:BB",
			IP:       "192.168.1.3",
		},
	}

	if len(devices.Devices) != 2 {
		t.Fatalf("Expected 2 devices, got %d", len(devices.Devices))
	}

	if !reflect.DeepEqual(devices.Devices, expectedDevices) {
		t.Fatalf("Expected %+v, got %+v", expectedDevices, devices.Devices)
	}
}

func TestListDevices(t *testing.T) {

	defer os.Remove(filename)

	devices := deviceList{
		Devices: []device{
			{UUID: "1234", Name: "Device1", Hostname: "device1.local", MAC: "00:11:22:33:44:55", IP: "192.168.1.2"},
		},
	}
	saveDevices(devices)

	cmd := &cobra.Command{}
	var output bytes.Buffer
	cmd.SetOut(&output)

	listDevices(cmd, []string{})

	expected := "UUID: 1234, Name: Device1, Hostname: device1.local, MAC: 00:11:22:33:44:55, IP: 192.168.1.2\n"

	if output.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, output.String())
	}
}

func TestUpdateDevice(t *testing.T) {

	defer os.Remove(filename)

	devices := deviceList{
		Devices: []device{
			{UUID: "1234", Name: "Device1", Hostname: "device1.local", MAC: "00:11:22:33:44:55", IP: "192.168.1.2"},
		},
	}
	saveDevices(devices)

	cmd := &cobra.Command{}
	cmd.Flags().String("uuid", "1234", "")
	cmd.Flags().String("name", "UpdatedDevice", "")
	cmd.Flags().String("hostname", "updated.local", "")
	cmd.Flags().String("mac", "66:77:88:99:AA:BB", "")
	cmd.Flags().String("ip", "10.0.0.1", "")

	updateDevice(cmd, []string{})

	updatedDevices, err := loadDevices()
	if err != nil {
		t.Fatalf("Failed to load devices: %v", err)
	}

	expected := device{
		UUID:     "1234",
		Name:     "UpdatedDevice",
		Hostname: "updated.local",
		MAC:      "66:77:88:99:AA:BB",
		IP:       "10.0.0.1",
	}

	if len(updatedDevices.Devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(updatedDevices.Devices))
	}

	if !reflect.DeepEqual(updatedDevices.Devices[0], expected) {
		t.Fatalf("Expected %+v, got %+v", expected, updatedDevices.Devices[0])
	}
}

func TestDeleteDevice(t *testing.T) {

	defer os.Remove(filename)

	devices := deviceList{
		Devices: []device{
			{UUID: "1234", Name: "Device1", Hostname: "device1.local", MAC: "00:11:22:33:44:55", IP: "192.168.1.2"},
		},
	}
	saveDevices(devices)

	cmd := &cobra.Command{}
	cmd.Flags().String("uuid", "1234", "")

	deleteDevice(cmd, []string{})

	updatedDevices, err := loadDevices()
	if err != nil {
		t.Fatalf("Failed to load devices: %v", err)
	}

	if len(updatedDevices.Devices) != 0 {
		t.Fatalf("Expected 0 devices, got %d", len(updatedDevices.Devices))
	}
}
