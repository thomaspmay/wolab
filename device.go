package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func addDevice(cmd *cobra.Command, args []string) {
	uuid, _ := cmd.Flags().GetString("uuid")
	name, _ := cmd.Flags().GetString("name")
	hostname, _ := cmd.Flags().GetString("hostname")
	mac, _ := cmd.Flags().GetString("mac")
	ip, _ := cmd.Flags().GetString("ip")

	d := device{
		UUID:     uuid,
		Name:     name,
		Hostname: hostname,
		MAC:      mac,
		IP:       ip,
	}

	devices, err := loadDevices()
	if err != nil {
		fmt.Println("Error loading devices:", err)
		return
	}

	devices.Devices = append(devices.Devices, d)

	if err := saveDevices(devices); err != nil {
		fmt.Println("Error saving device:", err)
		return
	}

	fmt.Println("Device created successfully.")
}

func listDevices(cmd *cobra.Command, args []string) {
	devices, err := loadDevices()
	if err != nil {
		fmt.Println("Error loading devices:", err)
		return
	}

	for _, d := range devices.Devices {
		fmt.Printf("UUID: %s, Name: %s, Hostname: %s, MAC: %s, IP: %s\n", d.UUID, d.Name, d.Hostname, d.MAC, d.IP)
	}
}

func updateDevice(cmd *cobra.Command, args []string) {
	uuid, _ := cmd.Flags().GetString("uuid")
	name, _ := cmd.Flags().GetString("name")
	hostname, _ := cmd.Flags().GetString("hostname")
	mac, _ := cmd.Flags().GetString("mac")
	ip, _ := cmd.Flags().GetString("ip")

	devices, err := loadDevices()
	if err != nil {
		fmt.Println("Error loading devices:", err)
		return
	}

	for i, d := range devices.Devices {
		if d.UUID == uuid {
			if name != "" {
				devices.Devices[i].Name = name
			}
			if hostname != "" {
				devices.Devices[i].Hostname = hostname
			}
			if mac != "" {
				devices.Devices[i].MAC = mac
			}
			if ip != "" {
				devices.Devices[i].IP = ip
			}
			break
		}
	}

	if err := saveDevices(devices); err != nil {
		fmt.Println("Error saving device:", err)
		return
	}

	fmt.Println("Device updated successfully.")
}

func deleteDevice(cmd *cobra.Command, args []string) {
	uuid, _ := cmd.Flags().GetString("uuid")

	devices, err := loadDevices()
	if err != nil {
		fmt.Println("Error loading devices:", err)
		return
	}

	for i, d := range devices.Devices {
		if d.UUID == uuid {
			devices.Devices = append(devices.Devices[:i], devices.Devices[i+1:]...)
			break
		}
	}

	if err := saveDevices(devices); err != nil {
		fmt.Println("Error saving device:", err)
		return
	}

	fmt.Println("Device deleted successfully.")
}
