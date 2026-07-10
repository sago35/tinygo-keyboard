//go:build tinygo

package keyboard

import (
	k "machine/usb/hid/keyboard"
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

	active  int
	pressed []k.Keycode
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
	return nil
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
	return out.Down(c)
}

func (s *SwitchableKeyboard) Up(c k.Keycode) error {
	out := s.output()

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
		for _, p := range s.pressed {
			out.Up(p)
		}
	}
	s.pressed = s.pressed[:0]
	s.active = n
}

// Output returns the currently active output (OutputUSB or OutputBLE).
func (s *SwitchableKeyboard) Output() int {
	return s.active
}

// Unpair forwards keycodes.KeyBluetoothUnpair handling to the BLE output, if
// it supports it.
func (s *SwitchableKeyboard) Unpair() {
	if u, ok := s.BLE.(interface{ Unpair() }); ok {
		u.Unpair()
	}
}
