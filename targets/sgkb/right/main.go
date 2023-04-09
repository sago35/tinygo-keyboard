package main

import (
	"context"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
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
	}

	rowPins := []machine.Pin{
		machine.D5,
		machine.D10,
		machine.D9,
		machine.D8,
		machine.D7,
	}

	d.AddDuplexMatrixKeyboard(colPins, rowPins, [][][]keyboard.Keycode{
		{
			{0x0000, 0x0001, 0x0002, 0x0003, 0x0004},
			{0x0100, 0x0101, 0x0102, 0x0103, 0x0104},
			{0x0200, 0x0201, 0x0202, 0x0203, 0x0204},
			{0x0300, 0x0301, 0x0302, 0x0303, 0x0304},
			{0x0400, 0x0401, 0x0402, 0x0403, 0x0404},
			{0x0500, 0x0501, 0x0502, 0x0503, 0x0504},
			{0x0600, 0x0601, 0x0602, 0x0603, 0x0604},
			{0x0700, 0x0701, 0x0702, 0x0703, 0x0704},
			{0x0800, 0x0801, 0x0802, 0x0803, 0x0804},
			{0x0900, 0x0901, 0x0902, 0x0903, 0x0904},
		},
	})

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.NoPin})

	d.Keyboard = &keyboard.UartTxKeyboard{
		Uart: uart,
	}

	return d.Loop(context.Background())
}
