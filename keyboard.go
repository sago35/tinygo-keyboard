//go:build tinygo

package keyboard

import (
	"context"
	"machine"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

type Device struct {
	Keyboard UpDowner
	Mouse    Mouser
	Override [][]Keycode
	Debug    bool
	flashCh  chan bool
	flashCnt int

	kb []KBer

	layer   int
	pressed []Keycode
}

type KBer interface {
	Get() []State
	Key(layer, index int) Keycode
	SetKeycode(layer, index int, key Keycode)
	GetKeyCount() int
	Init() error
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

type Callback func(layer, index int, state State)

func New() *Device {
	kb := &Keyboard{
		Port: k.Port(),
	}
	d := &Device{
		Keyboard: kb,
		Mouse:    mouse.Port(),
		pressed:  make([]Keycode, 0, 10),
		flashCh:  make(chan bool, 10),
	}

	SetDevice(d)

	return d
}

func (d *Device) OverrideCtrlH() {
	d.Keyboard = &Keyboard{
		Port:          k.Port(),
		overrideCtrlH: true,
	}
}

func (d *Device) Init() error {
	for _, k := range d.kb {
		err := k.Init()
		if err != nil {
			return err
		}
	}

	d.flashCnt = 0

	// TODO: Allow change to match keyboard
	layers := LayerCount
	keyboards := len(d.kb)
	keys := d.GetMaxKeyCount()

	// TODO: refactor
	rbuf := make([]byte, 4+layers*keyboards*keys*2)
	_, err := machine.Flash.ReadAt(rbuf, 0)
	if err != nil {
		return err
	}
	sz := (int64(rbuf[0]) << 24) + (int64(rbuf[1]) << 16) + (int64(rbuf[2]) << 8) + int64(rbuf[3])
	if sz != machine.Flash.Size() {
		// No settings are saved
		return nil
	}

	offset := 4
	for layer := 0; layer < layers; layer++ {
		for keyboard := 0; keyboard < keyboards; keyboard++ {
			for key := 0; key < keys; key++ {
				kc := Keycode(rbuf[offset+2*key+0]) << 8
				kc += Keycode(rbuf[offset+2*key+1])
				device.SetKeycode(layer, keyboard, key, kc)
			}
			offset += keys * 2
		}
	}

	return nil
}

func (d *Device) GetMaxKeyCount() int {
	cnt := 0
	for _, k := range d.kb {
		if cnt < k.GetKeyCount() {
			cnt = k.GetKeyCount()
		}
	}

	return cnt
}

func (d *Device) Tick() error {
	pressToRelease := []Keycode{}

	select {
	case <-d.flashCh:
		d.flashCnt = 1
	default:
		if d.flashCnt >= 500 {
			d.flashCnt = 0
			err := Save()
			if err != nil {
				return err
			}
		} else if d.flashCnt > 0 {
			d.flashCnt++
		}
	}

	// read from key matrix
	for _, k := range d.kb {
		state := k.Get()
		for i := range state {
			switch state[i] {
			case None:
				// skip
			case NoneToPress:
				x := k.Key(d.layer, i)
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
				x := k.Key(d.layer, i)

				for i, p := range d.pressed {
					if x == p {
						d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
						pressToRelease = append(pressToRelease, x)
					}
				}
			}
		}
	}

	for i, x := range d.pressed {
		if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
			d.layer = int(x) & 0x0F
		} else if x == keycodes.KeyRestoreDefaultKeymap {
			// restore default keymap for QMK
			machine.Flash.EraseBlocks(0, 1)
		} else if x&0xF000 == 0xD000 {
			switch x & 0x00FF {
			case 0x01, 0x02, 0x04, 0x08, 0x10:
				d.Mouse.Press(mouse.Button(x & 0x00FF))
			case 0x20:
				d.Mouse.WheelDown()
				// ここ上手にキーリピートさせたい感じはある
				d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
				pressToRelease = append(pressToRelease, x)
			case 0x40:
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
			d.layer = 0

			pressed := []Keycode{}
			for _, p := range d.pressed {
				if p&0xF000 == 0xD000 {
					switch p & 0x00FF {
					case 0x01, 0x02, 0x04, 0x08, 0x10:
						d.Mouse.Release(mouse.Button(p & 0x00FF))
					case 0x20:
						//d.Mouse.WheelDown()
					case 0x40:
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
			case 0x01, 0x02, 0x04, 0x08, 0x10:
				d.Mouse.Release(mouse.Button(x & 0x00FF))
			case 0x20:
				//d.Mouse.WheelDown()
			case 0x40:
				//d.Mouse.WheelUp()
			}
		} else {
			d.Keyboard.Up(k.Keycode(x))
		}
	}

	return nil
}

func (d *Device) Loop(ctx context.Context) error {
	err := d.Init()
	if err != nil {
		return err
	}

	cont := true
	for cont {
		select {
		case <-ctx.Done():
			cont = false
			continue
		default:
		}

		err := d.Tick()
		if err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (d *Device) Key(layer, kbIndex, index int) Keycode {
	if kbIndex >= len(d.kb) {
		return 0
	}
	return d.kb[kbIndex].Key(layer, index)
}

func (d *Device) KeyVia(layer, kbIndex, index int) Keycode {
	//fmt.Printf("    KeyVia(%d, %d, %d)\n", layer, kbIndex, index)
	if kbIndex >= len(d.kb) {
		return 0
	}
	kc := d.kb[kbIndex].Key(layer, index)
	switch kc {
	case jp.MouseLeft:
		kc = 0x00D1
	case jp.MouseRight:
		kc = 0x00D2
	case jp.MouseMiddle:
		kc = 0x00D3
	case jp.MouseBack:
		kc = 0x00D4
	case jp.MouseForward:
		kc = 0x00D5
	case jp.WheelUp:
		kc = 0x00D9
	case jp.WheelDown:
		kc = 0x00DA
	case jp.KeyMediaVolumeInc:
		kc = 0x00A9
	case jp.KeyMediaVolumeDec:
		kc = 0x00AA
	case 0xFF00, 0xFF01, 0xFF02, 0xFF03, 0xFF04, 0xFF05:
		// MO(x)
		kc = 0x5220 | (kc & 0x000F)
	case keycodes.KeyRestoreDefaultKeymap:
		// restore default keymap for QMK
		kc = keycodes.KeyRestoreDefaultKeymap
	default:
		kc = kc & 0x0FFF
	}
	return kc
}

func (d *Device) SetKeycode(layer, kbIndex, index int, key Keycode) {
	if kbIndex >= len(d.kb) {
		return
	}
	d.kb[kbIndex].SetKeycode(layer, index, key)
}

func (d *Device) SetKeycodeVia(layer, kbIndex, index int, key Keycode) {
	if kbIndex >= len(d.kb) {
		return
	}
	//fmt.Printf("SetKeycodeVia(%d, %d, %d, %04X)\n", layer, kbIndex, index, key)
	kc := key | 0xF000

	switch key {
	case 0x00D1:
		kc = jp.MouseLeft
	case 0x00D2:
		kc = jp.MouseRight
	case 0x00D3:
		kc = jp.MouseMiddle
	case 0x00D4:
		kc = jp.MouseBack
	case 0x00D5:
		kc = jp.MouseForward
	case 0x00D9:
		kc = jp.WheelUp
	case 0x00DA:
		kc = jp.WheelDown
	case 0x00A9:
		kc = jp.KeyMediaVolumeInc
	case 0x00AA:
		kc = jp.KeyMediaVolumeDec
	case 0x5220, 0x5221, 0x5222, 0x5223, 0x5224, 0x5225:
		// MO(x)
		kc = 0xFF00 | (kc & 0x000F)
	case keycodes.KeyRestoreDefaultKeymap:
		kc = keycodes.KeyRestoreDefaultKeymap
	default:
	}

	d.kb[kbIndex].SetKeycode(layer, index, kc)
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

// UartTxKeyboard is a keyboard that simply sends row/col corresponding to key
// placement via UART. For instructions on how to set it up, see bellow.
//
//	./target/sgkb/right
type UartTxKeyboard struct {
	pressed []k.Keycode
	Uart    *machine.UART
}

func (k *UartTxKeyboard) Up(c k.Keycode) error {
	for i, p := range k.pressed {
		if c == p {
			k.pressed = append(k.pressed[:i], k.pressed[i+1:]...)
			row := byte(c >> 8)
			col := byte(c)
			_, err := k.Uart.Write([]byte{0x55, byte(row), byte(col)})
			return err
		}
	}
	return nil
}

func (k *UartTxKeyboard) Down(c k.Keycode) error {
	found := false
	for _, p := range k.pressed {
		if c == p {
			found = true
		}
	}
	if !found {
		k.pressed = append(k.pressed, c)
		row := byte(c >> 8)
		col := byte(c)
		_, err := k.Uart.Write([]byte{0xAA, byte(row), byte(col)})
		return err
	}
	return nil
}
