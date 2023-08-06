//go:build tinygo

package keyboard

import (
	"machine"

	rotary_encoder "github.com/bgould/tinygo-rotary-encoder"
)

type RotaryKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	enc      *rotary_encoder.Device
	oldValue int
}

func (d *Device) AddRotaryKeyboard(rotA, rotB machine.Pin, keys [][]Keycode) *RotaryKeyboard {
	state := make([]State, 2)

	enc := rotary_encoder.New(rotA, rotB)
	enc.Configure()

	k := &RotaryKeyboard{
		State:    state,
		Keys:     keys,
		callback: func(layer, index int, state State) {},

		enc: enc,
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *RotaryKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *RotaryKeyboard) Get() []State {
	rot := []bool{false, false}
	if newValue := d.enc.Value(); newValue != d.oldValue {
		if newValue < d.oldValue {
			rot[0] = true
		} else {
			rot[1] = true
		}
		d.oldValue = newValue
	}

	for c, current := range rot {
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

func (d *RotaryKeyboard) Key(layer, index int) Keycode {
	if layer >= len(d.Keys) {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *RotaryKeyboard) SetKeycode(layer, index int, key Keycode) {
	d.Keys[layer][index] = key
}

func (d *RotaryKeyboard) Init() error {
	return nil
}
