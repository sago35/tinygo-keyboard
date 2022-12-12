//go:build xiao_ble

package main

import (
	keyboard "github.com/sago35/tinygo-keyboard"
	"tinygo.org/x/bluetooth"
)

var (
	adapter     = bluetooth.DefaultAdapter
	serviceUUID = bluetooth.NewUUID([16]byte{0x70, 0xb3, 0xc2, 0xef, 0x07, 0x85, 0x54, 0x30, 0x4c, 0x81, 0x86, 0xb6, 0xbc, 0xaf, 0x1c, 0x7a})
	charUUID    = bluetooth.NewUUID([16]byte{0x21, 0x96, 0x7c, 0xbf, 0x76, 0x9a, 0x64, 0x08, 0x60, 0x3b, 0x66, 0x5f, 0xbb, 0x68, 0xa7, 0xb0})
)

func initialize(d *keyboard.Device) error {
	//time.Sleep(2 * time.Second)
	println("starting")
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "sgkb-0.1.0 server",
	}))
	must("start adv", adv.Start())

	var buf [4]byte
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  charUUID,
				Value: buf[:],
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 {
						return
					}

					if value[0] == 0 {
						d.Mod(1, false)
					} else {
						d.Mod(1, true)
					}
				},
			},
		},
	}))

	return nil
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
