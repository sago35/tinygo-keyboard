//go:build tinygo

package keyboard

import (
	"machine"
)

type DuplexMatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	Col []machine.Pin
	Row []machine.Pin
}

func (d *Device) AddDuplexMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode) *DuplexMatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*2*col)

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		colPins[c].High()
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	k := &DuplexMatrixKeyboard{
		Col:      colPins,
		Row:      rowPins,
		State:    state,
		Keys:     keys,
		callback: func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *DuplexMatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *DuplexMatrixKeyboard) Get() []State {
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
			case Press:
				if current {
				} else {
					d.State[idx] = PressToRelease
					d.callback(0, idx, PressToRelease)
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
			case Press:
				if current {
				} else {
					d.State[idx] = PressToRelease
					d.callback(0, idx, PressToRelease)
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
	}

	return d.State
}

func (d *DuplexMatrixKeyboard) Key(layer, index int) Keycode {
	return d.Keys[layer][index]
}

func (d *DuplexMatrixKeyboard) SetKeycode(layer, index int, key Keycode) {
	d.Keys[layer][index] = key
}

func (d *DuplexMatrixKeyboard) Init() error {
	return nil
}
