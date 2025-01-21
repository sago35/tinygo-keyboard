//go:build xiao_ble

package main

import (
	"machine"
)

func init() {
	sclPin = machine.SCL0_PIN
	sdaPin = machine.SDA0_PIN
}
