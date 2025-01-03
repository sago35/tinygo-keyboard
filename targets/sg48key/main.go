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
	jp "github.com/sago35/tinygo-keyboard/keycodes/japanese"
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

	sm := d.AddSquaredMatrixKeyboard(colPins, [][]keyboard.Keycode{
		{
			jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT, jp.KeyY, jp.KeyU, jp.KeyI, jp.KeyO, jp.KeyP, jp.KeyAt,
			jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG, jp.KeyH, jp.KeyJ, jp.KeyK, jp.KeyL, jp.KeySemicolon, jp.KeyColon,
			jp.KeyLeftShift, jp.KeyZ, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB, jp.KeyN, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash,
			jp.KeyEsc, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyTo1, jp.KeyPrintscreen, jp.KeyDelete,
		},
		{
			jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, jp.KeyY, jp.KeyU, jp.KeyTab, jp.KeyO, jp.WheelUp, jp.KeyAt,
			jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, jp.KeyLeft, jp.KeyDown, jp.KeyUp, jp.KeyRight, jp.KeyEnter, jp.KeyEsc,
			jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, jp.WheelDown, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash,
			jp.KeyEsc, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyTo2, jp.KeyPrintscreen, jp.KeyDelete,
		},
		{
			jp.KeyTab, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyBackspace,
			jp.KeyLeftCtrl, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyHome, jp.KeyPageDown, jp.KeyPageUp, jp.KeyEnd, jp.KeyEnter, jp.KeyEsc,
			jp.KeyLeftShift, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11,
			jp.KeyEsc, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyTo0, jp.KeyPrintscreen, jp.KeyF12,
		},
	})
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		layer = d.Layer()
		fmt.Printf("sm: %d %d %d\n", layer, index, state)
		callback(layer)
	})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	// Combos
	combos := []keyboard.Combo{
		{
			Keys:      [4]keyboard.Keycode{jp.KeyQ, jp.KeyZ},
			OutputKey: jp.KeyMediaMute,
		},
		{
			Keys:      [4]keyboard.Keycode{jp.KeyW, jp.KeyX},
			OutputKey: jp.KeyMediaVolumeDec,
		},
		{
			Keys:      [4]keyboard.Keycode{jp.KeyE, jp.KeyC},
			OutputKey: jp.KeyMediaVolumeInc,
		},
		{
			Keys:      [4]keyboard.Keycode{jp.KeyR, jp.KeyV},
			OutputKey: jp.KeyMediaBrightnessDown,
		},
		{
			Keys:      [4]keyboard.Keycode{jp.KeyT, jp.KeyB},
			OutputKey: jp.KeyMediaBrightnessUp,
		},
	}
	for i, c := range combos {
		d.SetCombo(i, c)
	}

	loadKeyboardDef()

	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	x := NewADCDevice(ax, 0x3000, 0xD000, true)
	y := NewADCDevice(ay, 0x3000, 0xD000, true)
	ticker := time.Tick(500 * time.Microsecond)
	cnt := 0
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}

		if cnt%(5*3) == 0 {
			xx := x.Get2()
			yy := y.Get2()
			//fmt.Printf("%04X %04X %4d %4d %4d %4d\n", x.RawValue, y.RawValue, xx, yy, x.Get(), y.Get())
			m.Move(int(xx), int(yy))
		}
		cnt++
	}

	return nil
}
