//go:build xiao

package main

import (
	"machine"
)

var (
	led1 = machine.LED
	led2 = machine.LED2
	led3 = machine.LED3
)

func init() {
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func callback(layer int) {
	led1.Set(layer != 0)
	led2.Set(layer == 0)
	led3.Set(layer == 0)
}
