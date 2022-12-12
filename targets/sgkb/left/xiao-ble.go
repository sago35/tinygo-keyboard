//go:build xiao_ble

package main

import (
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"tinygo.org/x/bluetooth"
)

var (
	adapter     = bluetooth.DefaultAdapter
	serviceUUID = bluetooth.NewUUID([16]byte{0x70, 0xb3, 0xc2, 0xef, 0x07, 0x85, 0x54, 0x30, 0x4c, 0x81, 0x86, 0xb6, 0xbc, 0xaf, 0x1c, 0x7a})
	charUUID    = bluetooth.NewUUID([16]byte{0x21, 0x96, 0x7c, 0xbf, 0x76, 0x9a, 0x64, 0x08, 0x60, 0x3b, 0x66, 0x5f, 0xbb, 0x68, 0xa7, 0xb0})
)

func initialize(d *keyboard.Device) error {
	go func() error {
		//return nil
		//time.Sleep(2 * time.Second)
		println("starting")
		must("enable BLE stack", adapter.Enable())
		// The address to connect to. Set during scanning and read afterwards.
		var foundDevice bluetooth.ScanResult

		// Scan for NUS peripheral.
		println("Scanning...")
		err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if result.LocalName() != "sgkb-0.1.0 server" {
				time.Sleep(100 * time.Millisecond)
				return
			}
			foundDevice = result

			// Stop the scan.
			err := adapter.StopScan()
			if err != nil {
				// Unlikely, but we can't recover from this.
				println("failed to stop the scan:", err.Error())
			}
		})
		if err != nil {
			return err
		}

		device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
		if err != nil {
			return err
		}

		println("Discovering service...")
		services, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
		if err != nil {
			return err
		}
		service := services[0]

		chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{charUUID})
		if err != nil {
			return err
		}

		d.Callback(func(layer int, down bool) {
			b := []byte{0}
			if down {
				b[0] = byte(layer)
			}
			chars[0].WriteWithoutResponse(b)
		})

		return nil
	}()

	return nil
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
