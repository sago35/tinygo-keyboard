package main

import (
	"context"
	_ "embed"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	usb.Product = "fric10key-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D9,
		machine.D8,
		machine.D7,
		machine.D6,
	}

	rowPins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
		machine.D5,
		machine.D10,
	}

	d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			jp.KeyF10, jp.KeyPageUp, jp.KeyUp, jp.KeyPageDown,
			jp.KeyF12, jp.KeyLeft, jp.KeyDown, jp.KeyRight,

			jp.KeyTab, jp.KeypadSlash, jp.KeypadAsterisk, jp.KeyBackspace,
			jp.Keypad7, jp.Keypad8, jp.Keypad9, jp.KeypadMinus,
			jp.Keypad4, jp.Keypad5, jp.Keypad6, jp.KeypadPlus,
			jp.Keypad1, jp.Keypad2, jp.Keypad3, jp.KeypadEnter,
			jp.Keypad0, jp.Key0, jp.KeypadPeriod, 0,
		},
		{
			jp.KeyF10, jp.KeyPageUp, jp.KeyUp, jp.KeyPageDown,
			jp.KeyF12, jp.KeyLeft, jp.KeyDown, jp.KeyRight,

			jp.KeyTab, jp.KeypadSlash, jp.KeypadAsterisk, jp.KeyBackspace,
			jp.Keypad7, jp.Keypad8, jp.Keypad9, jp.KeypadMinus,
			jp.Keypad4, jp.Keypad5, jp.Keypad6, jp.KeypadPlus,
			jp.Keypad1, jp.Keypad2, jp.Keypad3, jp.KeypadEnter,
			jp.Keypad0, jp.Key0, jp.KeypadPeriod, 0,
		},
	})

	// for Vial
	loadKeyboardDef()

	d.Debug = true
	return d.Loop(context.Background())
}
