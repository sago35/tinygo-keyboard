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

	Col          []int
	Row          []int
	cycleCounter []uint8
	debounce     uint8
	expander     *mcp23017.Device
}

func (d *Device) AddExpanderKeyboard(expanderDevice *mcp23017.Device, colPins, rowPins []int, keys [][]Keycode, opt ...Option) *ExpanderKeyboard {
	col := len(colPins)
	row := len(rowPins)
	state := make([]State, row*col)
	cycleCnt := make([]uint8, len(state))

	o := Options{}
	for _, f := range opt {
		f(&o)
	}

	for _, c := range colPins {

		if !o.InvertDiode {
			_ = expanderDevice.Pin(c).SetMode(mcp23017.Output)
			_ = expanderDevice.Pin(c).Set(true)
		} else {
			_ = expanderDevice.Pin(c).SetMode(mcp23017.Input | mcp23017.Pullup)
		}

	}
	for _, r := range rowPins {
		if !o.InvertDiode {
			_ = expanderDevice.Pin(r).SetMode(mcp23017.Input | mcp23017.Pullup)
		} else {
			_ = expanderDevice.Pin(r).SetMode(mcp23017.Output)
			_ = expanderDevice.Pin(r).Set(true)
		}
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
		debounce:     2,
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
	for cIdx, c := range d.Col {
		for rIdx, r := range d.Row {
			current := false
			if !d.options.InvertDiode {
				_ = d.expander.Pin(c).Set(false)
				current, _ = d.expander.Pin(r).Get()
			} else {
				_ = d.expander.Pin(r).Set(false)
				current, _ = d.expander.Pin(c).Get()
			}
			idx := rIdx*len(d.Col) + cIdx
			switch d.State[idx] {
			case None:
				if !current {
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
				if !current {
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
				_ = d.expander.Pin(c).Set(true)
			} else {
				_ = d.expander.Pin(r).Set(true)
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
