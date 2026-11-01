package main

import (
	"context"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	jp "github.com/sago35/tinygo-keyboard/keycodes/japanese"
)

func main() {
	usb.Product = "teensy40-0.1.0"

	d := keyboard.New()

	// Connect pin 2 to pin 5 to GND to press a key.
	// Pin 0 and pin 1 (UART1) stay free for debug output.
	gpioPins := []machine.Pin{
		machine.D2,
		machine.D3,
		machine.D4,
		machine.D5,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			jp.KeyA,
			jp.KeyB,
			jp.KeyC,
			jp.KeyD,
		},
	})

	// for Vial
	loadKeyboardDef()

	d.Loop(context.Background())
}
