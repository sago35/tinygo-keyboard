//go:build tinygo

package keyboard

import (
	"machine"
)

type DuplexMatrixKeyboard struct {
	State [][]State
	Keys  [][][]Keycode

	Col []machine.Pin
	Row []machine.Pin
}

func (d *Device) AddDuplexMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][][]Keycode) {
	state := [][]State{}
	col := len(colPins)
	row := len(rowPins)

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	k := &DuplexMatrixKeyboard{
		Col:   colPins,
		Row:   rowPins,
		State: state,
		Keys:  keys,
	}

	d.kb = append(d.kb, k)
}

func (d *DuplexMatrixKeyboard) Get() [][]State {
	for c := range d.Col {
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Col[c].High()
		for r := range d.Row {
			//d.State[2*len(d.Row)-r-1][c] = d.Row[r].Get()
			current := d.Row[r].Get()
			switch d.State[2*len(d.Row)-r-1][c] {
			case None:
				if current {
					d.State[2*len(d.Row)-r-1][c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[2*len(d.Row)-r-1][c] = Press
				} else {
					d.State[2*len(d.Row)-r-1][c] = PressToRelease
				}
			case Press:
				if current {
				} else {
					d.State[2*len(d.Row)-r-1][c] = PressToRelease
				}
			case PressToRelease:
				if current {
					d.State[2*len(d.Row)-r-1][c] = NoneToPress
				} else {
					d.State[2*len(d.Row)-r-1][c] = None
				}
			}
		}
		d.Col[c].Low()
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	for r := range d.Row {
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Row[r].High()
		for c := range d.Col {
			//d.State[r][c] = d.Col[c].Get()
			current := d.Col[c].Get()
			switch d.State[r][c] {
			case None:
				if current {
					d.State[r][c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[r][c] = Press
				} else {
					d.State[r][c] = PressToRelease
				}
			case Press:
				if current {
				} else {
					d.State[r][c] = PressToRelease
				}
			case PressToRelease:
				if current {
					d.State[r][c] = NoneToPress
				} else {
					d.State[r][c] = None
				}
			}
		}
		d.Row[r].Low()
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return d.State
}

func (d *DuplexMatrixKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *DuplexMatrixKeyboard) Init() error {
	return nil
}
