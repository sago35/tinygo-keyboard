package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"runtime/volatile"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
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
	time.Sleep(2 * time.Second)

	pins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
	}

	sm := d.AddSquaredMatrixKeyboard(pins, [][][]keyboard.Keycode{
		{
			{jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD},
			{jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH},
			{jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyL},
		},
	})
	sm.SetCallback(func(layer, row, col int, state keyboard.State) {
		fmt.Printf("sm: %d %d %d %d\n", layer, row, col, state)
		rowx := row
		if col%2 == 1 {
			rowx = 3 - row - 1
		}
		wsLeds[rowx+3*col] = white
		if state == keyboard.PressToRelease {
			wsLeds[rowx+3*col] = black
		}
		fmt.Printf("%#v\n", wsLeds)
		changed.Set(1)
	})

	d.AddRotaryKeyboard(machine.D5, machine.D10, [][][]keyboard.Keycode{
		{
			{jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc},
		},
	})

	d.AddRotaryKeyboard(machine.D9, machine.D8, [][][]keyboard.Keycode{
		{
			{jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc},
		},
	})

	gpioPins := []machine.Pin{machine.D7, machine.D6}
	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	d.AddGpioKeyboard(gpioPins, [][][]keyboard.Keycode{
		{
			{jp.Key1, jp.Key2},
		},
	})

	cont := true
	ticker := time.Tick(32 * time.Millisecond)
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
