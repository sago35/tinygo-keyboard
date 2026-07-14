//go:build tinygo

package keyboard

import (
	k "machine/usb/hid/keyboard"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

// Output indices for SwitchableKeyboard.
const (
	OutputUSB = 0
	OutputBLE = 1
)

// SwitchableKeyboard routes key events to one of two outputs (typically USB
// and BLE) and can switch between them at runtime, either via SetOutput or
// via the keycodes.KeyOutputUSB / KeyOutputBLE / KeyOutputNext keycodes
// handled in Device.Tick. Both outputs are initialized up front so that
// switching is instant. The selected output is persisted to flash by Save()
// and restored by Device.Init(); Default is used when nothing is saved.
type SwitchableKeyboard struct {
	USB UpDowner // typically NewUSBKeyboard()
	BLE UpDowner // typically &BLETxKeyboard{}

	// Default is the output selected when no choice is saved in flash
	// (OutputUSB or OutputBLE).
	Default int

	// OverrideCtrlH translates a Ctrl+H chord into Backspace on whichever
	// output is active, the same as Device.OverrideCtrlH does for a plain
	// USB Keyboard. It has to live here (rather than wrapping d.Keyboard)
	// because Device.Tick reaches d.Keyboard as a concrete *SwitchableKeyboard
	// to handle output switching and unpair.
	OverrideCtrlH bool

	// MuteBLEOnUSB stops advertising and disconnects the current BLE central
	// (if any) whenever USB becomes the active output, and resumes
	// advertising when switching back to BLE. Off by default: switching back
	// to BLE then requires reconnecting rather than being instant. Only takes
	// effect if BLE implements MuteRadio()/UnmuteRadio() (true for
	// *BLETxKeyboard).
	MuteBLEOnUSB bool

	active   int
	pressed  []k.Keycode
	override []k.Keycode
}

func (s *SwitchableKeyboard) Init() error {
	if s.USB != nil {
		err := s.USB.Init()
		if err != nil {
			return err
		}
	}
	if s.BLE != nil {
		err := s.BLE.Init()
		if err != nil {
			return err
		}
	}
	s.active = s.Default
	if s.output() == nil {
		s.active = OutputUSB
	}
	s.applyRadioMute()
	return nil
}

// applyRadioMute mutes or unmutes the BLE radio according to MuteBLEOnUSB and
// the currently active output. It is a no-op unless MuteBLEOnUSB is set and
// BLE implements MuteRadio()/UnmuteRadio() (true for *BLETxKeyboard).
func (s *SwitchableKeyboard) applyRadioMute() {
	if !s.MuteBLEOnUSB {
		return
	}
	r, ok := s.BLE.(interface {
		MuteRadio() error
		UnmuteRadio() error
	})
	if !ok {
		return
	}
	if s.active == OutputUSB {
		r.MuteRadio()
	} else {
		r.UnmuteRadio()
	}
}

func (s *SwitchableKeyboard) output() UpDowner {
	if s.active == OutputBLE {
		return s.BLE
	}
	return s.USB
}

func (s *SwitchableKeyboard) Down(c k.Keycode) error {
	for _, p := range s.pressed {
		if c == p {
			return nil // already pressed
		}
	}
	s.pressed = append(s.pressed, c)

	out := s.output()
	if out == nil {
		return nil
	}

	// Same Ctrl+H -> Backspace translation as Keyboard.Down (keyboard.go),
	// but routed through the active output so it works on USB and BLE alike.
	if s.OverrideCtrlH && len(s.pressed) == 2 &&
		s.pressed[0] == keycodes.KeyLeftCtrl && s.pressed[1] == keycodes.KeyH {
		for _, p := range s.pressed {
			out.Up(p)
		}
		s.override = append(s.override, keycodes.KeyBackspace)
		return out.Down(keycodes.KeyBackspace)
	}

	if len(s.override) > 0 {
		for _, p := range s.override {
			out.Up(p)
		}
		s.override = s.override[:0]
		for _, p := range s.pressed {
			out.Down(p)
		}
	}
	return out.Down(c)
}

func (s *SwitchableKeyboard) Up(c k.Keycode) error {
	out := s.output()

	if out != nil && len(s.override) > 0 {
		for _, p := range s.override {
			out.Up(p)
		}
		s.override = s.override[:0]
		for _, p := range s.pressed {
			// When overriding, do not press the last key again.
			if c != p && p != s.pressed[len(s.pressed)-1] {
				out.Down(p)
			}
		}
	}

	for i, p := range s.pressed {
		if c == p {
			s.pressed = append(s.pressed[:i], s.pressed[i+1:]...)
			if out != nil {
				return out.Up(c)
			}
			return nil
		}
	}
	return nil
}

func (s *SwitchableKeyboard) Write(b []byte) (n int, err error) {
	if out := s.output(); out != nil {
		return out.Write(b)
	}
	return len(b), nil
}

// SetOutput selects the active output (OutputUSB or OutputBLE). All keys
// currently pressed are released on the previous output first, so no key
// stays stuck there. Keys that are physically still held are not re-sent to
// the new output; they take effect from the next press.
func (s *SwitchableKeyboard) SetOutput(n int) {
	if n != OutputUSB && n != OutputBLE {
		return
	}
	if n == s.active {
		return
	}
	if out := s.output(); out != nil {
		// While the Ctrl+H override is active, the key actually held down on
		// the output is the override key (Backspace), not the ones in
		// pressed, so it has to be released too.
		for _, p := range s.override {
			out.Up(p)
		}
		for _, p := range s.pressed {
			out.Up(p)
		}
	}
	s.pressed = s.pressed[:0]
	s.override = s.override[:0]
	s.active = n
	s.applyRadioMute()
}

// Output returns the currently active output (OutputUSB or OutputBLE).
func (s *SwitchableKeyboard) Output() int {
	return s.active
}

// SetBatteryLevel forwards the remaining charge in percent (0-100) to the
// BLE output, if it supports it. USB has no equivalent, so the value is
// reported regardless of the active output.
func (s *SwitchableKeyboard) SetBatteryLevel(percent uint8) error {
	if b, ok := s.BLE.(interface{ SetBatteryLevel(uint8) error }); ok {
		return b.SetBatteryLevel(percent)
	}
	return nil
}

// Unpair forwards keycodes.KeyBluetoothUnpair handling to the BLE output, if
// it supports it.
func (s *SwitchableKeyboard) Unpair() {
	if u, ok := s.BLE.(interface{ Unpair() }); ok {
		u.Unpair()
	}
}
