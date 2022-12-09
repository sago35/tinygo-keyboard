package main

import (
	"fmt"
	"log"
	"machine"
	k "machine/usb/hid/keyboard"
	"time"

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
	}, [][]k.Keycode{
		{jp.KeyB, jp.KeyY, jp.KeyH, jp.KeyN},
		{jp.Key7, jp.KeyU, jp.KeyJ, jp.KeyM, jp.KeyHenkan},
		{jp.Key8, jp.KeyI, jp.KeyK, jp.KeyComma},
		{jp.Key9, jp.KeyO, jp.KeyL, jp.KeyPeriod},
		{jp.Key0, jp.KeyP, jp.KeySemicolon, jp.KeySlash},
		{jp.KeyMinus, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyLeft},
		{jp.KeyHat, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyUp, jp.KeyDown},
		{jp.KeyBackslash2, jp.KeyBackspace, jp.KeyEnter, 0, jp.KeyRight},
	})

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
