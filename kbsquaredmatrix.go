//go:build tinygo

package keyboard

import (
	"machine"
)

type SquaredMatrixKeyboard struct {
	State    [][]State
	Keys     [][][]Keycode
	callback Callback

	Pins []machine.Pin
}

func (d *Device) AddSquaredMatrixKeyboard(pins []machine.Pin, keys [][][]Keycode) *SquaredMatrixKeyboard {
	state := [][]State{}
	col := len(pins)
	row := len(pins) - 1

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	for i := range pins {
		pins[i].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	k := &SquaredMatrixKeyboard{
		Pins:     pins,
		State:    state,
		Keys:     keys,
		callback: func(layer, row, col int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *SquaredMatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *SquaredMatrixKeyboard) Get() [][]State {
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

			switch d.State[r][c] {
			case None:
				if current {
					d.State[r][c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[r][c] = Press
					d.callback(0, r, c, Press)
				} else {
					d.State[r][c] = PressToRelease
					d.callback(0, r, c, Press)
					d.callback(0, r, c, PressToRelease)
				}
			case Press:
				if current {
				} else {
					d.State[r][c] = PressToRelease
					d.callback(0, r, c, PressToRelease)
				}
			case PressToRelease:
				if current {
					d.State[r][c] = NoneToPress
					d.callback(0, r, c, Press)
				} else {
					d.State[r][c] = None
				}
			}
			d.Pins[j].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
		cols = cols[:0]
	}

	return d.State
}

func (d *SquaredMatrixKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *SquaredMatrixKeyboard) Init() error {
	return nil
}
