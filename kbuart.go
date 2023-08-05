//go:build tinygo

package keyboard

import (
	"machine"
)

type UartKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	uart *machine.UART
	buf  []byte
}

func (d *Device) AddUartKeyboard(size int, uart *machine.UART, keys [][]Keycode) *UartKeyboard {
	state := make([]State, size)

	k := &UartKeyboard{
		State:    state,
		Keys:     keys,
		callback: func(layer, row, col int, state State) {},
		uart:     uart,
		buf:      make([]byte, 0, 3),
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *UartKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *UartKeyboard) Get() []State {
	uart := d.uart

	for i := range d.State {
		switch d.State[i] {
		case NoneToPress:
			d.State[i] = Press
		case PressToRelease:
			d.State[i] = None
		}
	}

	for uart.Buffered() > 0 {
		data, _ := uart.ReadByte()
		d.buf = append(d.buf, data)

		if len(d.buf) == 3 {
			index := int(d.buf[1]) + int(d.buf[2])*10
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

			switch d.State[index] {
			case None:
				if current {
					d.State[index] = NoneToPress
					d.callback(0, index, 0, Press)
				} else {
				}
			case NoneToPress:
				if current {
					d.State[index] = Press
				} else {
					d.State[index] = PressToRelease
					d.callback(0, index, 0, PressToRelease)
				}
			case Press:
				if current {
				} else {
					d.State[index] = PressToRelease
					d.callback(0, index, 0, PressToRelease)
				}
			case PressToRelease:
				if current {
					d.State[index] = NoneToPress
					d.callback(0, index, 0, Press)
				} else {
					d.State[index] = None
				}
			}
			d.buf = d.buf[:0]
		}
	}
	return d.State
}

func (d *UartKeyboard) Key(layer, index int) Keycode {
	return d.Keys[layer][index]
}

func (d *UartKeyboard) Init() error {
	for d.uart.Buffered() > 0 {
		d.uart.ReadByte()
	}
	return nil
}
