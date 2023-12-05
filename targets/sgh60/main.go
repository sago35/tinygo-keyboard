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
	ax := machine.ADC{Pin: machine.A2}
	ax.Configure(machine.ADCConfig{})
	ay := machine.ADC{Pin: machine.A3}
	ay.Configure(machine.ADCConfig{})

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
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyKana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyLeft, jp.KeyDown, jp.KeyRight,
		},
		{
			jp.KeyEsc, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, jp.KeyF12, jp.KeyBackslash2, jp.KeyBackspace,
			jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, jp.KeyY, jp.KeyU, jp.KeyTab, jp.KeyO, jp.WheelUp, jp.KeyAt, jp.KeyLeftBrace, jp.KeyEnter,
			jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, jp.KeyLeft, jp.KeyDown, jp.KeyUp, jp.KeyRight, jp.KeyEnter, jp.KeyEsc, jp.KeyRightBrace,
			jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, jp.WheelDown, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash, jp.KeyPageUp, jp.KeyDelete,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyKana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyHome, jp.KeyPageDown, jp.KeyEnd,
		},
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			jp.KeyTab, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyBackspace, 0, 0,
			jp.KeyLeftCtrl, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyHome, jp.KeyPageDown, jp.KeyPageUp, jp.KeyEnd, jp.KeyEnter, jp.KeyEsc, 0,
			jp.KeyLeftShift, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, 0, 0,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, jp.KeySpace, jp.KeyMod2, jp.KeyKana, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyF12, 0, 0,
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
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}

		x := ax.Get()
		y := ay.Get()

		xx := ad2xy(x, 0x7E00)
		yy := ad2xy(y, 0x8600)
		m.Move(-1*xx, -1*yy)

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func ad2xy(x, initValue uint16) int {
	xx := int(x) - int(initValue)

	switch {
	case abs(xx) < 0x1000:
		return 0
	case abs(xx) < 0x2000:
		return 1 * sign(xx)
	case abs(xx) < 0x3000:
		return 2 * sign(xx)
	case abs(xx) < 0x3500:
		return 3 * sign(xx)
	case abs(xx) < 0x3A00:
		return 4 * sign(xx)
	case abs(xx) < 0x4000:
		return 5 * sign(xx)
	case abs(xx) < 0x4500:
		return 6 * sign(xx)
	case abs(xx) < 0x4A00:
		return 7 * sign(xx)
	case abs(xx) < 0x5000:
		return 8 * sign(xx)
	case abs(xx) < 0x5500:
		return 9 * sign(xx)
	case abs(xx) < 0x5A00:
		return 10 * sign(xx)
	case abs(xx) < 0x6000:
		return 11 * sign(xx)
	case abs(xx) < 0x6400:
		return 12 * sign(xx)
	case abs(xx) < 0x6800:
		return 13 * sign(xx)
	case abs(xx) < 0x6C00:
		return 14 * sign(xx)
	case abs(xx) < 0x7000:
		return 15 * sign(xx)
	}
	return 20 * sign(xx)
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}
