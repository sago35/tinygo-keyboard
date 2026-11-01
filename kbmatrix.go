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

// driveRows reports whether the rows are the driven axis. An active low scan
// drives the cathode side of the diode, so the axes swap.
func (d *MatrixKeyboard) driveRows() bool {
	return d.options.InvertDiode != d.options.InvertButtonState
}

func (d *MatrixKeyboard) inputMode() machine.PinMode {
	if d.options.InvertButtonState {
		return machine.PinInputPullup
	}
	return machine.PinInputPulldown
}

// selectLine drives one line to the active level.
func (d *MatrixKeyboard) selectLine(p machine.Pin) {
	p.Configure(machine.PinConfig{Mode: machine.PinOutput})
	p.Set(!d.options.InvertButtonState)
}

// unselectLine drives the line back to the idle level, then leaves it as an input.
func (d *MatrixKeyboard) unselectLine(p machine.Pin) {
	p.Set(d.options.InvertButtonState)
	p.Configure(machine.PinConfig{Mode: d.inputMode()})
}

// discharge drives the sense lines to the idle level for a moment, then returns
// them to input. The lines start from the idle level before the next select.
func (d *MatrixKeyboard) discharge(pins []machine.Pin) {
	for _, p := range pins {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.Set(d.options.InvertButtonState)
	}
	mode := d.inputMode()
	for _, p := range pins {
		p.Configure(machine.PinConfig{Mode: mode})
	}
}

func (d *MatrixKeyboard) read(p machine.Pin) bool {
	return p.Get() != d.options.InvertButtonState
}

func (d *MatrixKeyboard) Get() []State {
	ncol := len(d.Col)
	if d.driveRows() {
		for r := range d.Row {
			d.selectLine(d.Row[r])
			for c := range d.Col {
				d.step(r*ncol+c, d.read(d.Col[c]))
			}
			d.unselectLine(d.Row[r])
			d.discharge(d.Col)
		}
		return d.State
	}
	for c := range d.Col {
		d.selectLine(d.Col[c])
		for r := range d.Row {
			d.step(r*ncol+c, d.read(d.Row[r]))
		}
		d.unselectLine(d.Col[c])
		d.discharge(d.Row)
	}
	return d.State
}

func (d *MatrixKeyboard) step(idx int, current bool) {
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
