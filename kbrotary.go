package keyboard

import (
	"fmt"
	"machine"

	rotary_encoder "github.com/bgould/tinygo-rotary-encoder"
)

type RotaryKeyboard struct {
	State [][]State
	Keys  [][][]Keycode

	enc      *rotary_encoder.Device
	oldValue int
}

func (d *Device) AddRotaryKeyboard(rotA, rotB machine.Pin, keys [][][]Keycode) {
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
}

func (d *RotaryKeyboard) Get() [][]State {
	rot := []bool{false, false}
	if newValue := d.enc.Value(); newValue != d.oldValue {
		if newValue < d.oldValue {
			rot[0] = true
		} else {
			rot[1] = true
		}
		fmt.Printf("%#v\n", rot)
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
			} else {
				d.State[0][c] = PressToRelease
			}
		case Press:
			if current {
			} else {
				d.State[0][c] = PressToRelease
			}
		case PressToRelease:
			if current {
				d.State[0][c] = NoneToPress
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
