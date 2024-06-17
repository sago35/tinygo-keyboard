//go:build tinygo

package keyboard

import (
	"tinygo.org/x/drivers/mcp23017"
)

type ExpanderKeyboard struct {
	State    []State
	Keys     [][]Keycode
	options  Options
	callback Callback

	Col          []mcp23017.Pin
	Row          []mcp23017.Pin
	cycleCounter []uint8
	debounce     uint8
	expander     *mcp23017.Device
}

func (d *Device) AddExpanderKeyboard(expanderDevice *mcp23017.Device, colPins, rowPins []mcp23017.Pin, keys [][]Keycode, opt ...Option) *ExpanderKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)
	cycleCnt := make([]uint8, len(state))

	for c := range colPins {
		colPins[c].SetMode(mcp23017.Output)
	}
	for r := range rowPins {
		rowPins[r].SetMode(mcp23017.Output)
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

	k := &ExpanderKeyboard{
		Col:          colPins,
		Row:          rowPins,
		State:        state,
		Keys:         keydef,
		options:      o,
		callback:     func(layer, index int, state State) {},
		cycleCounter: cycleCnt,
		debounce:     8,
		expander:     expanderDevice,
	}
	d.kb = append(d.kb, k)
	return k
}

func (d *ExpanderKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *ExpanderKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
}

func (d *ExpanderKeyboard) Get() []State {
	for c := range d.Col {
		for r := range d.Row {
			current := false
			if !d.options.InvertDiode {
				d.Col[c].SetMode(mcp23017.Output)
				d.Col[c].High()
				current, _ = d.Row[r].Get()
			} else {
				d.Row[r].SetMode(mcp23017.Output)
				d.Row[r].High()
				current, _ = d.Col[c].Get()
			}
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
			if !d.options.InvertDiode {
				d.Col[c].Low()
			} else {
				d.Row[r].Low()
			}
		}
	}

	return d.State
}

func (d *ExpanderKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *ExpanderKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *ExpanderKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *ExpanderKeyboard) Init() error {
	return nil
}
