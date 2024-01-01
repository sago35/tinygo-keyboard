package keycodes

const (
	ModKeyMask = 0xFF00
	ToKeyMask  = 0x0010
)

const (
	TypeMediaKey = 0xE400
	TypeNormal   = 0xF000
	TypeMouse    = 0xD000
	TypeModKey   = 0xFF00
	TypeToKey    = 0xFF10
)

const (
	KeyLeftCtrl   = TypeNormal | 0xE0
	KeyLeftShift  = TypeNormal | 0xE1
	KeyLeftAlt    = TypeNormal | 0xE2
	KeyWindows    = TypeNormal | 0xE3
	KeyRightCtrl  = TypeNormal | 0xE4
	KeyRightShift = TypeNormal | 0xE5

	KeyH         = TypeNormal | 0x0B
	KeyBackspace = TypeNormal | 0x2A
)

const (
	KeyMod0 = TypeModKey | 0x00
	KeyMod1 = TypeModKey | 0x01
	KeyMod2 = TypeModKey | 0x02
	KeyMod3 = TypeModKey | 0x03
	KeyMod4 = TypeModKey | 0x04
	KeyMod5 = TypeModKey | 0x05

	KeyTo0 = TypeToKey | 0x00
	KeyTo1 = TypeToKey | 0x01
	KeyTo2 = TypeToKey | 0x02
	KeyTo3 = TypeToKey | 0x03
	KeyTo4 = TypeToKey | 0x04
	KeyTo5 = TypeToKey | 0x05
)

const (
	// restore default keymap for QMK
	KeyRestoreDefaultKeymap = 0x7C03
)
