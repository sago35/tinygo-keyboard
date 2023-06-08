//go:build tinygo

package keyboard

import (
	"machine"
)

type UartKeyboard struct {
	State    [][]State
	Keys     [][][]Keycode
	callback Callback

	uart *machine.UART
	buf  []byte
}

func (d *Device) AddUartKeyboard(row, col int, uart *machine.UART, keys [][][]Keycode) *UartKeyboard {
	state := [][]State{}

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	k := &UartKeyboard{
		State: state,
		Keys:  keys,
		uart:  uart,
		buf:   make([]byte, 0, 3),
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *UartKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *UartKeyboard) Get() [][]State {
	uart := d.uart

	for row := range d.State {
		for col := range d.State[row] {
			switch d.State[row][col] {
			case NoneToPress:
				d.State[row][col] = Press
			case PressToRelease:
				d.State[row][col] = None
			}
		}
	}

	for uart.Buffered() > 0 {
		data, _ := uart.ReadByte()
		d.buf = append(d.buf, data)

		if len(d.buf) == 3 {
			row, col := int(d.buf[1]), int(d.buf[2])
			current := false
			switch d.buf[0] {
			case 0xAA: // press
				current = true
			case 0x55: // release
				current = false
			default:
				d.buf[0], d.buf[1] = d.buf[1], d.buf[2]
				d.buf = d.buf[:2]
				continue
			}

			switch d.State[row][col] {
			case None:
				if current {
					d.State[row][col] = NoneToPress
					d.callback(0, row, col, Press)
				} else {
				}
			case NoneToPress:
				if current {
					d.State[row][col] = Press
				} else {
					d.State[row][col] = PressToRelease
					d.callback(0, row, col, PressToRelease)
				}
			case Press:
				if current {
				} else {
					d.State[row][col] = PressToRelease
					d.callback(0, row, col, PressToRelease)
				}
			case PressToRelease:
				if current {
					d.State[row][col] = NoneToPress
					d.callback(0, row, col, Press)
				} else {
					d.State[row][col] = None
				}
			}
			d.buf = d.buf[:0]
		}
	}
	return d.State
}

func (d *UartKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *UartKeyboard) Init() error {
	for d.uart.Buffered() > 0 {
		d.uart.ReadByte()
	}
	return nil
}
