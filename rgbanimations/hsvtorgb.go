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
