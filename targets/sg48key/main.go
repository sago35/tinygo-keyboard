package main

import (
	_ "embed"
	"fmt"
	"log"
	"machine"
	"machine/usb"
	"machine/usb/hid/mouse"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	usb.Product = "sg48key-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	machine.InitADC()
	ax := machine.A0
	ay := machine.A1

	m := mouse.Port()

	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D2,
		machine.D3,
		machine.D6,
		machine.D7,
		machine.D8,
		machine.D9,
		machine.D10,
		machine.D4, // not connected
	}

	led1 := machine.LED
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2 := machine.LED2
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led3 := machine.LED3
	led3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	sm := d.AddSquaredMatrixKeyboard(colPins, [][]keyboard.Keycode{
		{
			jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT, jp.KeyY, jp.KeyU, jp.KeyI, jp.KeyO, jp.KeyP, jp.KeyAt,
			jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG, jp.KeyH, jp.KeyJ, jp.KeyK, jp.KeyL, jp.KeySemicolon, jp.KeyColon,
			jp.KeyLeftShift, jp.KeyZ, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB, jp.KeyN, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash,
			jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeyTo1, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyDelete,
		},
		{
			jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, jp.KeyY, jp.KeyU, jp.KeyTab, jp.KeyO, jp.WheelUp, jp.KeyAt,
			jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, jp.KeyLeft, jp.KeyDown, jp.KeyUp, jp.KeyRight, jp.KeyEnter, jp.KeyEsc,
			jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, jp.WheelDown, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash,
			jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeyTo2, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyDelete,
		},
		{
			jp.KeyTab, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyBackspace,
			jp.KeyLeftCtrl, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyHome, jp.KeyPageDown, jp.KeyPageUp, jp.KeyEnd, jp.KeyEnter, jp.KeyEsc,
			jp.KeyLeftShift, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11,
			jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeyTo0, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyF12,
		},
	})
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		layer = d.Layer()
		fmt.Printf("sm: %d %d %d\n", layer, index, state)
		led1.Set(layer != 0)
		led2.Set(layer == 0)
		led3.Set(layer == 0)
	})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	loadKeyboardDef()

	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	x := NewADCDevice(ax, 0x2000, 0xDC00, true)
	y := NewADCDevice(ay, 0x2400, 0xD400, true)
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}

		xx := x.Get2()
		yy := y.Get2()
		//fmt.Printf("%04X %04X %4d %4d %4d %4d\n", x.RawValue, y.RawValue, xx, yy, x.Get(), y.Get())
		m.Move(int(xx), int(yy))

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
