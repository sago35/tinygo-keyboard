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

	colPins := []machine.Pin{
		machine.BCM5,
		machine.BCM6,
	}

	rowPins := []machine.Pin{
		machine.BCM13,
		machine.BCM19,
		machine.BCM26,
	}

	d.AddDuplexMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{jp.KeyT, jp.KeyI},
			{jp.KeyN, jp.KeyY},
			{jp.KeyG, jp.KeyO},
		},
	})

	gpioPins := []machine.Pin{
		machine.WIO_KEY_A,
		machine.WIO_KEY_B,
		machine.WIO_KEY_C,
		machine.WIO_5S_UP,
		machine.WIO_5S_LEFT,
		machine.WIO_5S_RIGHT,
		machine.WIO_5S_DOWN,
		machine.WIO_5S_PRESS,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInput})
	}

	d.AddGpioKeyboard(gpioPins, [][][]keyboard.Keycode{
		{
			{jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyMediaVolumeInc, jp.KeyMediaPrevTrack, jp.KeyMediaNextTrack, jp.KeyMediaVolumeDec, jp.KeyMediaPlayPause},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
