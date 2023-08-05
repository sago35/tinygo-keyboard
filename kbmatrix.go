//go:build tinygo

package keyboard

import (
	"machine"
)

type MatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	options  Options
	callback Callback

	Col []machine.Pin
	Row []machine.Pin
}

func (d *Device) AddMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode, opt ...Option) *MatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	o := Options{}
	for _, f := range opt {
		f(&o)
	}

	k := &MatrixKeyboard{
		Col:      colPins,
		Row:      rowPins,
		State:    state,
		Keys:     keys,
		options:  o,
		callback: func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *MatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *MatrixKeyboard) Get() []State {
	for c := range d.Col {
		for r := range d.Row {
			//d.State[r*len(d.Col)+c] = d.Row[r].Get()
			current := false
			if !d.options.InvertDiode {
				d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
				d.Col[c].High()
				current = d.Row[r].Get()
			} else {
				d.Row[r].Configure(machine.PinConfig{Mode: machine.PinOutput})
				d.Row[r].High()
				current = d.Col[c].Get()
			}
			idx := r*len(d.Col) + c
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
			if !d.options.InvertDiode {
				d.Col[c].Low()
				d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
			} else {
				d.Row[r].Low()
				d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
			}
		}
	}

	return d.State
}

func (d *MatrixKeyboard) Key(layer, index int) Keycode {
	return d.Keys[layer][index]
}

func (d *MatrixKeyboard) Init() error {
	return nil
}
