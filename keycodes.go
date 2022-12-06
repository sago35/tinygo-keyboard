package main

import (
	k "machine/usb/hid/keyboard"
)

// for Japanese Keyboard
const (
	KeyWindows    k.Keycode = 0xF000 | 0xE3
	KeyHenkan               = 0xF000 | 0x8A
	KeyMuhenkan             = 0xF000 | 0x8B
	KeyAt                   = 0xF000 | 0x2F
	KeyColon                = 0xF000 | 0x34
	KeyBackslash            = 0xF000 | 0x87 // \ |
	KeyBackslash2           = 0xF000 | 0x89 // \ _
	KeyHat                  = 0xF000 | 0x2E
	KeyLeftBrace            = 0xF000 | 0x30
	KeyRightBrace           = 0xF000 | 0x32
)
