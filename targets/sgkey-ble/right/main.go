package main

import (
	_ "embed"
	"log"
	"machine"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type RCS struct {
	row, col int
	state    keyboard.State
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D8,
		machine.D9,
		machine.D10,
	}

	rowPins := []machine.Pin{
		machine.D1,
		machine.D2,
	}

	d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			0x0000, 0x0001, 0x0002,
			0x0003, 0x0004, 0x0005,
		},
	})

	bleKeyboard := keyboard.BleTxKeyboard{
		RxBleName: "sgkey-left",
	}
	d.Keyboard = &bleKeyboard

	err := bleKeyboard.Connect()
	if err != nil {
		return err
	}

	err = d.Init()
	if err != nil {
		return err
	}

	cont := true
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
