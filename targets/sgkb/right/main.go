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
			{jp.KeyB, jp.KeyY, jp.KeyH, jp.KeyN},
			{jp.Key7, jp.KeyU, jp.KeyJ, jp.KeyM, jp.KeyHenkan},
			{jp.Key8, jp.KeyI, jp.KeyK, jp.KeyComma},
			{jp.Key9, jp.KeyO, jp.KeyL, jp.KeyPeriod},
			{jp.Key0, jp.KeyP, jp.KeySemicolon, jp.KeySlash},
			{jp.KeyMinus, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyLeft},
			{jp.KeyHat, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyUp, jp.KeyDown},
			{jp.KeyBackslash2, jp.KeyBackspace, jp.KeyEnter, 0, jp.KeyRight},
		},
	})

	return d.Loop(context.Background())
}
