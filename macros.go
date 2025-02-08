//go:build tinygo

package keyboard

import (
	"bytes"
	"fmt"
	"time"
)

type MacroDown Keycode

type MacroUp Keycode

// resetMacros resets the macro index table based on the current macro buffer.
// It updates the `d.Macros` array so that each index points to the start of a macro in `d.MacroBuf`.
// Macros are stored sequentially in `d.MacroBuf`, each ending with a null byte (0x00).
func (d *Device) resetMacros() {
	offset := 0
	for i := range d.Macros {
		d.Macros[i] = offset
		next := bytes.IndexByte(d.MacroBuf[offset:], 0x00)
		if next < 0 {
			break
		}
		offset += next + 1
	}
}

// calcMacroSize calculates the required buffer size to store the given macro sequence.
// The macro sequence consists of various types: strings, keycodes, time durations, etc.
// Note: Originally, we wanted to store macros using TGK keycodes, but due to the complexity
// of handling interactions with Vial, we abandoned that approach.
func calcMacroSize(m ...any) int {
	sz := 0
	for _, a := range m {
		switch v := a.(type) {
		case string:
			sz += len(v)
		case time.Duration:
			sz += 4
		case int:
			kc := keycodeTGKtoVia(Keycode(v))
			sz += 3
			if kc > 0xFF {
				sz++
			}
		case Keycode:
			kc := keycodeTGKtoVia(Keycode(v))
			sz += 3
			if kc > 0xFF {
				sz++
			}
		case MacroDown:
			kc := keycodeTGKtoVia(Keycode(v))
			sz += 3
			if kc > 0xFF {
				sz++
			}
		case MacroUp:
			kc := keycodeTGKtoVia(Keycode(v))
			sz += 3
			if kc > 0xFF {
				sz++
			}
		default:
			// Skip unsupported types
		}
	}
	return sz
}

// SetMacros sets a macro at the specified index in `d.MacroBuf`.
// It first validates the index, then updates the macro buffer while ensuring
// proper alignment of stored macros. If the new macro size differs from the old one,
// it shifts the subsequent macros accordingly.
// Note: Originally, we wanted to store macros using TGK keycodes, but due to the complexity
// of handling interactions with Vial, we abandoned that approach.
func (d *Device) SetMacro(index int, m ...any) error {
	if index < 0 || len(d.Macros) <= index {
		return fmt.Errorf("invalid macro index: %d", index)
	}

	d.resetMacros()

	macro := d.MacroBuf[d.Macros[index]:]
	totalSize := 0
	for _, v := range d.Macros {
		idx := bytes.IndexByte(d.MacroBuf[v:], 0x00)
		if idx >= 0 {
			totalSize += idx
		}
	}

	size := calcMacroSize(m...)
	oldSize := bytes.IndexByte(macro, 0x00)
	if size > oldSize {
		// If the new macro is larger, shift the subsequent macros backward
		if index < len(d.Macros) {
			for i := totalSize - 1; i >= d.Macros[index] && (i+size) < len(d.MacroBuf); i-- {
				d.MacroBuf[i+size] = d.MacroBuf[i+oldSize]
			}
		}
		d.resetMacros()
	} else if size < oldSize {
		// If the new macro is smaller, shift the subsequent macros forward
		if index < len(d.Macros) {
			for i := d.Macros[index]; i < totalSize; i++ {
				d.MacroBuf[i+size] = d.MacroBuf[i+oldSize]
			}
		}
		d.resetMacros()
	}

	ofs := 0
	for _, a := range m {
		switch v := a.(type) {
		case string:
			copy(macro[ofs:], []byte(v))
			ofs += len(v)
		case int:
			kc := keycodeTGKtoVia(Keycode(v))
			macro[ofs+0] = 0x01
			macro[ofs+1] = 0x01
			macro[ofs+2] = byte(kc)
			if kc <= 0x00FF {
				ofs += 3
			} else {
				macro[ofs+1] += 4
				macro[ofs+3] = byte(kc >> 8)
				ofs += 4
			}
		case time.Duration:
			ms := v.Milliseconds()
			if ms == 0 && v > 0 {
				ms = 1
			}
			macro[ofs+0] = 0x01
			macro[ofs+1] = 0x04
			macro[ofs+2] = byte(ms%255 + 1)
			macro[ofs+3] = byte(ms/255 + 1)
			ofs += 4
		case Keycode:
			kc := keycodeTGKtoVia(Keycode(v))
			macro[ofs+0] = 0x01
			macro[ofs+1] = 0x01
			macro[ofs+2] = byte(kc)
			if kc <= 0x00FF {
				ofs += 3
			} else {
				macro[ofs+1] += 4
				macro[ofs+3] = byte(kc >> 8)
				ofs += 4
			}
		case MacroDown:
			kc := keycodeTGKtoVia(Keycode(v))
			macro[ofs+0] = 0x01
			macro[ofs+1] = 0x02
			macro[ofs+2] = byte(kc)
			if kc <= 0x00FF {
				ofs += 3
			} else {
				macro[ofs+1] += 4
				macro[ofs+3] = byte(kc >> 8)
				ofs += 4
			}
		case MacroUp:
			kc := keycodeTGKtoVia(Keycode(v))
			macro[ofs+0] = 0x01
			macro[ofs+1] = 0x03
			macro[ofs+2] = byte(kc)
			if kc <= 0x00FF {
				ofs += 3
			} else {
				macro[ofs+1] += 4
				macro[ofs+3] = byte(kc >> 8)
				ofs += 4
			}
		default:
			// Skip unsupported types
		}
	}

	macro[ofs] = 0x00
	ofs++

	// Update macro index table
	for i := index + 1; i < len(d.Macros); i++ {
		d.Macros[i] += ofs
	}

	return nil
}
