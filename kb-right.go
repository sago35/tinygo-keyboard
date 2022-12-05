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
		{k.Key7, k.KeyU, k.KeyJ, k.KeyM, k.KeyF22 /* M2 */}, // M2 がない
		{k.Key8, k.KeyI, k.KeyK, k.KeyComma},
		{k.Key9, k.KeyO, k.KeyL, k.KeyPeriod},
		{k.Key0, k.KeyP, k.KeySemicolon, k.KeySlash},
		{k.KeyMinus, k.KeyLeftBrace /* @ */, 0xF000 | 52 /* : */, k.KeyF23 /* \ */, k.KeyLeft},
		{0xF000 | 46 /* ^ */, k.KeyRightBrace /* [ */, k.KeyBackslash /* ] */, k.KeyUp, k.KeyDown},
		{k.KeyF23 /* \ */, k.KeyBackspace, k.KeyEnter, 0, k.KeyRight},
	},
	)
	// @ : KeyLeftBrace
	// [ : KeyBackslash

	kb := k.Port()

	code := k.Keycode(0)
	for {
		d.Get()

		for row := range d.State {
			for col := range d.State[row] {
				switch d.State[row][col] {
				case None:
					// skip
				case NoneToPress:
					if false {
						kb.Press(0xF000 | code)
						fmt.Printf("code %d\r\n", int(code))
						code++
					} else {
						kb.Down(d.Keys[row][col])
						fmt.Printf("%2d %2d down\r\n", row, col)
					}
				case Press:
				case PressToRelease:
					if false {
					} else {
						kb.Up(d.Keys[row][col])
						fmt.Printf("%2d %2d up\r\n", row, col)
					}
				}
			}
		}

		time.Sleep(32*time.Millisecond - wait*3)
	}

	return nil
}
