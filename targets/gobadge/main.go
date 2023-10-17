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
			jp.Key0,   // Left
			jp.Key0,   // Up
			jp.Key0,   // Down
			jp.Key0,   // Right
			jp.KeyTo1, // Select
			jp.KeyTo1, // Start
			jp.KeyA,   // A
			jp.KeyB,   // B
		},
		{
			jp.Key1,    // Left
			jp.Key1,    // Up
			jp.Key1,    // Down
			jp.Key1,    // Right
			jp.KeyTo2,  // Select
			jp.KeyTo2,  // Start
			jp.KeyMod3, // A
			jp.KeyB,    // B
		},
		{
			jp.Key2,    // Left
			jp.Key2,    // Up
			jp.Key2,    // Down
			jp.Key2,    // Right
			jp.KeyTo0,  // Select
			jp.KeyTo0,  // Start
			jp.KeyMod4, // A
			jp.KeyB,    // B
		},
		{
			jp.Key3,    // Left
			jp.Key3,    // Up
			jp.Key3,    // Down
			jp.Key3,    // Right
			jp.Key3,    // Select
			jp.Key3,    // Start
			jp.KeyMod3, // A
			jp.KeyB,    // B
		},
		{
			jp.Key4,    // Left
			jp.Key4,    // Up
			jp.Key4,    // Down
			jp.Key4,    // Right
			jp.Key4,    // Select
			jp.Key4,    // Start
			jp.KeyMod4, // A
			jp.KeyB,    // B
		},
	})

	// for Vial
	loadKeyboardDef()

	d.Debug = true
	return d.Loop(context.Background())
}
