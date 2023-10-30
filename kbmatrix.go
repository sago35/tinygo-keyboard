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

	Col          []machine.Pin
	Row          []machine.Pin
	cycleCounter []uint8
}

const matrixCyclesToPreventChattering = uint8(4)

func (d *Device) AddMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode, opt ...Option) *MatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)
	cycleCnt := make([]uint8, len(state))

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

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	k := &MatrixKeyboard{
		Col:          colPins,
		Row:          rowPins,
		State:        state,
		Keys:         keydef,
		options:      o,
		callback:     func(layer, index int, state State) {},
		cycleCounter: cycleCnt,
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
					if d.cycleCounter[idx] >= matrixCyclesToPreventChattering {
						d.State[idx] = NoneToPress
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
					}
				} else {
					d.cycleCounter[idx] = 0
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
					d.cycleCounter[idx] = 0
				} else {
					if d.cycleCounter[idx] >= matrixCyclesToPreventChattering {
						d.State[idx] = PressToRelease
						d.callback(0, idx, PressToRelease)
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
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
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *MatrixKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *MatrixKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *MatrixKeyboard) Init() error {
	return nil
}
