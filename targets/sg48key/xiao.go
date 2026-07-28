//go:build xiao

package main

import (
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
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

// setupKeyboard keeps the default USB-only output; Ctrl+H is translated to
// Backspace.
func setupKeyboard(d *keyboard.Device) {
	d.OverrideCtrlH()
}

// enableUSB is a no-op: on this target USB is visible from reset.
func enableUSB() {
}

// tickBoard is a no-op: this target has no board-specific periodic work.
func tickBoard(cnt int) {
}

// allowIdle disables the idle (slow scan) mode: this target is USB powered,
// so slowing down would only add latency.
func allowIdle() bool {
	return false
}
