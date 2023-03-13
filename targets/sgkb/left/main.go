package main

import (
	"context"
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
	d := keyboard.New([]machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
	}, []machine.Pin{
		machine.D5,
		machine.D10,
		machine.D9,
		machine.D8,
		machine.D6,
	}, [][][]keyboard.Keycode{
		{
			{0, 0, 0, 0, jp.Key1},
			{jp.KeyEsc, jp.KeyTab, jp.KeyLeftCtrl, jp.KeyLeftShift, jp.Key2},
			{jp.Key1, jp.KeyQ, jp.KeyA, jp.KeyZ, jp.KeyWindows},
			{jp.Key2, jp.KeyW, jp.KeyS, jp.KeyX, jp.KeyLeftAlt},
			{jp.Key3, jp.KeyE, jp.KeyD, jp.KeyC, jp.KeyMod1},
			{jp.Key4, jp.KeyR, jp.KeyF, jp.KeyV, jp.KeySpace},
			{jp.Key5, jp.KeyT, jp.KeyG, jp.KeyB},
			{jp.Key6},
			{},
			{},
		},
		{
			{0, 0, 0, 0, jp.Key1},
			{jp.KeyEsc, jp.KeyTab, jp.KeyLeftCtrl, jp.KeyLeftShift, jp.Key2},
			{jp.KeyF1, jp.KeyQ, jp.KeyHome, jp.KeyF13, jp.KeyWindows},
			{jp.KeyF2, jp.KeyF15, jp.KeyS, jp.KeyF14, jp.KeyLeftAlt},
			{jp.KeyF3, jp.KeyEnd, jp.MouseRight, jp.MouseMiddle, jp.KeyMod1},
			{jp.KeyF4, jp.KeyF17, jp.MouseLeft, jp.KeyF16, jp.KeySpace},
			{jp.KeyF5, jp.KeyF18, jp.KeyG, jp.KeyB},
			{jp.KeyF6},
			{},
			{},
		},
	})

	d.AddUartKeyboard(5, 5, [][][]keyboard.Keycode{
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
			{0, 0, 0, jp.KeyB},
			{jp.KeyF6, jp.KeyY, jp.KeyLeft, jp.WheelDown, jp.KeySpace},
			{jp.KeyF7, jp.KeyU, jp.KeyDown, jp.KeyM, jp.KeyHenkan},
			{jp.KeyF8, jp.KeyTab, jp.KeyUp, jp.KeyComma, jp.KeyMod1},
			{jp.KeyF9, jp.KeyO, jp.KeyRight, jp.KeyPeriod, jp.KeyLeftAlt},
			{jp.KeyF10, jp.WheelUp, jp.KeySemicolon, jp.KeySlash, jp.KeyPrintscreen},
			{jp.KeyF11, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyEnd},
			{jp.KeyF12, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyPageUp, jp.KeyPageDown},
			{jp.KeyBackslash2, jp.KeyEnter, 0, jp.KeyDelete, jp.KeyHome},
			{jp.KeyDelete},
		},
	})

	// 後で、いい感じの場所に移動する
	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.NoPin, RX: machine.UART_RX_PIN})

	// override ctrl-h to BackSpace
	d.OverrideCtrlH()

	return d.LoopUartRx(context.Background())
}
