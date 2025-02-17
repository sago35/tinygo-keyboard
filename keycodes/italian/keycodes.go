package italian

import (
	"machine/usb/hid/keyboard"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

func init() {
	keycodes.CharToKeyCodeMap = &CharToKeyCodeMap
}

// for Italian Keyboard
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
	KeySingleQuote = keycodes.TypeNormal | 0x2D
	KeyIGrave      = keycodes.TypeNormal | 0x2E
	KeyEGrave      = keycodes.TypeNormal | 0x2F
	KeyPlus        = keycodes.TypeNormal | 0x30
	KeyUGrave      = keycodes.TypeNormal | 0x32
	KeyOGrave      = keycodes.TypeNormal | 0x33
	KeyAGrave      = keycodes.TypeNormal | 0x34
	KeyBackslash   = keycodes.TypeNormal | 0x35
	KeyComma       = keycodes.TypeNormal | 0x36
	KeyPeriod      = keycodes.TypeNormal | 0x37
	KeyMinus       = keycodes.TypeNormal | 0x38
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
	KeyLessThan    = keycodes.TypeNormal | 0x64
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
	KeyLeftCtrl    = keycodes.TypeNormal | 0xE0
	KeyLeftShift   = keycodes.TypeNormal | 0xE1
	KeyLeftAlt     = keycodes.TypeNormal | 0xE2
	KeyWindows     = keycodes.TypeNormal | 0xE3
	KeyRightCtrl   = keycodes.TypeNormal | 0xE4
	KeyRightShift  = keycodes.TypeNormal | 0xE5
	KeyRightAlt    = keycodes.TypeNormal | 0xE6
)

const (
	KeyMediaBrightnessUp   = keycodes.KeyMediaBrightnessUp
	KeyMediaBrightnessDown = keycodes.KeyMediaBrightnessDown
	KeyMediaPlay           = keycodes.KeyMediaPlay
	KeyMediaPause          = keycodes.KeyMediaPause
	KeyMediaRecord         = keycodes.KeyMediaRecord
	KeyMediaFastForward    = keycodes.KeyMediaFastForward
	KeyMediaRewind         = keycodes.KeyMediaRewind
	KeyMediaNextTrack      = keycodes.KeyMediaNextTrack
	KeyMediaPrevTrack      = keycodes.KeyMediaPrevTrack
	KeyMediaStop           = keycodes.KeyMediaStop
	KeyMediaEject          = keycodes.KeyMediaEject
	KeyMediaRandomPlay     = keycodes.KeyMediaRandomPlay
	KeyMediaPlayPause      = keycodes.KeyMediaPlayPause
	KeyMediaPlaySkip       = keycodes.KeyMediaPlaySkip
	KeyMediaMute           = keycodes.KeyMediaMute
	KeyMediaVolumeInc      = keycodes.KeyMediaVolumeInc
	KeyMediaVolumeDec      = keycodes.KeyMediaVolumeDec
)

const (
	MouseLeft    = keycodes.MouseLeft
	MouseRight   = keycodes.MouseRight
	MouseMiddle  = keycodes.MouseMiddle
	MouseBack    = keycodes.MouseBack
	MouseForward = keycodes.MouseForward
	WheelDown    = keycodes.WheelDown
	WheelUp      = keycodes.WheelUp
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

const (
	KeyMacro0  = keycodes.KeyMacro0
	KeyMacro1  = keycodes.KeyMacro1
	KeyMacro2  = keycodes.KeyMacro2
	KeyMacro3  = keycodes.KeyMacro3
	KeyMacro4  = keycodes.KeyMacro4
	KeyMacro5  = keycodes.KeyMacro5
	KeyMacro6  = keycodes.KeyMacro6
	KeyMacro7  = keycodes.KeyMacro7
	KeyMacro8  = keycodes.KeyMacro8
	KeyMacro9  = keycodes.KeyMacro9
	KeyMacro10 = keycodes.KeyMacro10
	KeyMacro11 = keycodes.KeyMacro11
	KeyMacro12 = keycodes.KeyMacro12
	KeyMacro13 = keycodes.KeyMacro13
	KeyMacro14 = keycodes.KeyMacro14
	KeyMacro15 = keycodes.KeyMacro15
)

var CharToKeyCodeMap = [256]keyboard.Keycode{
	keyboard.ASCII00,
	keyboard.ASCII01,
	keyboard.ASCII02,
	keyboard.ASCII03,
	keyboard.ASCII04,
	keyboard.ASCII05,
	keyboard.ASCII06,
	keyboard.ASCII07,
	keyboard.ASCII08,
	keyboard.ASCII09,
	keyboard.ASCII0A,
	keyboard.ASCII0B,
	keyboard.ASCII0C,
	keyboard.ASCII0D,
	keyboard.ASCII0E,
	keyboard.ASCII0F,
	keyboard.ASCII10,
	keyboard.ASCII11,
	keyboard.ASCII12,
	keyboard.ASCII13,
	keyboard.ASCII14,
	keyboard.ASCII15,
	keyboard.ASCII16,
	keyboard.ASCII17,
	keyboard.ASCII18,
	keyboard.ASCII19,
	keyboard.ASCII1A,
	keyboard.ASCII1B,
	keyboard.ASCII1C,
	keyboard.ASCII1D,
	keyboard.ASCII1E,
	keyboard.ASCII1F,

	KeySpace,                            //  32   SPACE
	Key1 | keycodes.ShiftMask,           //  33   !
	Key2 | keycodes.ShiftMask,           //  34   "
	KeyAGrave | keycodes.AltGrMask,      //  35   #
	Key4 | keycodes.ShiftMask,           //  36   $
	Key5 | keycodes.ShiftMask,           //  37   %
	Key6 | keycodes.ShiftMask,           //  38   &
	KeySingleQuote,                      //  39   '
	Key8 | keycodes.ShiftMask,           //  40   (
	Key9 | keycodes.ShiftMask,           //  41   )
	KeyPlus | keycodes.ShiftMask,        //  42   *
	KeyPlus,                             //  43   +
	KeyComma,                            //  44   ,
	KeyMinus,                            //  45   -
	KeyPeriod,                           //  46   .
	Key7 | keycodes.ShiftMask,           //  47   /
	Key0,                                //  48   0
	Key1,                                //  49   1
	Key2,                                //  50   2
	Key3,                                //  51   3
	Key4,                                //  52   4
	Key5,                                //  53   5
	Key6,                                //  54   6
	Key7,                                //  55   7
	Key8,                                //  55   8
	Key9,                                //  57   9
	KeyPeriod | keycodes.ShiftMask,      //  58   :
	KeyComma | keycodes.ShiftMask,       //  59   ;
	KeyLessThan,                         //  60   <
	Key0 | keycodes.ShiftMask,           //  61   =
	KeyLessThan | keycodes.ShiftMask,    //  62   >
	KeySingleQuote | keycodes.ShiftMask, //  63   ?
	KeyOGrave | keycodes.AltGrMask,      //  64   @
	KeyA | keycodes.ShiftMask,           //  65   A
	KeyB | keycodes.ShiftMask,           //  66   B
	KeyC | keycodes.ShiftMask,           //  67   C
	KeyD | keycodes.ShiftMask,           //  68   D
	KeyE | keycodes.ShiftMask,           //  69   E
	KeyF | keycodes.ShiftMask,           //  70   F
	KeyG | keycodes.ShiftMask,           //  71   G
	KeyH | keycodes.ShiftMask,           //  72   H
	KeyI | keycodes.ShiftMask,           //  73   I
	KeyJ | keycodes.ShiftMask,           //  74   J
	KeyK | keycodes.ShiftMask,           //  75   K
	KeyL | keycodes.ShiftMask,           //  76   L
	KeyM | keycodes.ShiftMask,           //  77   M
	KeyN | keycodes.ShiftMask,           //  78   N
	KeyO | keycodes.ShiftMask,           //  79   O
	KeyP | keycodes.ShiftMask,           //  80   P
	KeyQ | keycodes.ShiftMask,           //  81   Q
	KeyR | keycodes.ShiftMask,           //  82   R
	KeyS | keycodes.ShiftMask,           //  83   S
	KeyT | keycodes.ShiftMask,           //  84   T
	KeyU | keycodes.ShiftMask,           //  85   U
	KeyV | keycodes.ShiftMask,           //  86   V
	KeyW | keycodes.ShiftMask,           //  87   W
	KeyX | keycodes.ShiftMask,           //  88   X
	KeyY | keycodes.ShiftMask,           //  89   Y
	KeyZ | keycodes.ShiftMask,           //  90   Z
	KeyEGrave | keycodes.AltGrMask,      //  91   [
	KeyBackslash,                        //  92   \
	KeyPlus | keycodes.AltGrMask,        //  93   ]
	KeyIGrave | keycodes.ShiftMask,      //  94   ^
	KeyMinus | keycodes.ShiftMask,       //  95   _
	KeySingleQuote | keycodes.AltGrMask, //  96   `
	KeyA,                                //  97   a
	KeyB,                                //  98   b
	KeyC,                                //  99   c
	KeyD,                                // 100   d
	KeyE,                                // 101   e
	KeyF,                                // 102   f
	KeyG,                                // 103   g
	KeyH,                                // 104   h
	KeyI,                                // 105   i
	KeyJ,                                // 106   j
	KeyK,                                // 107   k
	KeyL,                                // 108   l
	KeyM,                                // 109   m
	KeyN,                                // 110   n
	KeyO,                                // 111   o
	KeyP,                                // 112   p
	KeyQ,                                // 113   q
	KeyR,                                // 114   r
	KeyS,                                // 115   s
	KeyT,                                // 116   t
	KeyU,                                // 117   u
	KeyV,                                // 118   v
	KeyW,                                // 119   w
	KeyX,                                // 120   x
	KeyY,                                // 121   y
	KeyZ,                                // 122   z
	KeyEGrave | keycodes.ShiftMask | keycodes.AltGrMask, // 123   {
	KeyBackslash | keycodes.ShiftMask,                   // 124   |
	KeyPlus | keycodes.ShiftMask | keycodes.AltGrMask,   // 125   }
	KeyIGrave | keycodes.AltGrMask,                      // 126   ~
	KeyDelete,                                           // 127   DEL
}
