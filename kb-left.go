//go:build kb_left

package main

import (
	"fmt"
	"machine"
	k "machine/usb/hid/keyboard"
	"time"
)

func run() error {
	wait := 1 * time.Millisecond

	d := New([]machine.Pin{
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
		{k.Key1, k.KeyQ, k.KeyA, k.KeyZ, k.KeyF21}, // Win がないのでいったん F21 へ
		{k.Key2, k.KeyW, k.KeyS, k.KeyX, k.KeyLeftAlt},
		{k.Key3, k.KeyE, k.KeyD, k.KeyC, k.KeyF20}, // Mod がない
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
				case None:
					// skip
				case NoneToPress:
					kb.Down(d.Keys[row][col])
					fmt.Printf("%2d %2d down\r\n", row, col)
				case Press:
				case PressToRelease:
					kb.Up(d.Keys[row][col])
					fmt.Printf("%2d %2d up\r\n", row, col)
				}
			}
		}

		time.Sleep(32*time.Millisecond - wait*3)
	}

	return nil
}
