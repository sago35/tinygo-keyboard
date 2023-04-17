//go:build tinygo

package keyboard

import (
	"machine/usb/hid/mouse"
)

type Mouser interface {
	Move(vx, vy int)
	Click(btn mouse.Button)
	Press(btn mouse.Button)
	Release(btn mouse.Button)
	Wheel(v int)
	WheelDown()
	WheelUp()
}
