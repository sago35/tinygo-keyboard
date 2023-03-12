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
		machine.D7,
	}, [][][]keyboard.Keycode{
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
			{jp.KeyF6, jp.KeyY, jp.KeyLeft, jp.KeyN, jp.KeySpace},
			{jp.KeyF7, jp.KeyU, jp.KeyDown, jp.KeyM, jp.KeyHenkan},
			{jp.KeyF8, jp.KeyI, jp.KeyUp, jp.KeyComma, jp.KeyMod1},
			{jp.KeyF9, jp.KeyO, jp.KeyRight, jp.KeyPeriod, jp.KeyLeftAlt},
			{jp.KeyF10, jp.KeyP, jp.KeySemicolon, jp.KeySlash, jp.KeyPrintscreen},
			{jp.KeyF11, jp.KeyAt, jp.KeyColon, jp.KeyBackslash, jp.KeyLeft},
			{jp.KeyF12, jp.KeyLeftBrace, jp.KeyRightBrace, jp.KeyUp, jp.KeyDown},
			{jp.KeyBackslash2, jp.KeyEnter, 0, jp.KeyDelete, jp.KeyRight},
			{jp.KeyDelete},
		},
	})

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.NoPin})

	return d.LoopUartTx(context.Background())
}
