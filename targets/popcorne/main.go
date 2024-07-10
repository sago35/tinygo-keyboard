package main

import (
	"context"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"machine"
	"machine/usb"

	"tinygo.org/x/drivers/mcp23017"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinydraw"
)

func main() {
	usb.Product = "popcorne"

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
		Frequency: 400 * machine.KHz,
		SDA:       machine.P1_13,
		SCL:       machine.P0_29,
	})

	ch := make(chan RCS, 16)

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Width:  128,
		Height: 32,
	})
	display.ClearDisplay()

	expander, _ := mcp23017.NewI2C(machine.I2C0, 0) // TODO actual address

	d := keyboard.New()

	colPinsLeft := []machine.Pin{
		machine.P0_13,
		machine.P0_14,
		machine.P0_15,
		machine.P0_16,
		machine.P0_24,
		machine.P0_25,
	}

	rowPinsLeft := []machine.Pin{
		machine.P0_03,
		machine.P0_28,
		machine.P0_02,
		machine.P0_31,
	}

	var colPinsRight []mcp23017.Pin

	var rowPinsRight []mcp23017.Pin

	mkLeft := d.AddMatrixKeyboard(colPinsLeft, rowPinsLeft, [][]keyboard.Keycode{
		{
			jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT,
			jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG,
			jp.KeyLeftShift, jp.KeyY, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB,
			0, jp.KeyMod1, jp.KeyMod2, jp.KeyWindows, jp.KeySpace, 0,
		},
		{
			jp.KeyEsc, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5,
			0, 0, jp.KeyUp, 0, 0, 0,
			0, jp.KeyLeft, jp.KeyDown, jp.KeyRight, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
		{
			jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
	}, keyboard.InvertDiode(true))
	mkLeft.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / len(colPinsLeft)
		col := index % len(colPinsLeft)
		fmt.Printf("mkLeft: %d %d %d %d\n", layer, row, col, state)
		select {
		case ch <- RCS{row: row, col: col, state: state}:
		}
	})

	mkRight := d.AddExpanderKeyboard(expander, colPinsRight, rowPinsRight, [][]keyboard.Keycode{
		{
			jp.KeyZ, jp.KeyU, jp.KeyI, jp.KeyO, jp.KeyP, jp.KeyMinus,
			jp.KeyH, jp.KeyJ, jp.KeyK, jp.KeyL, jp.KeyColon, jp.KeySemicolon,
			jp.KeyN, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyRightShift,
			jp.KeyEnter, jp.KeyMod3, jp.KeyBackspace, jp.KeyDelete, 0, 0,
		},
		{
			jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyHat,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
		{
			jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, jp.KeyF12,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		},
	}, keyboard.InvertDiode(true))
	mkRight.SetCallback(func(layer, index int, state keyboard.State) {
		row := index/len(colPinsRight) + len(rowPinsLeft)
		col := index%len(colPinsRight) + len(colPinsLeft)
		fmt.Printf("mkRight: %d %d %d %d\n", layer, row, col, state)
		select {
		case ch <- RCS{row: row, col: col, state: state}:
		}
	})

	go func() {
		for {
			select {
			case x := <-ch:
				c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
				if x.state == keyboard.PressToRelease {
					c = color.RGBA{A: 255}
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
