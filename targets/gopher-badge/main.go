package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinydraw"
)

// This example is based on
// https://github.com/conejoninja/gopherbadge/blob/main/tutorial/basics/step5/main.go

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	circle = color.RGBA{0, 100, 250, 255}
	white  = color.RGBA{255, 255, 255, 255}
	ring   = color.RGBA{200, 0, 0, 255}
)

func run() error {
	// Setup the screen pins
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	})

	display := st7789.New(machine.SPI0,
		machine.TFT_RST,       // TFT_RESET
		machine.TFT_WRX,       // TFT_DC
		machine.TFT_CS,        // TFT_CS
		machine.TFT_BACKLIGHT) // TFT_LITE
	display.Configure(st7789.Config{
		Rotation: st7789.ROTATION_90,
		Height:   320,
	})
	fb := FB{}
	fb.FillScreen(white)

	// Draw blue circles to represent each of the buttons
	tinydraw.FilledCircle(&fb, 25, 120, 14, circle) // LEFT
	tinydraw.FilledCircle(&fb, 95, 120, 14, circle) // RIGHT
	tinydraw.FilledCircle(&fb, 60, 85, 14, circle)  // UP
	tinydraw.FilledCircle(&fb, 60, 155, 14, circle) // DOWN

	tinydraw.FilledCircle(&fb, 260, 120, 14, circle) // B
	tinydraw.FilledCircle(&fb, 295, 85, 14, circle)  // A

	display.DrawRGBBitmap8(0, 0, fb.buf[:], 320, 240)

	d := keyboard.New()

	gpioPins := []machine.Pin{
		machine.BUTTON_A,
		machine.BUTTON_B,
		machine.BUTTON_LEFT,
		machine.BUTTON_UP,
		machine.BUTTON_RIGHT,
		machine.BUTTON_DOWN,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInput})
	}

	gk := d.AddGpioKeyboard(gpioPins, [][][]keyboard.Keycode{
		{
			{
				jp.KeyA,
				jp.KeyB,
				jp.KeyLeft,
				jp.KeyUp,
				jp.KeyRight,
				jp.KeyDown,
			},
		},
	})

	xy := []struct{ x, y int16 }{
		{x: 295, y: 85},  // A
		{x: 260, y: 120}, // B
		{x: 25, y: 120},  // LEFT
		{x: 60, y: 85},   // UP
		{x: 95, y: 120},  // RIGHT
		{x: 60, y: 155},  // DOWN
	}
	gk.SetCallback(func(layer, row, col int, state keyboard.State) {
		fmt.Printf("gk: %d %d %d %d\n", layer, row, col, state)
		if state == keyboard.Press {
			tinydraw.Circle(&display, xy[col].x, xy[col].y, 16, ring)
		} else {
			tinydraw.Circle(&display, xy[col].x, xy[col].y, 16, white)
		}
	})

	return d.Loop(context.Background())
}