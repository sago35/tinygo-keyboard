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

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	k := &RotaryKeyboard{
		State:    state,
		Keys:     keydef,
		callback: func(layer, index int, state State) {},

		enc: enc,
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *RotaryKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *RotaryKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
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

func (d *RotaryKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *RotaryKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *RotaryKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *RotaryKeyboard) Init() error {
	return nil
}
