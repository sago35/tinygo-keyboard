package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"machine"
	"machine/usb"
	"math/rand"
	"runtime/volatile"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/ble"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	usb.Product = "xiao-kb01-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	white = color.RGBA{0x3F, 0x3F, 0x3F, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

func run() error {
	var changed volatile.Register8
	changed.Set(0)

	neo := machine.D4
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	wsLeds := [12]color.RGBA{}
	for i := range wsLeds {
		wsLeds[i] = black
	}

	d := keyboard.New()

	pins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
	}

	sm := d.AddSquaredMatrixKeyboard(pins, [][]keyboard.Keycode{
		{
			0x0000, 0x0001, 0x0002, 0x0003,
			0x0004, 0x0005, 0x0006, 0x0007,
			0x0008, 0x0009, 0x000A, 0x000B,
		},
	})
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / 4
		col := index % 4
		fmt.Printf("sm: %d %d %d %d\n", layer, row, col, state)
		rowx := row
		if col%2 == 1 {
			rowx = 3 - row - 1
		}
		c := rand.Int()
		wsLeds[rowx+3*col] = color.RGBA{
			byte(c>>16) & 0x3F,
			byte(c>>8) & 0x3F,
			byte(c>>0) & 0x3F,
			0xFF,
		}
		if state == keyboard.PressToRelease {
			wsLeds[rowx+3*col] = black
		}
		fmt.Printf("%#v\n", wsLeds)
		changed.Set(1)
	})

	d.AddRotaryKeyboard(machine.D5, machine.D10, [][]keyboard.Keycode{
		{
			0x000C, 0x000D,
		},
	})

	d.AddRotaryKeyboard(machine.D9, machine.D8, [][]keyboard.Keycode{
		{
			0x000E, 0x000F,
		},
	})

	gpioPins := []machine.Pin{machine.D7, machine.D6}
	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			0x0010, 0x0011,
		},
	})

	time.Sleep(2 * time.Second)
	bk := ble.NewSplitKeyboard(usb.Product)
	err := bk.Connect()
	if err != nil {
		return err
	}

	d.Keyboard = bk

	// for Vial
	//loadKeyboardDef()

	err = d.Init()
	if err != nil {
		return err
	}

	cont := true
	ticker := time.Tick(4 * time.Millisecond)
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}
		if changed.Get() != 0 {
			ws.WriteColors(wsLeds[:])
			changed.Set(0)
		}
	}

	return nil
}
