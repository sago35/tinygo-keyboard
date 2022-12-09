package main

import (
	"fmt"
	"machine"
	k "machine/usb/hid/keyboard"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
)

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
	}, [][]k.Keycode{
		{k.KeyEsc, k.KeyTab, k.KeyLeftCtrl, k.KeyLeftShift},
		{k.Key1, k.KeyQ, k.KeyA, k.KeyZ, keyboard.KeyWindows},
		{k.Key2, k.KeyW, k.KeyS, k.KeyX, k.KeyLeftAlt},
		{k.Key3, k.KeyE, k.KeyD, k.KeyC, keyboard.KeyMuhenkan},
		{k.Key4, k.KeyR, k.KeyF, k.KeyV, k.KeySpace},
		{k.Key5, k.KeyT, k.KeyG},
		{k.Key6},
	},
	)

	kb := k.Port()

	for {
		d.Get()

		for row := range d.State {
			for col := range d.State[row] {
				switch d.State[row][col] {
				case keyboard.None:
					// skip
				case keyboard.NoneToPress:
					kb.Down(d.Keys[row][col])
					fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[row][col])
				case keyboard.Press:
				case keyboard.PressToRelease:
					kb.Up(d.Keys[row][col])
					fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[row][col])
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
