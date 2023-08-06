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

	k := &GpioKeyboard{
		Col:      pins,
		State:    state,
		Keys:     keys,
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
	if layer >= len(d.Keys) {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *GpioKeyboard) SetKeycode(layer, index int, key Keycode) {
	d.Keys[layer][index] = key
}

func (d *GpioKeyboard) Init() error {
	return nil
}
