package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"machine"
	"machine/usb"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers/sh1106"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func main() {
	usb.Product = "macropad-rp2040-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

const (
	rawWhite = 0x3F3F3FFF
	rawRed   = 0x00FF00FF
	rawGreen = 0xFF0000FF
	rawBlue  = 0x0000FFFF
	rawBlack = 0x000000FF
)

func writeColors(s pio.StateMachine, ws *piolib.WS2812B, colors []uint32) {
	ws.WriteRaw(colors)
}

func run() error {
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 48000000,
	})
	display := sh1106.NewSPI(machine.SPI1, machine.OLED_DC, machine.OLED_RST, machine.OLED_CS)
	display.Configure(sh1106.Config{
		Width:  128,
		Height: 64,
	})
	display.ClearDisplay()

	wsPin := machine.WS2812
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, wsPin)
	err := ws.EnableDMA(true)
	if err != nil {
		return err
	}
	wsLeds := [12]uint32{}
	for i := range wsLeds {
		wsLeds[i] = rawBlack
	}
	writeColors(s, ws, wsLeds[:])

	d := keyboard.New()

	gpioPins := []machine.Pin{
		machine.KEY1,
		machine.KEY2,
		machine.KEY3,
		machine.KEY4,
		machine.KEY5,
		machine.KEY6,
		machine.KEY7,
		machine.KEY8,
		machine.KEY9,
		machine.KEY10,
		machine.KEY11,
		machine.KEY12,
	}

	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	gk := d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			jp.Key1,
			jp.Key2,
			jp.Key3,
			jp.Key4,
			jp.Key5,
			jp.Key6,
			jp.Key7,
			jp.Key8,
			jp.Key9,
			jp.KeyA,
			jp.KeyB,
			jp.KeyC,
		},
	})
	gk.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / 3
		col := index % 3
		fmt.Printf("gk: %d %d %d %d %d\n", layer, index, row, col, state)
		c := white
		wsLeds[index] = rawWhite
		if state == keyboard.PressToRelease {
			wsLeds[index] = rawBlack
			c = black
		}
		display.ClearBuffer()
		tinyfont.WriteLine(&display, &freemono.Regular9pt7b, 10, 20, fmt.Sprintf("Key%d", index+1), c)
		display.Display()
	})

	rk := d.AddRotaryKeyboard(machine.ROT_A, machine.ROT_B, [][]keyboard.Keycode{
		{
			jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc,
		},
	})
	rk.SetCallback(func(layer, index int, state keyboard.State) {
		fmt.Printf("rk: %d %d %d\n", layer, index, state)
	})

	// for Vial
	loadKeyboardDef()

	err = d.Init()
	if err != nil {
		return err
	}

	// Combos: A + B = J
	d.Combos[1][0] = 0x0004
	d.Combos[1][1] = 0x0005
	d.Combos[1][2] = 0x0000
	d.Combos[1][3] = 0x0000
	d.Combos[1][4] = 0x000D

	// Combos: A + B + C = K
	d.Combos[0][0] = 0x0004
	d.Combos[0][1] = 0x0005
	d.Combos[0][2] = 0x0006
	d.Combos[0][3] = 0x0000
	d.Combos[0][4] = 0x000E

	cont := true
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}
		writeColors(s, ws, wsLeds[:])
		time.Sleep(1 * time.Millisecond)
	}

	return nil
}
