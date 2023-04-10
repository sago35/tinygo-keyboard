package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/drivers/sh1106"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
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

	d.AddGpioKeyboard(gpioPins, [][][]keyboard.Keycode{
		{
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
		},
	})

	d.AddRotaryKeyboard(machine.ROT_A, machine.ROT_B, [][][]keyboard.Keycode{
		{
			{jp.KeyMediaVolumeDec, jp.KeyMediaVolumeInc},
		},
	})

	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 48000000,
	})
	display := sh1106.NewSPI(machine.SPI1, machine.OLED_DC, machine.OLED_RST, machine.OLED_CS)
	display.Configure(sh1106.Config{
		Width:  128,
		Height: 64,
	})
	display.ClearDisplay()

	err := d.Init()
	if err != nil {
		return err
	}

	ticker := time.Tick(10 * time.Millisecond)
	x := int16(0)
	y := int16(0)
	deltaX := int16(1)
	deltaY := int16(1)
	cont := true
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}

		pixel := display.GetPixel(x, y)
		c := color.RGBA{255, 255, 255, 255}
		if pixel {
			c = color.RGBA{0, 0, 0, 255}
		}
		display.SetPixel(x, y, c)
		display.Display()

		x += deltaX
		y += deltaY

		if x == 0 || x == 127 {
			deltaX = -deltaX
		}

		if y == 0 || y == 63 {
			deltaY = -deltaY
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

type celsius float32

func (c celsius) String() string {
	return fmt.Sprintf("%4.1f", c)
}
