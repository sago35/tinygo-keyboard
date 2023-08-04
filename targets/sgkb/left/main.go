package main

import (
	"context"
	"fmt"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
		machine.D5,
		machine.D8,
	}

	sm := d.AddSquaredMatrixKeyboard(colPins, [][][]keyboard.Keycode{
		{
			{jp.KeyEsc, jp.Key1, jp.Key2, jp.Key3, jp.Key4, jp.Key5, jp.Key6, 0},
			{jp.KeyTab, jp.KeyQ, jp.KeyW, jp.KeyE, jp.KeyR, jp.KeyT, 0},
			{jp.KeyLeftCtrl, jp.KeyA, jp.KeyS, jp.KeyD, jp.KeyF, jp.KeyG, 0},
			{jp.KeyLeftShift, jp.KeyZ, jp.KeyX, jp.KeyC, jp.KeyV, jp.KeyB, 0},
			{jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		},
		{
			{jp.KeyEsc, jp.KeyF1, jp.KeyF2, jp.KeyF3, jp.KeyF4, jp.KeyF5, jp.KeyF6, 0},
			{jp.KeyTab, jp.KeyQ, jp.KeyF15, jp.KeyEnd, jp.KeyF17, jp.KeyF18, 0},
			{jp.KeyLeftCtrl, jp.KeyHome, jp.KeyS, jp.MouseRight, jp.MouseLeft, jp.MouseBack, 0},
			{jp.KeyLeftShift, jp.KeyF13, jp.KeyF14, jp.MouseMiddle, jp.KeyF16, jp.MouseForward, 0},
			{jp.KeyMod1, jp.KeyLeftCtrl, jp.KeyWindows, jp.KeyLeftAlt, jp.KeyMod1, jp.KeySpace, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
		},
	})
	sm.SetCallback(func(layer, row, col int, state keyboard.State) {
		fmt.Printf("sm: %d %d %d %d\n", layer, row, col, state)
	})

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.NoPin, RX: machine.UART_RX_PIN})

	d.AddUartKeyboard(5, 5, uart, [][][]keyboard.Keycode{
		{
			{0, 0, 0, jp.KeyB},
			{jp.Key6, jp.KeyY, jp.KeyH, jp.KeyN, jp.KeySpace},
			{jp.Key7, jp.KeyU, jp.KeyJ, jp.KeyM, jp.KeyHenkan},
			{jp.Key8, jp.KeyI, jp.KeyK, jp.KeyComma, jp.KeyMod1},
			{jp.Key9, jp.KeyO, jp.KeyL, jp.KeyPeriod, jp.KeyLeftAlt},
			{jp.Key0, jp.KeyP, jp.KeySemicolon, jp.KeySlash, jp.KeyPrintscreen},
			{jp.KeyMinus, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyLeft},
			{jp.KeyHat, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyUp, jp.KeyDown},
			{jp.KeyBackslash2, jp.KeyEnter, 0, jp.KeyDelete, jp.KeyRight},
			{jp.KeyBackspace},
		},
		{
			{0, 0, 0, jp.MouseForward},
			{jp.KeyF6, jp.KeyY, jp.KeyLeft, jp.WheelDown, jp.KeySpace},
			{jp.KeyF7, jp.KeyU, jp.KeyDown, jp.KeyM, jp.KeyHenkan},
			{jp.KeyF8, jp.KeyTab, jp.KeyUp, jp.KeyComma, jp.KeyMod1},
			{jp.KeyF9, jp.KeyO, jp.KeyRight, jp.KeyPeriod, jp.KeyLeftAlt},
			{jp.KeyF10, jp.WheelUp, jp.KeySemicolon, jp.KeySlash, jp.KeyPrintscreen},
			{jp.KeyF11, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyHome},
			{jp.KeyF12, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyPageUp, jp.KeyPageDown},
			{jp.KeyBackslash2, jp.KeyEnter, 0, jp.KeyDelete, jp.KeyEnd},
			{jp.KeyDelete},
		},
	})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	return d.Loop(context.Background())
}
