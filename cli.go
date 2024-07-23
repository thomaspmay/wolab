package main

import "github.com/spf13/cobra"

func runCli() {
	var rootCmd = &cobra.Command{Use: "device_manager"}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new device",
		Run:   addDevice,
	}
	createCmd.Flags().String("uuid", "", "Device UUID")

	createCmd.Flags().String("name", "", "Device name")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("hostname", "", "Device hostname")
	createCmd.Flags().String("mac", "", "Device MAC address")
	createCmd.MarkFlagRequired("mac")
	createCmd.Flags().String("ip", "", "Device IP address")
	createCmd.MarkFlagRequired("ip")

	var readCmd = &cobra.Command{
		Use:   "read",
		Short: "Read all devices",
		Run:   listDevices,
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing device",
		Run:   updateDevice,
	}
	updateCmd.Flags().String("uuid", "", "Device UUID")
	updateCmd.MarkFlagRequired("uuid")
	updateCmd.Flags().String("name", "", "Device name")
	updateCmd.Flags().String("hostname", "", "Device hostname")
	updateCmd.Flags().String("mac", "", "Device MAC address")
	updateCmd.Flags().String("ip", "", "Device IP address")

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a device",
		Run:   deleteDevice,
	}
	deleteCmd.Flags().String("uuid", "", "Device UUID")
	deleteCmd.MarkFlagRequired("uuid")

	rootCmd.AddCommand(createCmd, readCmd, updateCmd, deleteCmd)
	rootCmd.Execute()
}
