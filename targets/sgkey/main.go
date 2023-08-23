package main

import (
	"context"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
)

func main() {
	usb.Product = "sgkey-0.1.0"

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
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})

	ch := make(chan RCS, 16)

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.ClearDisplay()

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
			jp.KeyT, jp.KeyI, jp.KeyN,
			jp.KeyY, jp.KeyG, jp.KeyO,
		},
	})
	mk.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / len(colPins)
		col := index % len(colPins)
		fmt.Printf("mk: %d %d %d %d\n", layer, row, col, state)
		select {
		case ch <- RCS{row: row, col: col, state: state}:
		}
	})

	go func() {
		for {
			select {
			case x := <-ch:
				c := color.RGBA{255, 255, 255, 255}
				if x.state == keyboard.PressToRelease {
					c = color.RGBA{0, 0, 0, 255}
				}
				tinydraw.FilledRectangle(&display, 10+20*int16(x.col), 10+20*int16(x.row), 18, 18, c)
				display.Display()
			}
		}
	}()

	// for Vial
	loadKeyboardDef()

	d.Debug = true
	return d.Loop(context.Background())
}
