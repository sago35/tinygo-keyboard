package main

import (
	"image/color"
)

type FB struct {
	buf [320 * 240 * 2]byte
}

func (fb *FB) Size() (x, y int16) {
	return 320, 240
}

func (fb *FB) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || 320 <= x || y < 0 || 240 <= y {
		return
	}
	c565 := RGBATo565(c)
	fb.buf[(int(y)*320+int(x))*2+0] = byte(c565 >> 8)
	fb.buf[(int(y)*320+int(x))*2+1] = byte(c565)
}

func (fb *FB) Display() error {
	return nil
}

func (fb *FB) FillScreen(c color.RGBA) {
	c565 := RGBATo565(c)
	c1 := uint8(c565 >> 8)
	c2 := uint8(c565)

	for i := range fb.buf {
		if i%2 == 0 {
			fb.buf[i] = c1
		} else {
			fb.buf[i] = c2
		}
	}
}

// RGBATo565 converts a color.RGBA to uint16 used in the display
func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}
