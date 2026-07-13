//go:build xiao_ble

package main

import (
	"device/nrf"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
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
	led3.Low()

	nrf.USBD.USBPULLUP.Set(0) // ホストにまだ見せない
}

func callback(layer int) {
	led1.Set(layer != 0)
	led2.Set(layer != 0)
	led3.Set(layer == 0)
}

// setupKeyboard routes the keyboard and the joystick mouse through a
// USB/BLE switchable output.
func setupKeyboard(d *keyboard.Device) {
	kb := keyboard.NewBLEKeyboard()
	kb.MuteBLEOnUSB = true
	kb.OverrideCtrlH = true // override ctrl-h to BackSpace (USB/BLE 両方)
	d.Keyboard = kb
	d.Mouse = keyboard.NewBLEMouse(kb)
}

// enableUSB shows the device to the USB host. It runs after d.Init() so the
// host never sees a half-configured device (VID/PID and HID are set by then).
func enableUSB() {
	nrf.USBD.USBPULLUP.Set(1)
}
