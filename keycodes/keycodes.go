package keycodes

const (
	ModKeyMask = 0xFF00
	ToKeyMask  = 0x0010
)

const (
	KeyLeftCtrl   = 0xF000 | 0xE0
	KeyLeftShift  = 0xF000 | 0xE1
	KeyLeftAlt    = 0xF000 | 0xE2
	KeyWindows    = 0xF000 | 0xE3
	KeyRightCtrl  = 0xF000 | 0xE4
	KeyRightShift = 0xF000 | 0xE5

	KeyH         = 0xF000 | 0x0B
	KeyBackspace = 0xF000 | 0x2A
)

const (
	KeyMod0 = ModKeyMask | 0x00
	KeyMod1 = ModKeyMask | 0x01
	KeyMod2 = ModKeyMask | 0x02
	KeyMod3 = ModKeyMask | 0x03
	KeyMod4 = ModKeyMask | 0x04
	KeyMod5 = ModKeyMask | 0x05

	KeyTo0 = ModKeyMask | ToKeyMask | 0x00
	KeyTo1 = ModKeyMask | ToKeyMask | 0x01
	KeyTo2 = ModKeyMask | ToKeyMask | 0x02
	KeyTo3 = ModKeyMask | ToKeyMask | 0x03
	KeyTo4 = ModKeyMask | ToKeyMask | 0x04
	KeyTo5 = ModKeyMask | ToKeyMask | 0x05
)

const (
	// restore default keymap for QMK
	KeyRestoreDefaultKeymap = 0x7C03
)
