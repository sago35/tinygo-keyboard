//go:build tinygo

package keyboard

import (
	"machine"
)

type GpioKeyboard struct {
	State    []State
	Keys     [][]Keycode
	options  Options
	callback Callback

	Col []machine.Pin
}

func (d *Device) AddGpioKeyboard(pins []machine.Pin, keys [][]Keycode, opt ...Option) *GpioKeyboard {
	col := len(pins)
	state := make([]State, col)

	o := Options{
		InvertButtonState: true,
	}
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

	k := &GpioKeyboard{
		Col:      pins,
		State:    state,
		Keys:     keydef,
		options:  o,
		callback: func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *GpioKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *GpioKeyboard) Get() []State {
	for c := range d.Col {
		current := d.Col[c].Get()

		if d.options.InvertButtonState {
			current = !current
		}

		switch d.State[c] {
		case None:
			if current {
				d.State[c] = NoneToPress
			} else {
			}
		case NoneToPress:
			if current {
				d.State[c] = Press
				d.callback(0, c, Press)
			} else {
				d.State[c] = PressToRelease
				d.callback(0, c, Press)
				d.callback(0, c, PressToRelease)
			}
		case Press:
			if current {
			} else {
				d.State[c] = PressToRelease
				d.callback(0, c, PressToRelease)
			}
		case PressToRelease:
			if current {
				d.State[c] = NoneToPress
				d.callback(0, c, Press)
			} else {
				d.State[c] = None
			}
		}
	}

	return d.State
}

func (d *GpioKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *GpioKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *GpioKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *GpioKeyboard) Init() error {
	return nil
}
