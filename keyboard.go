package keyboard

import (
	"context"
	"machine"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

type Device struct {
	Keyboard UpDowner
	Mouse    Mouser
	Override [][]Keycode
	Debug    bool

	dmk []*DuplexMatrixKeyboard
	uk  []*UartKeyboard

	modKeyCallback func(layer int, down bool)
	layer          int
	pressed        []Keycode
}

type DuplexMatrixKeyboard struct {
	State [][]State
	Keys  [][][]Keycode

	Col []machine.Pin
	Row []machine.Pin
}

type UartKeyboard struct {
	State [][]State
	Keys  [][][]Keycode

	uart *machine.UART
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

func New() *Device {
	kb := &Keyboard{
		Port: k.Port(),
	}
	d := &Device{
		Keyboard: kb,
		Mouse:    mouse.Port(),
		pressed:  make([]Keycode, 0, 10),
	}

	return d
}

func (d *Device) AddDuplexMatrixKeyboard(colPins, rowPins []machine.Pin, keys [][][]Keycode) {
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

	k := &DuplexMatrixKeyboard{
		Col:   colPins,
		Row:   rowPins,
		State: state,
		Keys:  keys,
	}

	d.dmk = append(d.dmk, k)
}

func (d *Device) OverrideCtrlH() {
	d.Keyboard = &Keyboard{
		Port:          k.Port(),
		overrideCtrlH: true,
	}
}

func (d *Device) AddUartKeyboard(row, col int, uart *machine.UART, keys [][][]Keycode) {
	state := [][]State{}

	for r := 0; r < row*2; r++ {
		column := make([]State, col)
		state = append(state, column)
	}

	u := &UartKeyboard{
		State: state,
		Keys:  keys,
		uart:  uart,
	}
	d.uk = append(d.uk, u)
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
	buf := make([]byte, 0, 3)

	for _, k := range d.uk {
		uart := k.uart
		for uart.Buffered() > 0 {
			uart.ReadByte()
		}
	}

	cont := true
	for cont {
		select {
		case <-ctx.Done():
			cont = false
			continue
		default:
		}

		pressToRelease := []Keycode{}

		// read from key matrix
		for _, k := range d.dmk {
			k.Get()
			for row := range k.State {
				for col := range k.State[row] {
					switch k.State[row][col] {
					case None:
						// skip
					case NoneToPress:
						x := k.Keys[d.layer][row][col]
						found := false
						for _, p := range d.pressed {
							if x == p {
								found = true
							}
						}
						if !found {
							d.pressed = append(d.pressed, x)
						}

					case Press:
					case PressToRelease:
						x := k.Keys[d.layer][row][col]

						for i, p := range d.pressed {
							if x == p {
								d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
								pressToRelease = append(pressToRelease, x)
							}
						}
					}
				}
			}
		}

		// read from uart
		for _, k := range d.uk {
			uart := k.uart
			for uart.Buffered() > 0 {
				data, _ := uart.ReadByte()
				buf = append(buf, data)

				if len(buf) == 3 {
					row, col := buf[1], buf[2]
					current := false
					switch buf[0] {
					case 0xAA: // press
						current = true
					case 0x55: // release
						current = false
					default:
						buf[0], buf[1] = buf[1], buf[2]
						buf = buf[:2]
						continue
					}

					switch k.State[row][col] {
					case None:
						if current {
							k.State[row][col] = NoneToPress
						} else {
						}
					case NoneToPress:
						if current {
							k.State[row][col] = Press
						} else {
							k.State[row][col] = PressToRelease
						}
					case Press:
						if current {
						} else {
							k.State[row][col] = PressToRelease
						}
					case PressToRelease:
						if current {
							k.State[row][col] = NoneToPress
						} else {
							k.State[row][col] = None
						}
					}

					switch k.State[row][col] {
					case None:
						// skip
					case NoneToPress:
						x := k.Keys[d.layer][row][col]
						found := false
						for _, p := range d.pressed {
							if x == p {
								found = true
							}
						}
						if !found {
							d.pressed = append(d.pressed, x)
						}
					case Press:
					case PressToRelease:
						x := k.Keys[d.layer][row][col]
						for i, p := range d.pressed {
							if x == p {
								d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
								pressToRelease = append(pressToRelease, x)
							}
						}
					}
					buf = buf[:0]
				}
			}
		}

		for i, x := range d.pressed {
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
					// ここ上手にキーリピートさせたい感じはある
					d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
					pressToRelease = append(pressToRelease, x)
				case 0x05:
					d.Mouse.WheelUp()
					// ここ上手にキーリピートさせたい感じはある
					d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
					pressToRelease = append(pressToRelease, x)
				}
			} else {
				d.Keyboard.Down(k.Keycode(x))
			}
		}

		for _, x := range pressToRelease {
			if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
				if d.modKeyCallback != nil {
					d.modKeyCallback(d.layer, false)
				}
				d.layer = 0

				pressed := []Keycode{}
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
						switch k.Keycode(p) {
						case keycodes.KeyLeftCtrl, keycodes.KeyRightCtrl:
							pressed = append(pressed, p)
						default:
							d.Keyboard.Up(k.Keycode(p))
						}
					}
				}
				d.pressed = d.pressed[:0]
				d.pressed = append(d.pressed, pressed...)

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
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (d *DuplexMatrixKeyboard) Get() [][]State {
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
