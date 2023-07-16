package main

import (
	"context"
	"fmt"
	"log"
	"machine"

	kbd "machine/usb/hid/keyboard"

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
	err := kbd.Load()
	if err != nil {
		return err
	}

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

	mk := d.AddMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{jp.KeyT, jp.KeyI, jp.KeyN},
			{jp.KeyY, jp.KeyG, jp.KeyO},
		},
	})
	mk.SetCallback(func(layer, row, col int, state keyboard.State) {
		fmt.Printf("mk: %d %d %d %d\n", layer, row, col, state)
	})

	d.Debug = true
	return d.Loop(context.Background())
}
