package main

import (
	"context"
	"log"
	"machine"
	"time"

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
	time.Sleep(3 * time.Second)
	err := kbd.Load()
	if err != nil {
		return err
	}

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

	d.AddMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{jp.KeyF10, jp.KeyPageUp, jp.KeyUp, jp.KeyPageDown},
			{jp.KeyF12, jp.KeyLeft, jp.KeyDown, jp.KeyRight},

			{jp.KeyTab, jp.KeypadSlash, jp.KeypadAsterisk, jp.KeyBackspace},
			{jp.Keypad7, jp.Keypad8, jp.Keypad9, jp.KeypadMinus},
			{jp.Keypad4, jp.Keypad5, jp.Keypad6, jp.KeypadPlus},
			{jp.Keypad1, jp.Keypad2, jp.Keypad3, jp.KeypadEnter},
			{jp.Keypad0, jp.Key0, jp.KeypadPeriod},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
