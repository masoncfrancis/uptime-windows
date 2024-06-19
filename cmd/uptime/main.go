package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	// Initialize COM
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Create WMI COM object
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		log.Fatalf("Failed to create WbemScripting.SWbemLocator object: %v", err)
	}
	defer unknown.Release()

	// Query the WMI service
	wmiService, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatalf("Failed to query interface: %v", err)
	}
	defer wmiService.Release()

	// Connect to the WMI namespace
	wmi, err := oleutil.CallMethod(wmiService, "ConnectServer", nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WMI namespace: %v", err)
	}
	defer wmi.Clear()

	// Query the Win32_OperatingSystem class
	result, err := oleutil.CallMethod(wmi.ToIDispatch(), "ExecQuery", "SELECT LastBootUpTime FROM Win32_OperatingSystem")
	if err != nil {
		log.Fatalf("Failed to query Win32_OperatingSystem class: %v", err)
	}
	defer result.Clear()

	// Iterate over the result set
	countVariant, err := oleutil.GetProperty(result.ToIDispatch(), "Count")
	if err != nil {
		log.Fatalf("Failed to get count: %v", err)
	}
	count := int(countVariant.Val)

	for i := 0; i < count; i++ {
		itemVariant, err := oleutil.CallMethod(result.ToIDispatch(), "ItemIndex", i)
		if err != nil {
			log.Fatalf("Failed to get item at index %d: %v", i, err)
		}
		item := itemVariant.ToIDispatch()
		defer item.Release()

		lastBootUpTimeVariant, err := oleutil.GetProperty(item, "LastBootUpTime")
		if err != nil {
			log.Fatalf("Failed to get LastBootUpTime: %v", err)
		}
		lastBootUpTime := lastBootUpTimeVariant.ToString()

		// Remove the time zone offset and handle the boot time separately
		timePart := lastBootUpTime[:len("20060102150405.000000")]
		offsetPart := lastBootUpTime[len("20060102150405.000000"):]

		// Parse the LastBootUpTime without the offset
		bootTime, err := time.Parse("20060102150405.000000", timePart)
		if err != nil {
			log.Fatalf("Failed to parse boot time: %v", err)
		}

		// Handle the time zone offset if it exists
		if strings.TrimSpace(offsetPart) != "" {
			offsetMinutes := 0
			fmt.Sscanf(offsetPart, "%d", &offsetMinutes)
			bootTime = bootTime.Add(time.Duration(-offsetMinutes) * time.Minute)
		}

		// Calculate uptime
		uptime := time.Since(bootTime)

		// Print the result in a nice format
		printUptime(uptime)
	}
}

func printUptime(uptime time.Duration) {
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	fmt.Printf("System Uptime: %d days, %02d hours, %02d minutes, %02d seconds\n", days, hours, minutes, seconds)
}
