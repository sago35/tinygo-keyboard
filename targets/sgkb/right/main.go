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
	d := keyboard.New([]machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
	}, []machine.Pin{
		machine.D10,
		machine.D9,
		machine.D8,
		machine.D7,
	}, [][][]keyboard.Keycode{
		{
			{jp.KeyB, jp.KeyY, jp.KeyH, jp.KeyN},
			{jp.Key7, jp.KeyU, jp.KeyJ, jp.KeyM, jp.KeyHenkan},
			{jp.Key8, jp.KeyI, jp.KeyK, jp.KeyComma},
			{jp.Key9, jp.KeyO, jp.KeyL, jp.KeyPeriod},
			{jp.Key0, jp.KeyP, jp.KeySemicolon, jp.KeySlash},
			{jp.KeyMinus, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyLeft},
			{jp.KeyHat, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyUp, jp.KeyDown},
			{jp.KeyBackslash2, jp.KeyBackspace, jp.KeyEnter, 0, jp.KeyRight},
		},
		{
			{jp.KeyB, jp.KeyY, jp.KeyLeft, jp.WheelDown},
			{jp.KeyF7, jp.KeyU, jp.KeyDown, jp.KeyM, jp.KeyHenkan},
			{jp.KeyF8, jp.KeyTab, jp.KeyUp, jp.KeyComma},
			{jp.KeyF9, jp.KeyO, jp.KeyRight, jp.KeyPeriod},
			{jp.KeyF10, jp.WheelUp, jp.KeySemicolon, jp.KeySlash},
			{jp.KeyF11, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyHome},
			{jp.KeyF12, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyPageUp, jp.KeyPageDown},
			{jp.KeyBackslash2, jp.KeyDelete, jp.KeyEnter, 0, jp.KeyEnd},
		},
	})

	err := initialize(d)
	if err != nil {
		return err
	}

	return d.Loop(context.Background())
}
