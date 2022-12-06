//go:build kb_right

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
		{k.KeyB, k.KeyY, k.KeyH, k.KeyN},
		{k.Key7, k.KeyU, k.KeyJ, k.KeyM, KeyHenkan},
		{k.Key8, k.KeyI, k.KeyK, k.KeyComma},
		{k.Key9, k.KeyO, k.KeyL, k.KeyPeriod},
		{k.Key0, k.KeyP, k.KeySemicolon, k.KeySlash},
		{k.KeyMinus, KeyAt, KeyColon, KeyBackslash, k.KeyLeft},
		{KeyHat, KeyLeftBrace, KeyRightBrace, k.KeyUp, k.KeyDown},
		{KeyBackslash2, k.KeyBackspace, k.KeyEnter, 0, k.KeyRight},
	},
	)
	// @ : KeyLeftBrace
	// [ : KeyBackslash

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
					fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[row][col])
				case Press:
				case PressToRelease:
					kb.Up(d.Keys[row][col])
					fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[row][col])
				}
			}
		}

		time.Sleep(32*time.Millisecond - wait*3)
	}

	return nil
}
