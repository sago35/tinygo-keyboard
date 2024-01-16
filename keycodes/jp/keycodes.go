package jp

import (
	"github.com/sago35/tinygo-keyboard/keycodes"
)

// for Japanese Keyboard
// based on machine/usb/hid/keyboard/keycode.go
const (
	KeyA           = keycodes.TypeNormal | 0x04
	KeyB           = keycodes.TypeNormal | 0x05
	KeyC           = keycodes.TypeNormal | 0x06
	KeyD           = keycodes.TypeNormal | 0x07
	KeyE           = keycodes.TypeNormal | 0x08
	KeyF           = keycodes.TypeNormal | 0x09
	KeyG           = keycodes.TypeNormal | 0x0A
	KeyH           = keycodes.TypeNormal | 0x0B
	KeyI           = keycodes.TypeNormal | 0x0C
	KeyJ           = keycodes.TypeNormal | 0x0D
	KeyK           = keycodes.TypeNormal | 0x0E
	KeyL           = keycodes.TypeNormal | 0x0F
	KeyM           = keycodes.TypeNormal | 0x10
	KeyN           = keycodes.TypeNormal | 0x11
	KeyO           = keycodes.TypeNormal | 0x12
	KeyP           = keycodes.TypeNormal | 0x13
	KeyQ           = keycodes.TypeNormal | 0x14
	KeyR           = keycodes.TypeNormal | 0x15
	KeyS           = keycodes.TypeNormal | 0x16
	KeyT           = keycodes.TypeNormal | 0x17
	KeyU           = keycodes.TypeNormal | 0x18
	KeyV           = keycodes.TypeNormal | 0x19
	KeyW           = keycodes.TypeNormal | 0x1A
	KeyX           = keycodes.TypeNormal | 0x1B
	KeyY           = keycodes.TypeNormal | 0x1C
	KeyZ           = keycodes.TypeNormal | 0x1D
	Key1           = keycodes.TypeNormal | 0x1E
	Key2           = keycodes.TypeNormal | 0x1F
	Key3           = keycodes.TypeNormal | 0x20
	Key4           = keycodes.TypeNormal | 0x21
	Key5           = keycodes.TypeNormal | 0x22
	Key6           = keycodes.TypeNormal | 0x23
	Key7           = keycodes.TypeNormal | 0x24
	Key8           = keycodes.TypeNormal | 0x25
	Key9           = keycodes.TypeNormal | 0x26
	Key0           = keycodes.TypeNormal | 0x27
	KeyEnter       = keycodes.TypeNormal | 0x28
	KeyEsc         = keycodes.TypeNormal | 0x29
	KeyBackspace   = keycodes.TypeNormal | 0x2A
	KeyTab         = keycodes.TypeNormal | 0x2B
	KeySpace       = keycodes.TypeNormal | 0x2C
	KeyMinus       = keycodes.TypeNormal | 0x2D
	KeyHat         = keycodes.TypeNormal | 0x2E
	KeyAt          = keycodes.TypeNormal | 0x2F
	KeyLeftBrace   = keycodes.TypeNormal | 0x30
	KeyRightBrace  = keycodes.TypeNormal | 0x32
	KeySemicolon   = keycodes.TypeNormal | 0x33
	KeyColon       = keycodes.TypeNormal | 0x34
	KeyHankaku     = keycodes.TypeNormal | 0x35
	KeyComma       = keycodes.TypeNormal | 0x36
	KeyPeriod      = keycodes.TypeNormal | 0x37
	KeySlash       = keycodes.TypeNormal | 0x38
	KeyCapsLock    = keycodes.TypeNormal | 0x39
	KeyF1          = keycodes.TypeNormal | 0x3A
	KeyF2          = keycodes.TypeNormal | 0x3B
	KeyF3          = keycodes.TypeNormal | 0x3C
	KeyF4          = keycodes.TypeNormal | 0x3D
	KeyF5          = keycodes.TypeNormal | 0x3E
	KeyF6          = keycodes.TypeNormal | 0x3F
	KeyF7          = keycodes.TypeNormal | 0x40
	KeyF8          = keycodes.TypeNormal | 0x41
	KeyF9          = keycodes.TypeNormal | 0x42
	KeyF10         = keycodes.TypeNormal | 0x43
	KeyF11         = keycodes.TypeNormal | 0x44
	KeyF12         = keycodes.TypeNormal | 0x45
	KeyPrintscreen = keycodes.TypeNormal | 0x46
	KeyScrollLock  = keycodes.TypeNormal | 0x47
	KeyPause       = keycodes.TypeNormal | 0x48
	KeyInsert      = keycodes.TypeNormal | 0x49
	KeyHome        = keycodes.TypeNormal | 0x4A
	KeyPageUp      = keycodes.TypeNormal | 0x4B
	KeyDelete      = keycodes.TypeNormal | 0x4C
	KeyEnd         = keycodes.TypeNormal | 0x4D
	KeyPageDown    = keycodes.TypeNormal | 0x4E
	KeyRight       = keycodes.TypeNormal | 0x4F
	KeyLeft        = keycodes.TypeNormal | 0x50
	KeyDown        = keycodes.TypeNormal | 0x51
	KeyUp          = keycodes.TypeNormal | 0x52
	KeyNumLock     = keycodes.TypeNormal | 0x53
	KeypadSlash    = keycodes.TypeNormal | 0x54
	KeypadAsterisk = keycodes.TypeNormal | 0x55
	KeypadMinus    = keycodes.TypeNormal | 0x56
	KeypadPlus     = keycodes.TypeNormal | 0x57
	KeypadEnter    = keycodes.TypeNormal | 0x58
	Keypad1        = keycodes.TypeNormal | 0x59
	Keypad2        = keycodes.TypeNormal | 0x5A
	Keypad3        = keycodes.TypeNormal | 0x5B
	Keypad4        = keycodes.TypeNormal | 0x5C
	Keypad5        = keycodes.TypeNormal | 0x5D
	Keypad6        = keycodes.TypeNormal | 0x5E
	Keypad7        = keycodes.TypeNormal | 0x5F
	Keypad8        = keycodes.TypeNormal | 0x60
	Keypad9        = keycodes.TypeNormal | 0x61
	Keypad0        = keycodes.TypeNormal | 0x62
	KeypadPeriod   = keycodes.TypeNormal | 0x63
	KeyNonUSBS     = keycodes.TypeNormal | 0x64
	KeyMenu        = keycodes.TypeNormal | 0x65
	KeyF13         = keycodes.TypeNormal | 0x68
	KeyF14         = keycodes.TypeNormal | 0x69
	KeyF15         = keycodes.TypeNormal | 0x6A
	KeyF16         = keycodes.TypeNormal | 0x6B
	KeyF17         = keycodes.TypeNormal | 0x6C
	KeyF18         = keycodes.TypeNormal | 0x6D
	KeyF19         = keycodes.TypeNormal | 0x6E
	KeyF20         = keycodes.TypeNormal | 0x6F
	KeyF21         = keycodes.TypeNormal | 0x70
	KeyF22         = keycodes.TypeNormal | 0x71
	KeyF23         = keycodes.TypeNormal | 0x72
	KeyF24         = keycodes.TypeNormal | 0x73
	KeyBackslash   = keycodes.TypeNormal | 0x87 // \ |
	KeyHiragana    = keycodes.TypeNormal | 0x88
	KeyBackslash2  = keycodes.TypeNormal | 0x89 // \ _
	KeyHenkan      = keycodes.TypeNormal | 0x8A
	KeyMuhenkan    = keycodes.TypeNormal | 0x8B
	KeyKana        = keycodes.TypeNormal | 0x90
	KeyEisu        = keycodes.TypeNormal | 0x91
	KeyLeftCtrl    = keycodes.TypeNormal | 0xE0
	KeyLeftShift   = keycodes.TypeNormal | 0xE1
	KeyLeftAlt     = keycodes.TypeNormal | 0xE2
	KeyWindows     = keycodes.TypeNormal | 0xE3
	KeyRightCtrl   = keycodes.TypeNormal | 0xE4
	KeyRightShift  = keycodes.TypeNormal | 0xE5
)

const (
	KeyMediaPlay        = keycodes.TypeMediaKey | 0xB0
	KeyMediaPause       = keycodes.TypeMediaKey | 0xB1
	KeyMediaRecord      = keycodes.TypeMediaKey | 0xB2
	KeyMediaFastForward = keycodes.TypeMediaKey | 0xB3
	KeyMediaRewind      = keycodes.TypeMediaKey | 0xB4
	KeyMediaNextTrack   = keycodes.TypeMediaKey | 0xB5
	KeyMediaPrevTrack   = keycodes.TypeMediaKey | 0xB6
	KeyMediaStop        = keycodes.TypeMediaKey | 0xB7
	KeyMediaEject       = keycodes.TypeMediaKey | 0xB8
	KeyMediaRandomPlay  = keycodes.TypeMediaKey | 0xB9
	KeyMediaPlayPause   = keycodes.TypeMediaKey | 0xCD
	KeyMediaPlaySkip    = keycodes.TypeMediaKey | 0xCE
	KeyMediaMute        = keycodes.TypeMediaKey | 0xE2
	KeyMediaVolumeInc   = keycodes.TypeMediaKey | 0xE9
	KeyMediaVolumeDec   = keycodes.TypeMediaKey | 0xEA
)

const (
	MouseLeft    = keycodes.TypeMouse | 0x01 // mouse.Left
	MouseRight   = keycodes.TypeMouse | 0x02 // mouse.Right
	MouseMiddle  = keycodes.TypeMouse | 0x04 // mouse.Middle
	MouseBack    = keycodes.TypeMouse | 0x08 // mouse.Back
	MouseForward = keycodes.TypeMouse | 0x10 // mouse.Forward
	WheelDown    = keycodes.TypeMouse | 0x20
	WheelUp      = keycodes.TypeMouse | 0x40
)

const (
	KeyMod0 = keycodes.KeyMod0
	KeyMod1 = keycodes.KeyMod1
	KeyMod2 = keycodes.KeyMod2
	KeyMod3 = keycodes.KeyMod3
	KeyMod4 = keycodes.KeyMod4
	KeyMod5 = keycodes.KeyMod5

	KeyTo0 = keycodes.KeyTo0
	KeyTo1 = keycodes.KeyTo1
	KeyTo2 = keycodes.KeyTo2
	KeyTo3 = keycodes.KeyTo3
	KeyTo4 = keycodes.KeyTo4
	KeyTo5 = keycodes.KeyTo5
)
