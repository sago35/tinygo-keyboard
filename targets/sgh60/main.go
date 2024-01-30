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
	usb.Product = "sgh60-0.4.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	machine.InitADC()
	ax := machine.A2
	ay := machine.A3

	m := mouse.Port()

	d := keyboard.New()

	colPins := []machine.Pin{
		machine.GPIO15,
		machine.GPIO14,
		machine.GPIO13,
		machine.GPIO12,
		machine.GPIO11,
		machine.GPIO10,
		machine.GPIO9,
		machine.GPIO8,
		machine.GPIO7,
	}

	sm := d.AddSquaredMatrixKeyboard(colPins, [][]keyboard.Keycode{
		{
			jp.KeyEsc, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyBackspace,
			jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT, jp.KeyY, jp.KeyU, jp.KeyI, jp.KeyO, jp.KeyP, jp.KeyAt, jp.KeyLeftBrace, jp.KeyEnter,
			jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG, jp.KeyH, jp.KeyJ, jp.KeyK, jp.KeyL, jp.KeySemicolon, jp.KeyColon, jp.KeyRightBrace,
			jp.KeyLeftShift, jp.KeyZ, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB, jp.KeyN, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash, jp.KeyUp, jp.KeyDelete,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyLeft, jp.KeyDown, jp.KeyRight,
		},
		{
			jp.KeyEsc, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, jp.KeyF12, jp.KeyBackslash2, jp.KeyBackspace,
			jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, jp.KeyY, jp.KeyU, jp.KeyTab, jp.KeyO, jp.WheelUp, jp.KeyAt, jp.KeyLeftBrace, jp.KeyEnter,
			jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, jp.KeyLeft, jp.KeyDown, jp.KeyUp, jp.KeyRight, jp.KeyEnter, jp.KeyEsc, jp.KeyRightBrace,
			jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, jp.WheelDown, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash, jp.KeyPageUp, jp.KeyDelete,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyHome, jp.KeyPageDown, jp.KeyEnd,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			jp.KeyTab, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyBackspace, 0, 0,
			jp.KeyLeftCtrl, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyHome, jp.KeyPageDown, jp.KeyPageUp, jp.KeyEnd, jp.KeyEnter, jp.KeyEsc, 0,
			jp.KeyLeftShift, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, 0, 0,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyHiragana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyF12, 0, 0,
		},
	})
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		fmt.Printf("sm: %d %d %d\n", layer, index, state)
	})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	loadKeyboardDef()

	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	x := NewADCDevice(ax, 0x3400, 0xD400, true)
	y := NewADCDevice(ay, 0x3800, 0xE990, true)
	ticker := time.Tick(1 * time.Millisecond)
	cnt := 0
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}

		if cnt%10 == 0 {
			xx := x.Get2()
			yy := y.Get2()
			//fmt.Printf("%04X %04X %4d %4d %4d %4d\n", x.RawValue, y.RawValue, xx, yy, x.Get(), y.Get())
			m.Move(int(xx), int(yy))
		}
		cnt++
	}

	return nil
}
