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

	kb := k.Port()

	layer := 0
	for {
		d.Get()

		for row := range d.State {
			for col := range d.State[row] {
				switch d.State[row][col] {
				case keyboard.None:
					// skip
				case keyboard.NoneToPress:
					if d.Keys[layer][row][col] == jp.KeyLeftAlt {
						layer = 1
					} else {
						kb.Down(d.Keys[layer][row][col])
					}
					fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[0][row][col])
				case keyboard.Press:
				case keyboard.PressToRelease:
					if d.Keys[layer][row][col] == jp.KeyLeftAlt {
						layer = 0
					} else {
						kb.Up(d.Keys[layer][row][col])
					}
					fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[0][row][col])
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
