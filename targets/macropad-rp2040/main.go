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
	"tinygo.org/x/drivers/sh1106"
	"tinygo.org/x/drivers/ws2812"
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

	neo := machine.WS2812
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	wsLeds := [12]color.RGBA{}
	for i := range wsLeds {
		wsLeds[i] = black
	}

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
		wsLeds[index] = white
		if state == keyboard.PressToRelease {
			wsLeds[index] = black
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

	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}
		ws.WriteColors(wsLeds[:])
		time.Sleep(1 * time.Millisecond)
	}

	return nil
}
