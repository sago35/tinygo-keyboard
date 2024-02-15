//go:build xiao_rp2040

package main

import (
	"machine"
)

var (
	led1 = machine.LED_RED
	led2 = machine.LED_GREEN
	led3 = machine.LED_BLUE
)

func init() {
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	led1.Low()
	led2.Low()
	led3.High()
}

func callback(layer int) {
	led1.Set(layer != 0)
	led2.Set(layer != 0)
	led3.Set(layer == 0)
}
