package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"math/rand"
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
	time.Sleep(2 * time.Second)

	pins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
	}

	sm := d.AddSquaredMatrixKeyboard(pins, [][]keyboard.Keycode{
		{
			jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD,
			jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH,
			jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyL,
		},
	})
	sm.SetCallback(func(layer, row, col int, state keyboard.State) {
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
			jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc,
		},
	})

	d.AddRotaryKeyboard(machine.D9, machine.D8, [][]keyboard.Keycode{
		{
			jp.WheelDown, jp.WheelUp,
		},
	})

	gpioPins := []machine.Pin{machine.D7, machine.D6}
	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	d.AddGpioKeyboard(gpioPins, [][]keyboard.Keycode{
		{
			jp.MouseLeft, jp.MouseRight,
		},
	})

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
