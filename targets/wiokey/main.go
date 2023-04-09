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
			{jp.KeyT, jp.KeyY},
			{jp.KeyI, jp.KeyG},
			{jp.KeyN, jp.KeyO},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
