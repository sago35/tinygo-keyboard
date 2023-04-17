//go:build tinygo

package keyboard

import (
	"machine"
)

type GpioKeyboard struct {
	State   [][]State
	Keys    [][][]Keycode
	options Options

	Col []machine.Pin
}

func (d *Device) AddGpioKeyboard(pins []machine.Pin, keys [][][]Keycode, opt ...Option) {
	state := [][]State{}
	col := len(pins)

	column := make([]State, col)
	state = append(state, column)

	o := Options{
		InvertButtonState: true,
	}
	for _, f := range opt {
		f(&o)
	}

	k := &GpioKeyboard{
		Col:     pins,
		State:   state,
		Keys:    keys,
		options: o,
	}

	d.kb = append(d.kb, k)
}

func (d *GpioKeyboard) Get() [][]State {
	for c := range d.Col {
		current := d.Col[c].Get()

		if d.options.InvertButtonState {
			current = !current
		}

		switch d.State[0][c] {
		case None:
			if current {
				d.State[0][c] = NoneToPress
			} else {
			}
		case NoneToPress:
			if current {
				d.State[0][c] = Press
			} else {
				d.State[0][c] = PressToRelease
			}
		case Press:
			if current {
			} else {
				d.State[0][c] = PressToRelease
			}
		case PressToRelease:
			if current {
				d.State[0][c] = NoneToPress
			} else {
				d.State[0][c] = None
			}
		}
	}

	return d.State
}

func (d *GpioKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *GpioKeyboard) Init() error {
	return nil
}
