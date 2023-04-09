package main

import (
	"context"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := keyboard.New()

	gpioPins := []machine.Pin{
		machine.KEY1,
		machine.KEY2,
		machine.KEY3,
		machine.KEY4,
		machine.KEY5,
		machine.KEY6,
		machine.KEY7,
		machine.KEY8,
		machine.KEY9,
		machine.KEY10,
		machine.KEY11,
		machine.KEY12,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	d.AddGpioKeyboard(gpioPins, [][][]keyboard.Keycode{
		{
			{
				jp.Key1,
				jp.Key2,
				jp.Key3,
				jp.Key4,
				jp.Key5,
				jp.Key6,
				jp.Key7,
				jp.Key8,
				jp.Key9,
				jp.KeyA,
				jp.KeyB,
				jp.KeyC,
			},
		},
	})

	d.AddRotaryKeyboard(machine.ROT_A, machine.ROT_B, [][][]keyboard.Keycode{
		{
			{jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
