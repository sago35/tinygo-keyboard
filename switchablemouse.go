//go:build tinygo

package keyboard

import (
	"machine/usb/hid/mouse"
)

// SwitchableMouse routes mouse events to one of two outputs (typically USB
// and BLE) and can switch between them at runtime via SetOutput. It is the
// Mouser counterpart of SwitchableKeyboard; the two are driven together by
// Device.Tick's KeyOutputUSB/KeyOutputBLE/KeyOutputNext handling so that
// switching keyboard output also switches mouse output. See NewBLEMouse.
type SwitchableMouse struct {
	USB Mouser // typically mouse.Port()
	BLE Mouser // typically the *BLETxKeyboard shared with a SwitchableKeyboard

	// Default is the output selected initially (OutputUSB or OutputBLE).
	Default int

	active  int
	pressed mouse.Button
}

// NewBLEMouse returns a mouse output that shares the BLE connection (and
// HID service) with kb, the *SwitchableKeyboard returned by NewBLEKeyboard.
// It switches USB/BLE output in lockstep with kb, driven by Device.Tick's
// output-switching keycodes. kb.BLE must implement Mouser (true for
// *BLETxKeyboard); otherwise the BLE output is left unset and this behaves
// as a USB-only mouse.
func NewBLEMouse(kb *SwitchableKeyboard) *SwitchableMouse {
	sm := &SwitchableMouse{
		USB:     mouse.Port(),
		Default: kb.Default,
	}
	if m, ok := kb.BLE.(Mouser); ok {
		sm.BLE = m
	}
	sm.active = sm.Default
	if sm.output() == nil {
		sm.active = OutputUSB
	}
	return sm
}

func (s *SwitchableMouse) output() Mouser {
	if s.active == OutputBLE {
		return s.BLE
	}
	return s.USB
}

func (s *SwitchableMouse) Move(vx, vy int) {
	if out := s.output(); out != nil {
		out.Move(vx, vy)
	}
}

func (s *SwitchableMouse) Click(btn mouse.Button) {
	if out := s.output(); out != nil {
		out.Click(btn)
	}
}

func (s *SwitchableMouse) Press(btn mouse.Button) {
	s.pressed |= btn
	if out := s.output(); out != nil {
		out.Press(btn)
	}
}

func (s *SwitchableMouse) Release(btn mouse.Button) {
	s.pressed &^= btn
	if out := s.output(); out != nil {
		out.Release(btn)
	}
}

func (s *SwitchableMouse) Wheel(v int) {
	if out := s.output(); out != nil {
		out.Wheel(v)
	}
}

func (s *SwitchableMouse) WheelDown() {
	if out := s.output(); out != nil {
		out.WheelDown()
	}
}

func (s *SwitchableMouse) WheelUp() {
	if out := s.output(); out != nil {
		out.WheelUp()
	}
}

// SetOutput selects the active output (OutputUSB or OutputBLE). Any mouse
// buttons currently held are released on the previous output first, so no
// button stays stuck there.
func (s *SwitchableMouse) SetOutput(n int) {
	if n != OutputUSB && n != OutputBLE {
		return
	}
	if n == s.active {
		return
	}
	if out := s.output(); out != nil && s.pressed != 0 {
		for _, btn := range []mouse.Button{mouse.Left, mouse.Right, mouse.Middle, mouse.Back, mouse.Forward} {
			if s.pressed&btn != 0 {
				out.Release(btn)
			}
		}
	}
	s.pressed = 0
	s.active = n
}

// Output returns the currently active output (OutputUSB or OutputBLE).
func (s *SwitchableMouse) Output() int {
	return s.active
}
