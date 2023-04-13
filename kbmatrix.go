package keyboard

import (
	"machine"
)

type MatrixKeyboard struct {
	State [][]State
	Keys  [][][]Keycode

	Col []machine.Pin
	Row []machine.Pin
}

func (d *Device) AddMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][][]Keycode) {
	state := [][]State{}
	col := len(colPins)
	row := len(rowPins)

	for r := 0; r < row; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	k := &MatrixKeyboard{
		Col:   colPins,
		Row:   rowPins,
		State: state,
		Keys:  keys,
	}

	d.kb = append(d.kb, k)
}

func (d *MatrixKeyboard) Get() [][]State {
	for c := range d.Col {
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Col[c].High()
		for r := range d.Row {
			//d.State[r][c] = d.Row[r].Get()
			current := d.Row[r].Get()
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
		d.Col[c].Low()
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return d.State
}

func (d *MatrixKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *MatrixKeyboard) Init() error {
	return nil
}
