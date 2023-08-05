package main

import (
	"context"
	"fmt"
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
		machine.D5,
		machine.D10,
		machine.D9,
		machine.D8,
		machine.D7,
	}

	rowPins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
		machine.D4,
	}

	dm := d.AddDuplexMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		{
			0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007, 0x0008, 0x0009,
			0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F, 0x0010, 0x0011, 0x0012, 0x0013,
			0x0014, 0x0015, 0x0016, 0x0017, 0x0018, 0x0019, 0x001A, 0x001B, 0x001C, 0x001D,
			0x001E, 0x001F, 0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
			0x0028, 0x0029, 0x002A, 0x002B, 0x002C, 0x002D, 0x002E, 0x002F, 0x0030, 0x0031,
		},
	})
	dm.SetCallback(func(layer, index int, state keyboard.State) {
		fmt.Printf("dm: %d %d %d\n", layer, index, state)
	})

	uart := machine.UART0
	uart.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.NoPin})

	d.Keyboard = &keyboard.UartTxKeyboard{
		Uart: uart,
	}

	return d.Loop(context.Background())
}
