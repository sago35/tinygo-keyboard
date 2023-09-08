package main

import (
	"context"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	d := keyboard.New()

	gpioPins := []machine.Pin{
		machine.D0,
		machine.D3,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			jp.KeyA,
			jp.KeyB,
		},
	})

	d.Loop(context.Background())
}
