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
			{jp.KeyEsc, jp.KeyTab, jp.KeyLeftCtrl, jp.KeyLeftShift},
			{jp.Key1, jp.KeyQ, jp.KeyA, jp.KeyZ, jp.KeyWindows},
			{jp.Key2, jp.KeyW, jp.KeyS, jp.KeyX, jp.KeyLeftAlt},
			{jp.Key3, jp.KeyE, jp.KeyD, jp.KeyC, jp.KeyMuhenkan},
			{jp.Key4, jp.KeyR, jp.KeyF, jp.KeyV, jp.KeySpace},
			{jp.Key5, jp.KeyT, jp.KeyG},
			{jp.Key6},
		},
		{
			{jp.KeyEsc, jp.KeyTab, jp.KeyLeftCtrl, jp.KeyLeftShift},
			{jp.KeyF1, jp.KeyQ, jp.KeyHome, jp.KeyF13, jp.KeyWindows},
			{jp.KeyF2, jp.KeyF15, jp.KeyS, jp.KeyF14, jp.KeyLeft},
			{jp.KeyF3, jp.KeyEnd, jp.MouseRight, jp.MouseMiddle, jp.KeyMuhenkan},
			{jp.KeyF4, jp.KeyF17, jp.MouseLeft, jp.KeyF16, jp.KeySpace},
			{jp.KeyF5, jp.KeyF18, jp.KeyG},
			{jp.KeyF6},
		},
	})

	err := initialize(d)
	if err != nil {
		return err
	}

	return d.Loop(context.Background())
}
