package main

import (
	"context"
	_ "embed"
	"log"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/shifter"
)

func main() {
	usb.Product = "gobadge-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type RCS struct {
	row, col int
	state    keyboard.State
}

func run() error {
	d := keyboard.New()

	buttons := shifter.NewButtons()
	buttons.Configure()

	d.AddShifterKeyboard(buttons, [][]keyboard.Keycode{
		{
			jp.KeyT, // Left
			jp.KeyI, // Up
			jp.KeyN, // Down
			jp.KeyY, // Right
			jp.KeyG, // Select
			jp.KeyO, // Start
			jp.KeyA, // A
			jp.KeyB, // B
		},
	})

	// for Vial
	loadKeyboardDef()

	d.Debug = true
	return d.Loop(context.Background())
}
