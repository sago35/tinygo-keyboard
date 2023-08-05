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
		callback: func(layer, row, col int, state State) {},
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
			switch d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] {
			case None:
				if current {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = Press
					d.callback(0, r, 2*len(d.Col)-1-c, Press)
				} else {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = PressToRelease
					d.callback(0, r, 2*len(d.Col)-1-c, Press)
					d.callback(0, r, 2*len(d.Col)-1-c, PressToRelease)
				}
			case Press:
				if current {
				} else {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = PressToRelease
					d.callback(0, r, 2*len(d.Col)-1-c, PressToRelease)
				}
			case PressToRelease:
				if current {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = NoneToPress
					d.callback(0, r, 2*len(d.Col)-1-c, Press)
				} else {
					d.State[r*2*len(d.Col)+2*len(d.Col)-1-c] = None
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
			switch d.State[r*2*len(d.Col)+c] {
			case None:
				if current {
					d.State[r*2*len(d.Col)+c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[r*2*len(d.Col)+c] = Press
					d.callback(0, r, c, Press)
				} else {
					d.State[r*2*len(d.Col)+c] = PressToRelease
					d.callback(0, r, c, Press)
					d.callback(0, r, c, PressToRelease)
				}
			case Press:
				if current {
				} else {
					d.State[r*2*len(d.Col)+c] = PressToRelease
					d.callback(0, r, c, PressToRelease)
				}
			case PressToRelease:
				if current {
					d.State[r*2*len(d.Col)+c] = NoneToPress
					d.callback(0, r, c, Press)
				} else {
					d.State[r*2*len(d.Col)+c] = None
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

func (d *DuplexMatrixKeyboard) Init() error {
	return nil
}
