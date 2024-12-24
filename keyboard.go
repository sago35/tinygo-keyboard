//go:build tinygo

package keyboard

import (
	"bytes"
	"context"
	"fmt"
	"machine"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
	"golang.org/x/exp/slices"
)

type Device struct {
	Keyboard UpDowner
	Mouse    Mouser
	Override [][]Keycode
	Macros   [2048]byte
	Combos   [32][5]Keycode

	Debug    bool
	flashCh  chan bool
	flashCnt int

	kb []KBer

	layer      int
	layerStack []int
	baseLayer  int
	pressed    []uint32
	repeat     map[uint32]time.Time

	combosTimer    time.Time
	combosPressed  map[uint32]struct{}
	combosReleased []uint32
	combosKey      uint32
	combosFounds   []Keycode

	tapOrHold map[uint32]time.Time

	pressToReleaseBuf []uint32
	noneToPressBuf    []uint32
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
	Write(b []byte) (n int, err error)
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

		combosPressed:  map[uint32]struct{}{},
		combosReleased: make([]uint32, 0, 10),
		combosKey:      0xFFFFFFFF,
		combosFounds:   make([]Keycode, 10),

		tapOrHold: map[uint32]time.Time{},

		pressToReleaseBuf: make([]uint32, 0, 20),
		noneToPressBuf:    make([]uint32, 0, 20),
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
	rbuf := make([]byte, 4+layers*keyboards*keys*2+len(device.Macros)+
		len(device.Combos)*len(device.Combos[0])*2)
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

	macroSize := len(device.Macros)
	for i, b := range rbuf[offset : offset+macroSize] {
		if b == 0xFF {
			b = 0
		}
		device.Macros[i] = b
	}
	offset += macroSize

	for idx := range device.Combos {
		device.Combos[idx][0] = Keycode(rbuf[offset+0]) + Keycode(rbuf[offset+1])<<8 // key 1
		device.Combos[idx][1] = Keycode(rbuf[offset+2]) + Keycode(rbuf[offset+3])<<8 // key 2
		device.Combos[idx][2] = Keycode(rbuf[offset+4]) + Keycode(rbuf[offset+5])<<8 // key 3
		device.Combos[idx][3] = Keycode(rbuf[offset+6]) + Keycode(rbuf[offset+7])<<8 // key 4
		device.Combos[idx][4] = Keycode(rbuf[offset+8]) + Keycode(rbuf[offset+9])<<8 // Output key

		// Reinitialize to 0 when reading a value (0xFFFF) from uninitialized flash.
		if device.Combos[idx][0] == 0xFFFF {
			device.Combos[idx][0] = 0x0000
		}
		if device.Combos[idx][1] == 0xFFFF {
			device.Combos[idx][1] = 0x0000
		}
		if device.Combos[idx][2] == 0xFFFF {
			device.Combos[idx][2] = 0x0000
		}
		if device.Combos[idx][3] == 0xFFFF {
			device.Combos[idx][3] = 0x0000
		}
		if device.Combos[idx][4] == 0xFFFF {
			device.Combos[idx][4] = 0x0000
		}

		offset += len(device.Combos[0]) * 2
	}

	return nil
}

func (d *Device) GetKeyboardCount() int {
	return len(d.kb)
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
	pressToRelease := d.pressToReleaseBuf[:0]

	select {
	case <-d.flashCh:
		d.flashCnt = 1
	default:
		if d.flashCnt >= 5000 {
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
	noneToPress := d.noneToPressBuf[:0]
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
					noneToPress = append(noneToPress, x)
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

	d.combosFounds = d.combosFounds[:0]
	if d.combosKey == 0xFFFFFFFF {
		for _, xx := range noneToPress {
			kbidx, layer, index := decKey(xx)
			x := d.kb[kbidx].Key(layer, index)
			for _, combo := range d.Combos {
				for _, ckey := range combo[:4] {
					if keycodeViaToTGK(ckey) == x {
						uniq := true
						for _, f := range d.combosFounds {
							if f == ckey {
								uniq = false
							}
						}
						if uniq {
							d.combosFounds = append(d.combosFounds, ckey)
						}
						if d.combosTimer.IsZero() {
							d.combosTimer = time.Now().Add(48 * time.Millisecond)
						}
						d.combosPressed[xx] = struct{}{}
					}
				}
			}
		}
		if len(d.combosFounds) == len(noneToPress) {
			// Remove the keys pressed before the Combos are completed.
			noneToPress = noneToPress[:0]
		} else {
			// Cancel the Combos waiting state if a key unrelated to Combos is pressed.
			// Sort `d.combosPressed` so that the ones pressed earlier are triggered first.
			d.combosTimer = time.Time{}
			ofs := len(noneToPress)
			noneToPress = d.noneToPressBuf[:ofs+len(d.combosPressed)]
			for i := range noneToPress[:ofs] {
				noneToPress[len(noneToPress)-1-i] = noneToPress[ofs-1-i]
			}
			idx := 0
			for xx := range d.combosPressed {
				noneToPress[idx] = xx
				delete(d.combosPressed, xx)
				idx++
			}
			pressToRelease = append(d.combosReleased, pressToRelease...)
			d.combosReleased = d.combosReleased[:0]
		}
	}

	if !d.combosTimer.IsZero() {
		if time.Now().Before(d.combosTimer) {
			// In the Combos waiting state, only record the events.
			for _, xx := range pressToRelease {
				d.combosReleased = append(d.combosReleased, xx)
			}
			return nil
		} else {
			// When the Combos waiting state is complete, press the key if there is a perfect match.
			// If there is no match, reset `noneToPress` and `pressToRelease`.

			matched := false
			matchMax := 0
			for _, combo := range d.Combos {
				matchCnt := 0
				zero := 0
				for _, ckey := range combo[:4] {
					if ckey == 0x0000 {
						zero++
					} else {
						for xx := range d.combosPressed {
							kbidx, layer, index := decKey(xx)
							x := d.kb[kbidx].Key(layer, index)
							if keycodeViaToTGK(ckey) == x {
								matchCnt++
							}
						}
					}
				}
				if matchCnt >= 2 && zero+matchCnt == 4 && matchCnt > matchMax {
					matched = true
					matchMax = matchCnt
					d.combosKey = 0xFF000000 | uint32(keycodeViaToTGK(combo[4]))
				}
			}

			if matched {
				noneToPress = append(noneToPress, d.combosKey)

				d.combosReleased = d.combosReleased[:0]
			} else {
				for k := range d.combosPressed {
					noneToPress = append(noneToPress, k)
					delete(d.combosPressed, k)
				}
				for _, k := range d.combosReleased {
					pressToRelease = append(pressToRelease, k)
				}
				d.combosReleased = d.combosReleased[:0]
			}
			d.combosTimer = time.Time{}
		}
	}

	for _, xx := range noneToPress {
		kbidx, layer, index := decKey(xx)
		if kbidx < len(d.kb) {
			x := d.kb[kbidx].Key(layer, index)
			switch x & keycodes.QuantumMask {
			case keycodes.TypeLxxxT, keycodes.TypeRxxxT:
				d.tapOrHold[xx] = time.Now().Add(200 * time.Millisecond)
			}
		}
	}

	for xx, tt := range d.tapOrHold {
		if tt.IsZero() {
			// hold release
			for _, yy := range pressToRelease {
				if xx == yy {
					kbidx, layer, index := decKey(xx)
					x := d.kb[kbidx].Key(layer, index)
					switch x & keycodes.QuantumMask {
					case keycodes.TypeLxxxT:
						if x&keycodes.TypeXCtl > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyLeftCtrl))
						}
						if x&keycodes.TypeXSft > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyLeftShift))
						}
						if x&keycodes.TypeXAlt > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyLeftAlt))
						}
						if x&keycodes.TypeXGui > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyWindows))
						}
					case keycodes.TypeRxxxT:
						if x&keycodes.TypeXCtl > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyRightCtrl))
						}
						if x&keycodes.TypeXSft > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyRightShift))
						}
						if x&keycodes.TypeXAlt > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyLeftAlt))
						}
						if x&keycodes.TypeXGui > 0 {
							pressToRelease = append(pressToRelease, uint32(0xFF000000)|uint32(keycodes.KeyWindows))
						}
					}
					delete(d.tapOrHold, xx)
				}
			}
		} else if time.Now().Before(tt) {
			// tap
			for _, yy := range pressToRelease {
				if xx == yy {
					kbidx, layer, index := decKey(xx)
					x := d.kb[kbidx].Key(layer, index)
					switch x & keycodes.QuantumMask {
					case keycodes.TypeLxxxT, keycodes.TypeRxxxT:
						kc := uint32(0xFF000000) | uint32(keycodeViaToTGK(x&0x00FF))
						noneToPress = append(noneToPress, kc)
						pressToRelease = append(pressToRelease, kc)
					}
					delete(d.tapOrHold, xx)
				}
			}
		} else {
			// hold
			kbidx, layer, index := decKey(xx)
			x := d.kb[kbidx].Key(layer, index)
			switch x & keycodes.QuantumMask {
			case keycodes.TypeLxxxT:
				if x&keycodes.TypeXCtl > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyLeftCtrl))
				}
				if x&keycodes.TypeXSft > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyLeftShift))
				}
				if x&keycodes.TypeXAlt > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyLeftAlt))
				}
				if x&keycodes.TypeXGui > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyWindows))
				}
			case keycodes.TypeRxxxT:
				if x&keycodes.TypeXCtl > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyRightCtrl))
				}
				if x&keycodes.TypeXSft > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyRightShift))
				}
				if x&keycodes.TypeXAlt > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyLeftAlt))
				}
				if x&keycodes.TypeXGui > 0 {
					noneToPress = append(noneToPress, uint32(0xFF000000)|uint32(keycodes.KeyWindows))
				}
			}
			d.tapOrHold[xx] = time.Time{}
		}
	}

	for _, xx := range noneToPress {
		kbidx, layer, index := decKey(xx)
		var x Keycode
		if kbidx < len(d.kb) {
			x = d.kb[kbidx].Key(layer, index)
		} else if kbidx >= 0xFF {
			// Combos
			x = Keycode(xx & 0x00FFFFFF)
		} else {
			return fmt.Errorf("kbidx error : %d", kbidx)
		}
		if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
			d.layer = int(x) & 0x0F
			if x&keycodes.ToKeyMask == keycodes.ToKeyMask {
				d.baseLayer = d.layer
			} else {
				d.layerStack = append(d.layerStack, d.layer)
			}
		} else if x&keycodes.QuantumMask == keycodes.TypeLxxx && x&keycodes.QuantumTypeMask != 0 {
			// TypeLxxx
			if x&keycodes.TypeXCtl > 0 {
				d.Keyboard.Down(keycodes.KeyLeftCtrl)
			}
			if x&keycodes.TypeXSft > 0 {
				d.Keyboard.Down(keycodes.KeyLeftShift)
			}
			if x&keycodes.TypeXAlt > 0 {
				d.Keyboard.Down(keycodes.KeyLeftAlt)
			}
			if x&keycodes.TypeXGui > 0 {
				d.Keyboard.Down(keycodes.KeyWindows)
			}
			d.Keyboard.Down(k.Keycode(x&0x00FF | keycodes.TypeNormal))
		} else if x&keycodes.QuantumMask == keycodes.TypeRxxx && x&keycodes.QuantumTypeMask != 0 {
			// TypeRxxx
			if x&keycodes.TypeXCtl > 0 {
				d.Keyboard.Down(keycodes.KeyRightCtrl)
			}
			if x&keycodes.TypeXSft > 0 {
				d.Keyboard.Down(keycodes.KeyRightShift)
			}
			if x&keycodes.TypeXAlt > 0 {
				d.Keyboard.Down(keycodes.KeyLeftAlt)
			}
			if x&keycodes.TypeXGui > 0 {
				d.Keyboard.Down(keycodes.KeyWindows)
			}
			d.Keyboard.Down(k.Keycode(x&0x00FF | keycodes.TypeNormal))
		} else if x == keycodes.KeyRestoreDefaultKeymap {
			// restore default keymap for QMK
			machine.Flash.EraseBlocks(0, 1)
		} else if x&0xFF00 == keycodes.TypeMacroKey {
			no := uint8(x & 0x00FF)
			d.RunMacro(no)
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
		if kbidx < len(d.kb) {
			d.kb[kbidx].Callback(layer, index, Press)
		}
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
		if _, ok := d.combosPressed[xx]; ok {
			delete(d.combosPressed, xx)
		}
	}
	if len(pressToRelease) > 0 && d.combosKey != 0xFFFFFFFF && len(d.combosPressed) == 0 {
		// Combos are deactivated when the last key that constitutes the combo is released.
		pressToRelease = append(pressToRelease, d.combosKey)
		d.combosKey = 0xFFFFFFFF
	}

	for _, xx := range pressToRelease {
		kbidx, layer, index := decKey(xx)
		var x Keycode
		if kbidx < len(d.kb) {
			x = d.kb[kbidx].Key(layer, index)
		} else if kbidx >= 0xFF {
			// Combos
			x = Keycode(xx & 0x00FFFFFF)
		} else {
			return fmt.Errorf("kbidx error : %d", kbidx)
		}
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
		} else if x&keycodes.QuantumMask == keycodes.TypeLxxx && x&keycodes.QuantumTypeMask != 0 {
			// TypeLxxx
			if x&keycodes.TypeXCtl > 0 {
				d.Keyboard.Up(keycodes.KeyLeftCtrl)
			}
			if x&keycodes.TypeXSft > 0 {
				d.Keyboard.Up(keycodes.KeyLeftShift)
			}
			if x&keycodes.TypeXAlt > 0 {
				d.Keyboard.Up(keycodes.KeyLeftAlt)
			}
			if x&keycodes.TypeXGui > 0 {
				d.Keyboard.Up(keycodes.KeyWindows)
			}
			d.Keyboard.Up(k.Keycode(x&0x00FF | keycodes.TypeNormal))
		} else if x&keycodes.QuantumMask == keycodes.TypeRxxx && x&keycodes.QuantumTypeMask != 0 {
			// TypeRxxx
			if x&keycodes.TypeXCtl > 0 {
				d.Keyboard.Up(keycodes.KeyRightCtrl)
			}
			if x&keycodes.TypeXSft > 0 {
				d.Keyboard.Up(keycodes.KeyRightShift)
			}
			if x&keycodes.TypeXAlt > 0 {
				d.Keyboard.Up(keycodes.KeyLeftAlt)
			}
			if x&keycodes.TypeXGui > 0 {
				d.Keyboard.Up(keycodes.KeyWindows)
			}
			d.Keyboard.Up(k.Keycode(x&0x00FF | keycodes.TypeNormal))
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
		if kbidx < len(d.kb) {
			d.kb[kbidx].Callback(layer, index, PressToRelease)
		}
	}

	return nil
}

func (d *Device) RunMacro(no uint8) error {
	macros := bytes.SplitN(d.Macros[:], []byte{0x00}, 16)

	macro := macros[no]

	for i := 0; i < len(macro); {
		if macro[i] == 0x01 {
			p := macro[i:]
			if p[1] == 0x04 {
				// delayMs
				delayMs := int(p[2]) + int(p[3]-1)*255
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
				i += 4
			} else {
				kc := keycodeViaToTGK(Keycode(p[2]))
				sz := 3
				if p[1] > 0x04 {
					kc = Keycode(p[2]) + Keycode(p[3])<<8
					sz += 1
				}
				i += sz
				kc = keycodeViaToTGK(kc)

				switch p[1] {
				case 0x01, 0x05:
					k.Keyboard.Down(k.Keycode(kc))
					k.Keyboard.Up(k.Keycode(kc))
				case 0x02, 0x06:
					k.Keyboard.Down(k.Keycode(kc))
				case 0x03, 0x07:
					k.Keyboard.Up(k.Keycode(kc))
				}
			}
		} else {
			idx := bytes.Index(macro[i:], []byte{0x01})
			if idx == -1 {
				idx = len(macro)
			} else {
				idx = i + idx
			}
			if keycodes.CharToKeyCodeMap != nil {
				for _, b := range macro[i:idx] {
					kc := keycodes.CharToKeyCodeMap[b]
					switch kc & keycodes.ModKeyMask {
					case keycodes.TypeNormal:
						k.Keyboard.Press(kc)
					case keycodes.TypeNormal | keycodes.ShiftMask:
						k.Keyboard.Down(keycodes.KeyLeftShift)
						k.Keyboard.Press(kc ^ keycodes.ShiftMask)
						k.Keyboard.Up(keycodes.KeyLeftShift)
					default:
						//skip
					}
					time.Sleep(10 * time.Millisecond)
				}
			} else {
				k.Keyboard.Write(macro[i:idx])
			}
			i = idx
		}
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
	case keycodes.MouseLeft:
		kc = 0x00D1
	case keycodes.MouseRight:
		kc = 0x00D2
	case keycodes.MouseMiddle:
		kc = 0x00D3
	case keycodes.MouseBack:
		kc = 0x00D4
	case keycodes.MouseForward:
		kc = 0x00D5
	case keycodes.WheelUp:
		kc = 0x00D9
	case keycodes.WheelDown:
		kc = 0x00DA
	case keycodes.KeyMediaBrightnessDown:
		kc = 0x00BE
	case keycodes.KeyMediaBrightnessUp:
		kc = 0x00BD
	case keycodes.KeyMediaMute:
		kc = 0x00A8
	case keycodes.KeyMediaVolumeInc:
		kc = 0x00A9
	case keycodes.KeyMediaVolumeDec:
		kc = 0x00AA
	case keycodes.KeyMediaStop:
		kc = 0x00AD
	case keycodes.KeyMediaPlay:
		kc = 0x00AE
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
		switch kc & keycodes.QuantumMask {
		case keycodes.TypeRxxx, keycodes.TypeLxxxT, keycodes.TypeRxxxT:
			// skip
		default:
			if kc&keycodes.QuantumMask == 0 && kc&keycodes.QuantumTypeMask != 0 {
				// skip (keycodes.TpeLxxx)
			} else {
				switch kc & keycodes.ModKeyMask {
				case keycodes.TypeMacroKey:
					// skip
				default:
					kc = kc & 0x0FFF
				}
			}
		}
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
	kc := keycodeViaToTGK(key)

	d.kb[kbIndex].SetKeycode(layer, index, kc)
}

func keycodeViaToTGK(key Keycode) Keycode {
	kc := key | 0xF000

	switch key {
	case 0x00D1:
		kc = keycodes.MouseLeft
	case 0x00D2:
		kc = keycodes.MouseRight
	case 0x00D3:
		kc = keycodes.MouseMiddle
	case 0x00D4:
		kc = keycodes.MouseBack
	case 0x00D5:
		kc = keycodes.MouseForward
	case 0x00D9:
		kc = keycodes.WheelUp
	case 0x00DA:
		kc = keycodes.WheelDown
	case 0x00BD:
		kc = keycodes.KeyMediaBrightnessUp
	case 0x00BE:
		kc = keycodes.KeyMediaBrightnessDown
	case 0x00A8:
		kc = keycodes.KeyMediaMute
	case 0x00A9:
		kc = keycodes.KeyMediaVolumeInc
	case 0x00AA:
		kc = keycodes.KeyMediaVolumeDec
	case 0x00AD:
		kc = keycodes.KeyMediaStop
	case 0x00AE:
		kc = keycodes.KeyMediaPlay
	case 0x5200, 0x5201, 0x5202, 0x5203, 0x5204, 0x5205:
		// TO(x)
		kc = 0xFF10 | (kc & 0x000F)
	case 0x5220, 0x5221, 0x5222, 0x5223, 0x5224, 0x5225:
		// MO(x)
		kc = 0xFF00 | (kc & 0x000F)
	case keycodes.KeyRestoreDefaultKeymap:
		kc = keycodes.KeyRestoreDefaultKeymap
	default:
		switch key & keycodes.QuantumMask {
		case keycodes.TypeRxxx, keycodes.TypeLxxxT, keycodes.TypeRxxxT:
			kc = key
		default:
			if key&keycodes.QuantumMask == 0 && key&keycodes.QuantumTypeMask != 0 {
				// keycodes.TpeLxxx
				kc = key
			} else {
				switch key & 0xFF00 {
				case keycodes.TypeMacroKey:
					kc = key
				}
			}
		}
	}
	return kc
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

func (k *Keyboard) Write(b []byte) (n int, err error) {
	return k.Port.Write(b)
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

func (k *UartTxKeyboard) Write(b []byte) (n int, err error) {
	return len(b), nil
}
