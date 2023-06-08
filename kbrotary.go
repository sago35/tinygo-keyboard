//go:build tinygo && macropad_rp2040

package keyboard

import (
	"machine"

	rotary_encoder "github.com/bgould/tinygo-rotary-encoder"
)

type RotaryKeyboard struct {
	State    [][]State
	Keys     [][][]Keycode
	callback Callback

	enc      *rotary_encoder.Device
	oldValue int
}

func (d *Device) AddRotaryKeyboard(rotA, rotB machine.Pin, keys [][][]Keycode) *RotaryKeyboard {
	state := [][]State{}

	column := make([]State, 2)
	state = append(state, column)

	enc := rotary_encoder.New(rotA, rotB)
	enc.Configure()

	k := &RotaryKeyboard{
		State: state,
		Keys:  keys,

		enc: enc,
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *RotaryKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *RotaryKeyboard) Get() [][]State {
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
		switch d.State[0][c] {
		case None:
			if current {
				d.State[0][c] = NoneToPress
			} else {
			}
		case NoneToPress:
			if current {
				d.State[0][c] = Press
				d.callback(0, 0, c, Press)
			} else {
				d.State[0][c] = PressToRelease
				d.callback(0, 0, c, Press)
				d.callback(0, 0, c, PressToRelease)
			}
		case Press:
			if current {
			} else {
				d.State[0][c] = PressToRelease
				d.callback(0, 0, c, PressToRelease)
			}
		case PressToRelease:
			if current {
				d.State[0][c] = NoneToPress
				d.callback(0, 0, c, Press)
			} else {
				d.State[0][c] = None
			}
		}
	}

	return d.State
}

func (d *RotaryKeyboard) Key(layer, row, col int) Keycode {
	return d.Keys[layer][row][col]
}

func (d *RotaryKeyboard) Init() error {
	return nil
}
