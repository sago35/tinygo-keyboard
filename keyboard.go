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
	"golang.org/x/exp/slices"
)

type Device struct {
	Keyboard UpDowner
	Mouse    Mouser
	Override [][]Keycode
	Debug    bool
	flashCh  chan bool
	flashCnt int

	kb []KBer

	layer      int
	layerStack []int
	baseLayer  int
	pressed    []uint32
	repeat     map[uint32]time.Time
}

type KBer interface {
	Get() []State
	Key(layer, index int) Keycode
	SetKeycode(layer, index int, key Keycode)
	GetKeyCount() int
	Init() error
	Callback(layer, index int, state State)
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
		Keyboard:   kb,
		Mouse:      mouse.Port(),
		pressed:    make([]uint32, 0, 10),
		flashCh:    make(chan bool, 10),
		layerStack: make([]int, 0, 6),
		repeat:     map[uint32]time.Time{},
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
	pressToRelease := []uint32{}

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
	noneToPresse := []uint32{}
	for kbidx, k := range d.kb {
		state := k.Get()
		for i := range state {
			switch state[i] {
			case None:
				// skip
			case NoneToPress:
				x := encKey(kbidx, d.layer, i)
				found := false
				for _, p := range d.pressed {
					if x == p {
						found = true
					}
				}
				if !found {
					noneToPresse = append(noneToPresse, x)
					d.pressed = append(d.pressed, x)
				}

			case Press:
			case PressToRelease:
				x := encKey(kbidx, d.layer, i)

				for i, p := range d.pressed {
					if (x & 0xFF00FFFF) == (p & 0xFF00FFFF) {
						d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
						pressToRelease = append(pressToRelease, p)
					}
				}
			}
		}
	}

	for _, xx := range noneToPresse {
		kbidx, layer, index := decKey(xx)
		x := d.kb[kbidx].Key(layer, index)
		if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
			d.layer = int(x) & 0x0F
			if x&keycodes.ToKeyMask == keycodes.ToKeyMask {
				d.baseLayer = d.layer
			} else {
				d.layerStack = append(d.layerStack, d.layer)
			}
		} else if x == keycodes.KeyRestoreDefaultKeymap {
			// restore default keymap for QMK
			machine.Flash.EraseBlocks(0, 1)
		} else if x&0xF000 == 0xD000 {
			switch x & 0x00FF {
			case 0x01, 0x02, 0x04, 0x08, 0x10:
				d.Mouse.Press(mouse.Button(x & 0x00FF))
			case 0x20:
				d.Mouse.WheelDown()
				d.repeat[xx] = time.Now().Add(500 * time.Millisecond)
			case 0x40:
				d.Mouse.WheelUp()
				d.repeat[xx] = time.Now().Add(500 * time.Millisecond)
			}
		} else {
			d.Keyboard.Down(k.Keycode(x))
		}
		d.kb[kbidx].Callback(layer, index, Press)
	}

	for xx, v := range d.repeat {
		if 0 < v.Unix() && v.Sub(time.Now()) < 0 {
			kbidx, layer, index := decKey(xx)
			x := d.kb[kbidx].Key(layer, index)
			if x&0xF000 == 0xD000 {
				switch x & 0x00FF {
				case 0x20:
					d.Mouse.WheelDown()
					d.repeat[xx] = time.Now().Add(100 * time.Millisecond)
				case 0x40:
					d.Mouse.WheelUp()
					d.repeat[xx] = time.Now().Add(100 * time.Millisecond)
				}
			}
		}
	}

	for _, xx := range pressToRelease {
		kbidx, layer, index := decKey(xx)
		x := d.kb[kbidx].Key(layer, index)
		if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
			if x&keycodes.ToKeyMask != keycodes.ToKeyMask {
				layer = int(x) & 0x0F
				idx := slices.Index(d.layerStack, layer)
				slices.Delete(d.layerStack, idx, idx+1)
				d.layerStack = d.layerStack[:len(d.layerStack)-1]
				if len(d.layerStack) == 0 {
					d.layer = d.baseLayer
				} else {
					d.layer = d.layerStack[len(d.layerStack)-1]
				}
			}
		} else if x&0xF000 == 0xD000 {
			switch x & 0x00FF {
			case 0x01, 0x02, 0x04, 0x08, 0x10:
				d.Mouse.Release(mouse.Button(x & 0x00FF))
			case 0x20:
				//d.Mouse.WheelDown()
				d.repeat[xx] = time.Time{}
			case 0x40:
				//d.Mouse.WheelUp()
				d.repeat[xx] = time.Time{}
			}
		} else {
			d.Keyboard.Up(k.Keycode(x))
		}
		d.kb[kbidx].Callback(layer, index, PressToRelease)
	}

	return nil
}

func encKey(kb, layer, index int) uint32 {
	return (uint32(kb) << 24) | (uint32(layer) << 16) | uint32(index)
}

func decKey(k uint32) (int, int, int) {
	kbidx := k >> 24
	layer := (k >> 16) & 0xFF
	index := k & 0x0000FFFF
	return int(kbidx), int(layer), int(index)
}

func (d *Device) Loop(ctx context.Context) error {
	err := d.Init()
	if err != nil {
		return err
	}

	ticker := time.Tick(1 * time.Millisecond)
	cont := true
	for cont {
		select {
		case <-ctx.Done():
			cont = false
			continue
		case <-ticker:
		default:
		}

		err := d.Tick()
		if err != nil {
			return err
		}
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
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF15:
		// TO(x)
		kc = 0x5200 | (kc & 0x000F)
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
	case 0x5200, 0x5201, 0x5202, 0x5203, 0x5204, 0x5205:
		// TO(x)
		kc = 0xFF10 | (kc & 0x000F)
	case 0x5220, 0x5221, 0x5222, 0x5223, 0x5224, 0x5225:
		// MO(x)
		kc = 0xFF00 | (kc & 0x000F)
	case keycodes.KeyRestoreDefaultKeymap:
		kc = keycodes.KeyRestoreDefaultKeymap
	default:
	}

	d.kb[kbIndex].SetKeycode(layer, index, kc)
}

func (d *Device) Layer() int {
	return d.layer
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
