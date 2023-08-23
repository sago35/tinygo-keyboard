package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	usb.Product = "sgkb-0.3.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D5,
		machine.D10,
		machine.D9,
		machine.D8,
		machine.D6,
	}

	rowPins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
	}

	dk := d.AddDuplexMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			0, jp.KeyEsc, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, 0, 0,
			0, jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT, 0, 0, 0,
			0, jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG, 0, 0, 0,
			0, jp.KeyLeftShift, jp.KeyZ, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB, 0, 0, 0,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, 0, 0, 0,
		},
		{
			0, jp.KeyEsc, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, 0, 0,
			0, jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, 0, 0, 0,
			0, jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, 0, 0, 0,
			0, jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, 0, 0, 0,
			jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, 0, 0, 0,
		},
	})
	dk.SetCallback(func(layer, index int, state keyboard.State) {
		fmt.Printf("dk: %d %d %d\n", layer, index, state)
	})

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.NoPin, RX: machine.UART_RX_PIN})

	uk := d.AddUartKeyboard(50, uart, [][]keyboard.Keycode{
		{
			0, jp.Key6, jp.Key7, jp.Key8, jp.Key9, jp.Key0, jp.KeyMinus, jp.KeyHat, jp.KeyBackslash2, jp.KeyBackspace,
			0, jp.KeyY, jp.KeyU, jp.KeyI, jp.KeyO, jp.KeyP, jp.KeyAt, jp.KeyLeftBrace, jp.KeyEnter, 0,
			0, jp.KeyH, jp.KeyJ, jp.KeyK, jp.KeyL, jp.KeySemicolon, jp.KeyColon, jp.KeyRightBrace, 0, 0,
			jp.KeyB, jp.KeyN, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash, jp.KeyUp, jp.KeyDelete, 0,
			0, jp.KeySpace, jp.KeyHenkan, jp.KeyMod1, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyLeft, jp.KeyDown, jp.KeyRight, 0,
		},
		{
			0, jp.KeyF6, jp.KeyF7, jp.KeyF8, jp.KeyF9, jp.KeyF10, jp.KeyF11, jp.KeyF12, jp.KeyBackslash2, jp.KeyBackspace,
			0, jp.KeyY, jp.KeyU, jp.KeyTab, jp.KeyO, jp.WheelUp, jp.KeyAt, jp.KeyLeftBrace, jp.KeyEnter, 0,
			0, jp.KeyLeft, jp.KeyDown, jp.KeyUp, jp.KeyRight, jp.KeySemicolon, jp.KeyColon, jp.KeyRightBrace, 0, 0,
			jp.MouseForward, jp.WheelDown, jp.KeyM, jp.KeyComma, jp.KeyPeriod, jp.KeySlash, jp.KeyBackslash, jp.KeyPageUp, jp.KeyDelete, 0,
			0, jp.KeySpace, jp.KeyHenkan, jp.KeyMod1, jp.KeyLeftAlt, jp.KeyPrintscreen, jp.KeyHome, jp.KeyPageDown, jp.KeyEnd, 0,
		},
	})
	uk.SetCallback(func(layer, index int, state keyboard.State) {
		fmt.Printf("uk: %d %d %d\n", layer, index, state)
	})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	loadKeyboardDef()

	return d.Loop(context.Background())
}
