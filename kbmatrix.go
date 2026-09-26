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
	debounce     uint8
}

func (d *Device) AddMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][]Keycode, opt ...Option) *MatrixKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)
	cycleCnt := make([]uint8, len(state))

	o := Options{}
	for _, f := range opt {
		f(&o)
	}

	inputMode := machine.PinInputPulldown
	if o.InvertButtonState {
		inputMode = machine.PinInputPullup
	}
	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: inputMode})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: inputMode})
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
		debounce:     8,
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *MatrixKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *MatrixKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
}

// scanPins returns the pin to drive and the pin to read for one key.
// An active low scan drives the cathode side of the diode, so the axes swap.
func (d *MatrixKeyboard) scanPins(c, r int) (drive, read machine.Pin) {
	if d.options.InvertDiode != d.options.InvertButtonState {
		return d.Row[r], d.Col[c]
	}
	return d.Col[c], d.Row[r]
}

func (d *MatrixKeyboard) Get() []State {
	activeLow := d.options.InvertButtonState
	inputMode := machine.PinInputPulldown
	if activeLow {
		inputMode = machine.PinInputPullup
	}
	for c := range d.Col {
		for r := range d.Row {
			drive, read := d.scanPins(c, r)
			drive.Configure(machine.PinConfig{Mode: machine.PinOutput})
			drive.Set(!activeLow)
			current := read.Get() != activeLow
			idx := r*len(d.Col) + c
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
			drive.Set(activeLow)
			drive.Configure(machine.PinConfig{Mode: inputMode})
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
