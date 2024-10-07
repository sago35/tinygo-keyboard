package rgbanimations

// Source: https://stackoverflow.com/a/14733008

func HSVToRGB(h, s, v uint8) (r, g, b, a uint8) {
	if s == 0 {
		return v, v, v, v
	}
	region := h / 43
	remainder := (h - (region * 43)) * 6

	var s16, v16 uint16
	s16 = uint16(s)
	v16 = uint16(v)

	p := uint8((v16 * (255 - s16)) >> 8)
	q := uint8((v16 * (255 - ((s16 * uint16(remainder)) >> 8))) >> 8)
	t := uint8((v16 * (255 - ((s16 * uint16(255-remainder)) >> 8))) >> 8)

	switch region {
	case 0:
		return v, t, p, v
	case 1:
		return q, v, p, v
	case 2:
		return p, v, t, v
	case 3:
		return p, q, v, v
	case 4:
		return t, p, v, v
	default:
		return v, p, q, v
	}
}

func Scale8(i uint8, scale uint8) uint8 {
	return uint8((uint16(i) * uint16(scale)) >> 8)
}

func Scale16by8(i uint16, scale uint8) uint16 {
	return (i * (1 + uint16(scale))) >> 8
}

func Abs8(i int8) uint8 {
	if i < 0 {
		return uint8(-i)
	}
	return uint8(i)
}

func Abs16(i int16) int16 {
	if i < 0 {
		return int16(-i)
	}
	return int16(i)
}

func Sin8(theta uint8) uint8 {
	offset := theta
	if theta&0x40 != 0 {
		offset = 0xFF - offset
	}
	offset &= 0x3F

	secoffset := offset & 0x0F
	if theta&0x40 != 0 {
		secoffset++
	}

	section := offset >> 4
	s2 := section * 2
	var p = []uint8{0, 49, 49, 41, 90, 27, 117, 10}
	b := p[s2]
	m16 := p[s2+1]

	mx := (m16 * secoffset) >> 4

	y := int8(mx + b)
	if theta&0x80 != 0 {
		y = -y
	}

	y += 127
	y++

	return uint8(y)
}

func Atan28(dy int16, dx int16) uint8 {
	if dy == 0 {
		if dx >= 0 {
			return 0
		} else {
			return 128
		}
	}

	var absY int16
	if dy < 0 {
		absY = -dy
	} else {
		absY = dy
	}
	var a int8
	if dx >= 0 {
		a = int8(32 - (32 * (dx - absY) / (dx + absY)))
	} else {
		a = int8(96 - (32 * (dx + absY) / (dx - absY)))
	}
	if dy < 0 {
		return uint8(-a)
	}
	return uint8(a)
}
