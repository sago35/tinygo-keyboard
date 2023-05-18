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
		machine.D8,
		machine.D9,
		machine.D10,
	}

	rowPins := []machine.Pin{
		machine.D1,
		machine.D2,
	}

	d.AddMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{jp.KeyT, jp.KeyI, jp.KeyN},
			{jp.KeyY, jp.KeyG, jp.KeyO},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
