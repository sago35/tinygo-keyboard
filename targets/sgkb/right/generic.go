//go:build !xiao_ble

package main

import keyboard "github.com/sago35/tinygo-keyboard"

func initialize(d *keyboard.Device) error {
	return nil
}
