//go:build tinygo

package keyboard

import (
	"machine"
)

type SquaredMatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	Pins []machine.Pin
}

func (d *Device) AddSquaredMatrixKeyboard(pins []machine.Pin, keys [][]Keycode) *SquaredMatrixKeyboard {
	state := make([]State, len(pins)*(len(pins)-1))

	for i := range pins {
		pins[i].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	k := &SquaredMatrixKeyboard{
		Pins:     pins,
		State:    state,
		Keys:     keys,
		callback: func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *SquaredMatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *SquaredMatrixKeyboard) Get() []State {
	c := int(0)
	cols := []int{}
	for i := range d.Pins {
		for j := range d.Pins {
			d.Pins[j].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
			if i == j {
				c = j
			} else {
				cols = append(cols, j)
			}
		}

		for r, j := range cols {
			d.Pins[j].Configure(machine.PinConfig{Mode: machine.PinOutput})
			d.Pins[j].Low()
			current := !d.Pins[c].Get()
			idx := r*(len(cols)+1) + c

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
			d.Pins[j].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
		cols = cols[:0]
	}

	return d.State
}

func (d *SquaredMatrixKeyboard) Key(layer, index int) Keycode {
	if layer >= len(d.Keys) {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *SquaredMatrixKeyboard) SetKeycode(layer, index int, key Keycode) {
	d.Keys[layer][index] = key
}

func (d *SquaredMatrixKeyboard) Init() error {
	return nil
}
