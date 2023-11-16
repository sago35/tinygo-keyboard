package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/ble"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/bluetooth"
)

func main() {
	usb.Product = "sgkey-ble-0.1.0"

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

	mk := d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			jp.KeyT, jp.KeyMediaPlay, jp.KeyMediaNextTrack,
			jp.KeyLeftShift, jp.KeyG, jp.KeyMod1,
		},
		{
			jp.KeyA, jp.KeyMediaVolumeInc, jp.KeyN,
			jp.KeyLeftShift, jp.KeyMediaVolumeDec, jp.KeyMod1,
		},
	})
	mk.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / len(colPins)
		col := index % len(colPins)
		fmt.Printf("mk: %d %d %d %d\n", layer, row, col, state)
	})

	d.AddBleSplitKeyboard(0x12, bluetooth.DefaultAdapter, "xiao-kb01-0.1.0", [][]keyboard.Keycode{
		{
			jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD,
			jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH,
			jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyL,

			jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc,

			jp.WheelDown, jp.WheelUp,

			jp.MouseLeft, jp.MouseRight,
		},
	})

	bk := ble.NewKeyboard(usb.Product)
	err := bk.Connect()
	if err != nil {
		return err
	}

	d.Keyboard = bk

	// for Vial
	loadKeyboardDef()

	d.Debug = true
	return d.Loop(context.Background())
}
