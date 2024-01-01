//go:build tinygo

package keyboard

import (
	"tinygo.org/x/drivers/shifter"
)

type ShifterKeyboard struct {
	State    []State
	Keys     [][]Keycode
	options  Options
	callback Callback

	Shifter shifter.Device
}

func (d *Device) AddShifterKeyboard(shifterDevice shifter.Device, keys [][]Keycode, opt ...Option) *ShifterKeyboard {
	state := make([]State, len(shifterDevice.Pins))

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

	k := &ShifterKeyboard{
		Shifter:  shifterDevice,
		State:    state,
		Keys:     keydef,
		options:  o,
		callback: func(layer, index int, state State) {},
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *ShifterKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *ShifterKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
}

func (d *ShifterKeyboard) Get() []State {
	d.Shifter.Read8Input()

	for c := 0; c < len(d.Shifter.Pins); c++ {
		current := d.Shifter.Pins[c].Get()

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
			} else {
				d.State[c] = PressToRelease
			}
		case Press:
			if current {
			} else {
				d.State[c] = PressToRelease
			}
		case PressToRelease:
			if current {
				d.State[c] = NoneToPress
			} else {
				d.State[c] = None
			}
		}
	}

	return d.State
}

func (d *ShifterKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *ShifterKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *ShifterKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *ShifterKeyboard) Init() error {
	return nil
}
