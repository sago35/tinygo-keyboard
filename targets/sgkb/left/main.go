package main

import (
	"context"
	"log"
	"machine"
	k "machine/usb/hid/keyboard"

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
	d := keyboard.New(k.Port(), []machine.Pin{
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
	}, [][][]k.Keycode{
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
			{jp.Key1, jp.KeyQ, jp.KeyHome, jp.KeyZ, jp.KeyWindows},
			{jp.Key2, jp.KeyW, jp.KeyS, jp.KeyX, jp.KeyLeftAlt},
			{jp.Key3, jp.KeyEnd, jp.KeyD, jp.KeyC, jp.KeyMuhenkan},
			{jp.Key4, jp.KeyR, jp.KeyF, jp.KeyV, jp.KeySpace},
			{jp.Key5, jp.KeyT, jp.KeyG},
			{jp.Key6},
		},
	})

	return d.Loop(context.Background())
}
