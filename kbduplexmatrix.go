//go:build tinygo

package keyboard

import (
	"machine"
)

type DuplexMatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	Col        []machine.Pin
	Row        []machine.Pin
	CounterCol []uint8
	CounterRow []uint8
}

func (d *Device) AddDuplexMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode) *DuplexMatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	counterCol := make([]uint8, len(colPins))
	counterRow := make([]uint8, len(rowPins))
	state := make([]State, row*2*col)

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		colPins[c].High()
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	k := &DuplexMatrixKeyboard{
		Col:        colPins,
		Row:        rowPins,
		CounterCol: counterCol,
		CounterRow: counterRow,
		State:      state,
		Keys:       keydef,
		callback:   func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *DuplexMatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *DuplexMatrixKeyboard) Get() []State {
	//waiting count for chattering prevention
	count := uint8(4)

	for c := range d.Col {
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Col[c].Low()
		for r := range d.Row {
			current := !d.Row[r].Get()
			idx := r*2*len(d.Col) + 2*len(d.Col) - 1 - c
			switch d.State[idx] {
			case None:
				if current {
					d.State[idx] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[idx] = Press
					d.callback(0, idx, Press)
				} else {
					d.State[idx] = PressToRelease
					d.callback(0, idx, Press)
					d.callback(0, idx, PressToRelease)
				}
				d.CounterCol[c] = 0
			case Press:
				if current {
					d.CounterCol[c] = 0
				} else {
					if d.CounterCol[c] >= count {
						d.State[idx] = PressToRelease
						d.callback(0, idx, PressToRelease)
					}
				}
			case PressToRelease:
				if current {
					d.State[idx] = NoneToPress
					d.callback(0, idx, Press)
				} else {
					d.State[idx] = None
				}
			}
		}
		d.Col[c].High()
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})

		d.CounterCol[c]++
		d.CounterCol[c] = d.CounterCol[c] % (count + 1)
	}

	for r := range d.Row {
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Row[r].Low()
		for c := range d.Col {
			current := !d.Col[c].Get()
			idx := r*2*len(d.Col) + c
			switch d.State[idx] {
			case None:
				if current {
					d.State[idx] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[idx] = Press
					d.callback(0, idx, Press)
				} else {
					d.State[idx] = PressToRelease
					d.callback(0, idx, Press)
					d.callback(0, idx, PressToRelease)
				}
				d.CounterRow[r] = 0
			case Press:
				if current {
					d.CounterRow[r] = 0
				} else {
					if d.CounterRow[r] >= count {
						d.State[idx] = PressToRelease
						d.callback(0, idx, PressToRelease)
					}
				}
			case PressToRelease:
				if current {
					d.State[idx] = NoneToPress
					d.callback(0, idx, Press)
				} else {
					d.State[idx] = None
				}
			}
		}
		d.Row[r].High()
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPullup})

		d.CounterRow[r]++
		d.CounterRow[r] = d.CounterRow[r] % (count + 1)
	}

	return d.State
}

func (d *DuplexMatrixKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *DuplexMatrixKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *DuplexMatrixKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *DuplexMatrixKeyboard) Init() error {
	return nil
}
