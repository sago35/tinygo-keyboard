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

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	k := &UartKeyboard{
		State:    state,
		Keys:     keydef,
		callback: func(layer, index int, state State) {},
		uart:     uart,
		buf:      make([]byte, 0, 3),
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *UartKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *UartKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
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
			index := (int(d.buf[1]) << 8) + int(d.buf[2])
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
				} else {
				}
			case NoneToPress:
				if current {
					d.State[index] = Press
				} else {
					d.State[index] = PressToRelease
				}
			case Press:
				if current {
				} else {
					d.State[index] = PressToRelease
				}
			case PressToRelease:
				if current {
					d.State[index] = NoneToPress
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
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *UartKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *UartKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *UartKeyboard) Init() error {
	for d.uart.Buffered() > 0 {
		d.uart.ReadByte()
	}
	return nil
}
