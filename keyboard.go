package keyboard

import (
	"context"
	"fmt"
	"machine"
	k "machine/usb/hid/keyboard"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

type Device struct {
	Col      []machine.Pin
	Row      []machine.Pin
	State    [][]State
	Keys     [][][]Keycode
	Keyboard UpDowner
	Debug    bool
}

type UpDowner interface {
	Up(c k.Keycode) error
	Down(c k.Keycode) error
}

type State uint8

const (
	None State = iota
	NoneToPress
	Press
	PressToRelease
)

func New(colPins, rowPins []machine.Pin, keys [][][]Keycode) *Device {
	state := [][]State{}
	col := len(colPins)
	row := len(rowPins)

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	for c := range colPins {
		colPins[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	for r := range rowPins {
		rowPins[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	d := &Device{
		Col:      colPins,
		Row:      rowPins,
		State:    state,
		Keys:     keys,
		Keyboard: k.Port(),
	}

	return d
}

func (d *Device) Loop(ctx context.Context) error {
	layer := 0
	cont := true
	for cont {
		select {
		case <-ctx.Done():
			cont = false
			continue
		default:
		}

		d.Get()

		for row := range d.State {
			for col := range d.State[row] {
				switch d.State[row][col] {
				case None:
					// skip
				case NoneToPress:
					if d.Keys[layer][row][col]&keycodes.ModKeyMask == keycodes.ModKeyMask {
						layer = int(d.Keys[layer][row][col]) & 0x0F
					} else {
						d.Keyboard.Down(k.Keycode(d.Keys[layer][row][col]))
					}
					if d.Debug {
						fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[0][row][col])
					}
				case Press:
				case PressToRelease:
					if d.Keys[layer][row][col]&keycodes.ModKeyMask == keycodes.ModKeyMask {
						layer = 0
					} else {
						d.Keyboard.Up(k.Keycode(d.Keys[layer][row][col]))
					}
					if d.Debug {
						fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[0][row][col])
					}
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (d *Device) Get() [][]State {
	wait := 1 * time.Millisecond

	for c := range d.Col {
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Col[c].High()
		time.Sleep(wait)
		for r := range d.Row {
			//d.State[2*len(d.Row)-r-1][c] = d.Row[r].Get()
			current := d.Row[r].Get()
			switch d.State[2*len(d.Row)-r-1][c] {
			case None:
				if current {
					d.State[2*len(d.Row)-r-1][c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[2*len(d.Row)-r-1][c] = Press
				} else {
					d.State[2*len(d.Row)-r-1][c] = PressToRelease
				}
			case Press:
				if current {
				} else {
					d.State[2*len(d.Row)-r-1][c] = PressToRelease
				}
			case PressToRelease:
				if current {
					d.State[2*len(d.Row)-r-1][c] = NoneToPress
				} else {
					d.State[2*len(d.Row)-r-1][c] = None
				}
			}
		}
		d.Col[c].Low()
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	for r := range d.Row {
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Row[r].High()
		time.Sleep(wait)
		for c := range d.Col {
			//d.State[r][c] = d.Col[c].Get()
			current := d.Col[c].Get()
			switch d.State[r][c] {
			case None:
				if current {
					d.State[r][c] = NoneToPress
				} else {
				}
			case NoneToPress:
				if current {
					d.State[r][c] = Press
				} else {
					d.State[r][c] = PressToRelease
				}
			case Press:
				if current {
				} else {
					d.State[r][c] = PressToRelease
				}
			case PressToRelease:
				if current {
					d.State[r][c] = NoneToPress
				} else {
					d.State[r][c] = None
				}
			}
		}
		d.Row[r].Low()
		d.Row[r].Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return d.State
}

type Keycode k.Keycode
