package keyboard

import (
	"context"
	"fmt"
	"machine"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

type Device struct {
	Col      []machine.Pin
	Row      []machine.Pin
	State    [][]State
	Keys     [][][]Keycode
	State2   [][]State     // for UART
	Keys2    [][][]Keycode // for UART
	Keyboard UpDowner
	Mouse    Mouser
	Override [][]Keycode
	Debug    bool

	modKeyCallback func(layer int, down bool)
	layer          int
	pressed        []Keycode
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

	kb := &Keyboard{
		Port: k.Port(),
	}
	d := &Device{
		Col:      colPins,
		Row:      rowPins,
		State:    state,
		Keys:     keys,
		Keyboard: kb,
		Mouse:    mouse.Port(),
		pressed:  make([]Keycode, 0, 10),
	}

	return d
}

func (d *Device) OverrideCtrlH() {
	d.Keyboard = &Keyboard{
		Port:          k.Port(),
		overrideCtrlH: true,
	}
}

func (d *Device) AddUartKeyboard(row, col int, keys [][][]Keycode) {
	state := [][]State{}

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	d.State2 = state
	d.Keys2 = keys
}

func (d *Device) Callback(fn func(layer int, down bool)) {
	d.modKeyCallback = fn
}

func (d *Device) Mod(layer int, down bool) {
	if down {
		d.layer = layer
	} else {
		d.layer = 0
		for _, p := range d.pressed {
			if p&0xF000 == 0xD000 {
				switch p & 0x00FF {
				case 0x01, 0x02, 0x03:
					d.Mouse.Release(mouse.Button(p & 0x00FF))
				case 0x04:
					//d.Mouse.WheelDown()
				case 0x05:
					//d.Mouse.WheelUp()
				}
			} else {
				d.Keyboard.Up(k.Keycode(p))
			}
		}
		d.pressed = d.pressed[:]
	}
}

func (d *Device) Loop(ctx context.Context) error {
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
					x := d.Keys[d.layer][row][col]
					found := false
					for _, p := range d.pressed {
						if x == p {
							found = true
						}
					}
					if !found {
						d.pressed = append(d.pressed, x)
					}

					if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
						d.layer = int(x) & 0x0F
						if d.modKeyCallback != nil {
							d.modKeyCallback(d.layer, true)
						}
					} else if x&0xF000 == 0xD000 {
						switch x & 0x00FF {
						case 0x01, 0x02, 0x03:
							d.Mouse.Press(mouse.Button(x & 0x00FF))
						case 0x04:
							d.Mouse.WheelDown()
						case 0x05:
							d.Mouse.WheelUp()
						}
					} else {
						d.Keyboard.Down(k.Keycode(x))
					}
					if d.Debug {
						fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[0][row][col])
					}
				case Press:
				case PressToRelease:
					x := d.Keys[d.layer][row][col]

					for i, p := range d.pressed {
						if x == p {
							d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
						}
					}

					if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
						if d.modKeyCallback != nil {
							d.modKeyCallback(d.layer, false)
						}
						d.layer = 0

						for _, p := range d.pressed {
							if p&0xF000 == 0xD000 {
								switch p & 0x00FF {
								case 0x01, 0x02, 0x03:
									d.Mouse.Release(mouse.Button(p & 0x00FF))
								case 0x04:
									//d.Mouse.WheelDown()
								case 0x05:
									//d.Mouse.WheelUp()
								}
							} else {
								d.Keyboard.Up(k.Keycode(p))
							}
						}
						d.pressed = d.pressed[:]

					} else if x&0xF000 == 0xD000 {
						switch x & 0x00FF {
						case 0x01, 0x02, 0x03:
							d.Mouse.Release(mouse.Button(x & 0x00FF))
						case 0x04:
							//d.Mouse.WheelDown()
						case 0x05:
							//d.Mouse.WheelUp()
						}
					} else {
						d.Keyboard.Up(k.Keycode(x))
					}
					if d.Debug {
						fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[0][row][col])
					}
				}
			}
		}

		time.Sleep(5 * time.Millisecond)
	}

	return nil
}

func (d *Device) Get() [][]State {
	for c := range d.Col {
		d.Col[c].Configure(machine.PinConfig{Mode: machine.PinOutput})
		d.Col[c].High()
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

type Keyboard struct {
	pressed       []k.Keycode
	override      []k.Keycode
	Port          UpDowner
	overrideCtrlH bool
}

func (k *Keyboard) Up(c k.Keycode) error {
	if len(k.override) > 0 {
		for _, p := range k.override {
			k.Port.Up(p)
		}
		k.override = k.override[:0]
		for _, p := range k.pressed {
			// When overriding, do not press the last key again
			if c != p && p != k.pressed[len(k.pressed)-1] {
				k.Port.Down(p)
			}
		}
	}

	for i, p := range k.pressed {
		if c == p {
			k.pressed = append(k.pressed[:i], k.pressed[i+1:]...)
			return k.Port.Up(c)
		}
	}
	return nil
}

func (k *Keyboard) Down(c k.Keycode) error {
	found := false
	for _, p := range k.pressed {
		if c == p {
			found = true
		}
	}
	if !found {
		k.pressed = append(k.pressed, c)

		if k.overrideCtrlH && len(k.pressed) == 2 && k.pressed[0] == keycodes.KeyLeftCtrl && k.pressed[1] == keycodes.KeyH {
			for _, p := range k.pressed {
				k.Port.Up(p)
			}
			k.override = append(k.override, keycodes.KeyBackspace)
			return k.Port.Down(keycodes.KeyBackspace)
		} else {
			if len(k.override) > 0 {
				for _, p := range k.override {
					k.Port.Up(p)
				}
				k.override = k.override[:0]
				for _, p := range k.pressed {
					k.Port.Down(p)
				}
			}
			return k.Port.Down(c)
		}
	}
	return nil
}
