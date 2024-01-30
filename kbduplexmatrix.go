//go:build tinygo

package keyboard

import (
	"machine"
)

type DuplexMatrixKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	Col          []machine.Pin
	Row          []machine.Pin
	cycleCounter []uint8
	debounce     uint8
}

func (d *Device) AddDuplexMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode) *DuplexMatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*2*col)
	cycleCnt := make([]uint8, len(state))

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
		Col:          colPins,
		Row:          rowPins,
		cycleCounter: cycleCnt,
		State:        state,
		Keys:         keydef,
		callback:     func(layer, index int, state State) {},
		debounce:     8,
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *DuplexMatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *DuplexMatrixKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
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
					if d.cycleCounter[idx] >= d.debounce {
						d.State[idx] = NoneToPress
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
					}
				} else {
					d.cycleCounter[idx] = 0
				}
			case NoneToPress:
				d.State[idx] = Press
			case Press:
				if current {
					d.cycleCounter[idx] = 0
				} else {
					if d.cycleCounter[idx] >= d.debounce {
						d.State[idx] = PressToRelease
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
					}
				}
			case PressToRelease:
				d.State[idx] = None
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
					if d.cycleCounter[idx] >= d.debounce {
						d.State[idx] = NoneToPress
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
					}
				} else {
					d.cycleCounter[idx] = 0
				}
			case NoneToPress:
				d.State[idx] = Press
			case Press:
				if current {
					d.cycleCounter[idx] = 0
				} else {
					if d.cycleCounter[idx] >= d.debounce {
						d.State[idx] = PressToRelease
						d.cycleCounter[idx] = 0
					} else {
						d.cycleCounter[idx]++
					}
				}
			case PressToRelease:
				d.State[idx] = None
			}
		}
		d.Row[r].High()
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
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
