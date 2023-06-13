package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
	})

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

	mk := d.AddMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{jp.KeyT, jp.KeyI, jp.KeyN},
			{jp.KeyY, jp.KeyG, jp.KeyO},
		},
	})
	mk.SetCallback(func(layer, row, col int, state keyboard.State) {
		fmt.Printf("mk: %d %d %d %d\n", layer, row, col, state)
		c := color.RGBA{255, 255, 255, 255}
		if state == keyboard.PressToRelease {
			c = color.RGBA{0, 0, 0, 255}
		}
		tinydraw.FilledRectangle(&display, 10+20*int16(col), 10+20*int16(row), 18, 18, c)
		display.Display()
	})

	d.Debug = true
	return d.Loop(context.Background())
}
