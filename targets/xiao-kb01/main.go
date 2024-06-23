package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"machine"
	"machine/usb"
	"math/rand"
	"runtime/interrupt"
	"runtime/volatile"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
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
	red   = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

func run() error {
	var changed volatile.Register8
	changed.Set(0)

	neo := machine.D4
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812(s, neo)
	wsLeds := [12]color.RGBA{}
	for i := range wsLeds {
		wsLeds[i] = black
	}
	writeColors(s, ws, wsLeds[:])

	d := keyboard.New()

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
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		row := index / 4
		col := index % 4
		fmt.Printf("sm: %d %d %d %d\n", layer, row, col, state)
		rowx := row
		if col%2 == 1 {
			rowx = 3 - row - 1
		}
		c := rand.Int()
		mask := interrupt.Disable()
		wsLeds[rowx+3*col] = color.RGBA{
			byte(c>>16) & 0x3F,
			byte(c>>8) & 0x3F,
			byte(c>>0) & 0x3F,
			0xFF,
		}
		if state == keyboard.PressToRelease {
			//wsLeds[rowx+3*col] = black
		}
		interrupt.Restore(mask)
		fmt.Printf("%#v\n", wsLeds)
		changed.Set(1)
	})

	d.AddRotaryKeyboard(machine.D5, machine.D10, [][]keyboard.Keycode{
		{
			jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc,
		},
	})

	r2 := d.AddRotaryKeyboard(machine.D9, machine.D8, [][]keyboard.Keycode{
		{
			jp.WheelDown, jp.WheelUp,
		},
	})
	r2idx := 0
	r2.SetCallback(func(layer, index int, state keyboard.State) {
		if state == 2 {
			if index == 1 {
				r2idx = (r2idx + 1) % 10
			} else {
				r2idx = (r2idx - 1 + 10) % 10
			}
			redidx := r2idx
			switch r2idx {
			case 0:
				redidx = 0
			case 1:
				redidx = 1
			case 2:
				redidx = 2
			case 3:
				redidx = 3
			case 4:
				redidx = 8
			case 5:
				redidx = 9
			case 6:
				redidx = 10
			case 7:
				redidx = 11
			case 8:
				redidx = 6
			case 9:
				redidx = 5
			}
			fmt.Printf("r2: %d %d %d %d\n", layer, index, r2idx, state)
			mask := interrupt.Disable()
			wsLeds[redidx] = red
			interrupt.Restore(mask)
			changed.Set(1)
		}
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

	// for Vial
	loadKeyboardDef()

	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	ticker := time.Tick(1 * time.Millisecond)
	cnt := 0
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}
		if cnt%4 == 0 {
			if changed.Get() != 0 {
				writeColors(s, ws, wsLeds[:])
				changed.Set(0)
			}
		}
		if cnt%32 == 0 {
			mask := interrupt.Disable()
			for i, c := range wsLeds {
				c.B >>= 1
				c.R >>= 1
				c.G >>= 1
				wsLeds[i] = c
			}
			writeColors(s, ws, wsLeds[:])
			interrupt.Restore(mask)
		}
		cnt++
	}

	return nil
}

func writeColors(s pio.StateMachine, ws *piolib.WS2812, colors []color.RGBA) {
	for _, c := range colors {
		for s.IsTxFIFOFull() {
		}
		ws.SetColor(c)
	}
}
