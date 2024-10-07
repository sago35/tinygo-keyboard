//go:build tinygo

package keyboard

import (
	"machine"
	"machine/usb"
	"machine/usb/descriptor"
)

func init() {
	// vial-gui requires the following magic word.
	usb.Serial = "vial:f64c2b3c"

	descriptor.CDCHID.Configuration[2] = 0x84
	descriptor.CDCHID.Configuration[3] = 0x00
	descriptor.CDCHID.Configuration[4] = 0x04

	descriptor.CDCHID.Configuration = append(descriptor.CDCHID.Configuration, []byte{
		// 32 byte

		// Interface Descriptor
		0x09, 0x04, 0x03, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00,
		// Length: 9 bytes
		// Descriptor Type: Interface (0x04)
		// Interface Number: 3
		// Alternate Setting: 0
		// Number of Endpoints: 2
		// Interface Class: 3 (HID - Human Interface Device)
		// Interface Subclass: 0
		// Interface Protocol: 0
		// Interface String Descriptor Index: 0 (No string descriptor)

		// HID Descriptor
		0x09, 0x21, 0x11, 0x01, 0x00, 0x01, 0x22, 0x22, 0x00,
		// Length: 9 bytes
		// Descriptor Type: HID (0x21)
		// HID Class Specification Release: 1.11
		// Country Code: 0 (Not localized)
		// Number of Descriptors: 1
		// Descriptor Type: Report (0x22)
		// Descriptor Length: 34 bytes (0x0022)

		// Endpoint Descriptor
		0x07, 0x05, 0x86, 0x03, 0x20, 0x00, 0x01,
		// Length: 7 bytes
		// Descriptor Type: Endpoint (0x05)
		// Endpoint Address: 0x86 (Endpoint 6, IN direction)
		// Attributes: 3 (Interrupt transfer type)
		// Maximum Packet Size: 32 bytes (0x0020)
		// Interval: 1 ms

		// Endpoint Descriptor
		0x07, 0x05, 0x07, 0x03, 0x20, 0x00, 0x01,
		// Length: 7 bytes
		// Descriptor Type: Endpoint (0x05)
		// Endpoint Address: 0x07 (Endpoint 7, OUT direction)
		// Attributes: 3 (Interrupt transfer type)
		// Maximum Packet Size: 32 bytes (0x0020)
		// Interval: 1 ms

	}...)

	descriptor.CDCHID.HID[3] = []byte{
		0x06, 0x60, 0xff, // Usage Page (Vendor-Defined 0xFF60)
		0x09, 0x61, // Usage (Vendor-Defined 0x61)
		0xa1, 0x01, // Collection (Application)
		0x09, 0x62, //   Usage (Vendor-Defined 0x62)
		0x15, 0x00, //   Logical Minimum (0)
		0x26, 0xff, 0x00, //   Logical Maximum (255)
		0x95, 0x20, //   Report Count (32)
		0x75, 0x08, //   Report Size (8)
		0x81, 0x02, //   Input (Data, Var, Abs)
		0x09, 0x63, //   Usage (Vendor-Defined 0x63)
		0x15, 0x00, //   Logical Minimum (0)
		0x26, 0xff, 0x00, //   Logical Maximum (255)
		0x95, 0x20, //   Report Count (32)
		0x75, 0x08, //   Report Size (8)
		0x91, 0x02, //   Output (Data, Var, Abs)
		0xc0, // End Collection
	}

	machine.ConfigureUSBEndpoint(descriptor.CDCHID,
		[]usb.EndpointConfig{
			{
				Index:     usb.MIDI_ENDPOINT_OUT,
				IsIn:      false,
				Type:      usb.ENDPOINT_TYPE_INTERRUPT,
				RxHandler: rxHandler,
			},
			{
				Index: usb.MIDI_ENDPOINT_IN,
				IsIn:  true,
				Type:  usb.ENDPOINT_TYPE_INTERRUPT,
			},
		},
		[]usb.SetupConfig{
			{
				Index:   usb.HID_INTERFACE,
				Handler: setupHandler,
			},
		})
}

var (
	txb         [256]byte
	Keys        [][][]Keycode // [row][col]Keycode
	Changed     bool
	Changed2    bool
	wbuf        []byte
	KeyboardDef []byte
	device      *Device
)

func SetDevice(d *Device) {
	device = d
}

func rxHandler(b []byte) {
	rxHandler2(b)
}

func rxHandler2(b []byte) bool {
	copy(txb[:32], b)
	switch b[0] {
	case 0x01:
		// GetProtocolVersionCount
		txb[2] = 0x09
	case 0x11:
		// DynamicKeymapGetLayerCountCommand
		txb[1] = 0x06
	case 0x12:
		// DynamicKeymapReadBufferCommand
		offset := (uint16(b[1]) << 8) + uint16(b[2])
		sz := b[3]
		//fmt.Printf("  offset : %04X + %d\n", offset, sz)
		cnt := device.GetMaxKeyCount()
		for i := 0; i < int(sz/2); i++ {
			//fmt.Printf("  %02X %02X\n", b[4+i+1], b[4+i+0])
			tmp := i + int(offset)/2
			layer := tmp / (cnt * device.GetKeyboardCount())
			tmp = tmp % (cnt * device.GetKeyboardCount())
			kbd := tmp / cnt
			idx := tmp % cnt
			//layer := 0
			//idx := tmp & 0xFF
			kc := device.KeyVia(layer, kbd, idx)
			//fmt.Printf("  (%d, %d, %d)\n", layer, kbd, idx)
			txb[4+2*i+1] = uint8(kc)
			txb[4+2*i+0] = uint8(kc >> 8)
		}

	case 0x0D:
		// DynamicKeymapMacroGetBufferSizeCommand
		sz := len(device.Macros)
		txb[1] = byte(sz >> 8)
		txb[2] = byte(sz)
	case 0x0C:
		// DynamicKeymapMacroGetCountCommand
		txb[1] = 0x10
	case 0x0E:
		// DynamicKeymapMacroGetBufferCommand
		offset := (uint16(b[1]) << 8) + uint16(b[2])
		sz := b[3]
		copy(txb[4:4+sz], device.Macros[offset:])
	case 0x0F:
		// CMD_VIA_MACRO_SET_BUFFER
		offset := (uint16(b[1]) << 8) + uint16(b[2])
		sz := b[3]
		copy(device.Macros[offset:], txb[4:4+sz])
		device.flashCh <- true
	case 0x02:
		// id_get_keyboard_value
		Changed = false
		Changed2 = false
		switch txb[1] {
		case 0x03:
			cols := device.GetMaxKeyCount()
			rowSize := (cols + 7) / 8
			for _, v := range device.pressed {
				row, _, col := decKey(v)
				idx := 2 + row*rowSize + (rowSize - 1) - col/8
				txb[idx] |= byte(1 << (col % 8))
			}
		}
	case 0x05:
		//fmt.Printf("XXXXXXXXX % X\n", b)
		//Keys[b[1]][b[2]][b[3]] = Keycode((uint16(b[4]) << 8) + uint16(b[5]))
		device.SetKeycodeVia(int(b[1]), int(b[2]), int(b[3]), Keycode((uint16(b[4])<<8)+uint16(b[5])))
		device.flashCh <- true
		//Changed = true
	case 0x07:
		// id_lighting_set_value
		if device.IsRGBMatrixEnabled() {
			switch b[1] {
			case 0x41:
				// VIALRGB_SET_MODE
				rgbId := uint16(b[2]) | (uint16(b[3]) << 8)
				device.SetCurrentRGBMode(rgbId)
				device.SetCurrentSpeed(b[4])
				device.SetCurrentHSV(b[5], b[6], b[7])
				device.flashCh <- true
			case 0x42:
				// VIALRGB_DIRECT_FASTSET
				if !device.IsDirectModeEnabled() {
					break
				}
				if len(b) < 4 {
					break
				}
				firstIndex := uint16(b[2]) | (uint16(b[3]) << 8)
				numLeds := uint16(b[4])
				if int(numLeds)*3 > len(b)-4 {
					break
				}
				var i uint16
				for i = 0; i < numLeds; i++ {
					if i+firstIndex >= device.GetRGBMatrixLEDCount() {
						break
					}
					device.SetDirectHSV(
						b[i*3+5],
						b[i*3+6],
						b[i*3+7],
						i+firstIndex,
					)
				}
			}
		}
	case 0x08:
		// id_lighting_get_value
		if device.IsRGBMatrixEnabled() {
			switch b[1] {
			case 0x40:
				// vialrgb_get_info
				const vialRgbProtocolVersion = 0x0001
				txb[2] = vialRgbProtocolVersion & 0xFF
				txb[3] = vialRgbProtocolVersion >> 8
				txb[4] = device.GetRGBMatrixMaximumBrightness()
			case 0x41:
				// vialrgb_get_mode
				currentEffect := device.GetCurrentRGBMode()
				txb[2] = byte(currentEffect & 0xFF)
				txb[3] = byte(currentEffect >> 8)
				txb[4] = device.GetCurrentSpeed()
				txb[5] = device.GetCurrentHue()
				txb[6] = device.GetCurrentSaturation()
				txb[7] = device.GetCurrentValue()
			case 0x42:
				// vialrgb_get_supported
				implementedEffects := device.GetSupportedRGBModes()
				var length int
				if len(implementedEffects) > 15 {
					length = 15
				} else {
					length = len(implementedEffects)
				}
				for i := 0; i < length; i++ {
					txb[i*2] = byte(implementedEffects[i].AnimationType & 0xFF)
					txb[i*2+1] = byte(implementedEffects[i].AnimationType >> 8)
				}
				if length < 16 {
					txb[length*2] = 0xFF
					txb[length*2+1] = 0xFF
				}
			case 0x43:
				// vialrgb_get_number_leds
				if !device.IsDirectModeEnabled() {
					break
				}
				txb[2] = byte(device.GetRGBMatrixLEDCount() & 0xFF)
				txb[3] = byte(device.GetRGBMatrixLEDCount() >> 8)
			case 0x44:
				//vialrgb_get_led_info
				if !device.IsDirectModeEnabled() {
					break
				}
				ledPos := uint16(b[2]&0xFF) | (uint16(b[3]) >> 8)
				position := device.GetRGBMatrixLEDMapping(ledPos)
				// physical position
				txb[2] = position.PhysicalX
				txb[3] = position.PhysicalY
				// flags
				txb[4] = position.LedFlags
				// matrix position
				txb[5] = position.KbIndex
				txb[6] = position.MatrixIndex
			}
		} else {
			txb[1] = 0x00
			txb[2] = 0x00
		}
	case 0xFE: // vial
		switch b[1] {
		case 0x00:
			// Get keyboard ID and Vial protocol version
			const vialProtocolVersion = 0x00000006
			txb[0] = vialProtocolVersion
			txb[1] = vialProtocolVersion >> 8
			txb[2] = vialProtocolVersion >> 16
			txb[3] = vialProtocolVersion >> 24
			txb[4] = 0x9D
			txb[5] = 0xD0
			txb[6] = 0xD5
			txb[7] = 0xE1
			txb[8] = 0x87
			txb[9] = 0xF3
			txb[10] = 0x54
			txb[11] = 0xE2
			if device.IsRGBMatrixEnabled() {
				txb[12] = 1
			}
		case 0x01:
			// Retrieve keyboard definition size
			size := len(KeyboardDef)
			txb[0] = uint8(size)
			txb[1] = uint8(size >> 8)
			txb[2] = uint8(size >> 16)
			txb[3] = uint8(size >> 24)
		case 0x02:
			// Retrieve 32-bytes block of the definition, page ID encoded within 2 bytes
			page := uint16(b[2]) + (uint16(b[3]) << 8)
			start := page * 32
			end := start + 32
			if end < start || int(start) >= len(KeyboardDef) {
				return false
			}
			if int(end) > len(KeyboardDef) {
				end = uint16(len(KeyboardDef))
			}
			//fmt.Printf("vial_get_def : page=%04X start=%04X end=%04X\n", page, start, end)
			copy(txb[:32], KeyboardDef[start:end])
		case 0x09:
			// vial_qmk_settings_query
			// 未対応
			for i := range txb[:32] {
				txb[i] = 0xFF
			}
		case 0x0D:
			// vial_dynamic_entry_op
			txb[0] = 0x00
			txb[1] = 0x00
			txb[2] = 0x00
		case 0x05:
			// vial_get_unlock_status
			txb[0] = 1 // unlocked
			txb[1] = 0 // unlock_in_progress

		default:
		}
	default:
		return false
	}
	machine.SendUSBInPacket(6, txb[:32])

	return true
}

func Save() error {
	layers := 6
	keyboards := device.GetKeyboardCount()

	cnt := device.GetMaxKeyCount()
	rgbStorageSize := 2 + 1 + 3 // currentEffect + speed + HSV
	wbuf := make([]byte, 4+layers*keyboards*cnt*2+len(device.Macros)+rgbStorageSize)
	needed := int64(len(wbuf)) / machine.Flash.EraseBlockSize()
	if needed == 0 {
		needed = 1
	}

	err := machine.Flash.EraseBlocks(0, needed)
	if err != nil {
		return err
	}

	// TODO: Size should be written last
	sz := machine.Flash.Size()
	wbuf[0] = byte(sz >> 24)
	wbuf[1] = byte(sz >> 16)
	wbuf[2] = byte(sz >> 8)
	wbuf[3] = byte(sz >> 0)

	offset := 4
	for layer := 0; layer < layers; layer++ {
		for keyboard := 0; keyboard < keyboards; keyboard++ {
			for key := 0; key < cnt; key++ {
				wbuf[offset+2*key+0] = byte(device.Key(layer, keyboard, key) >> 8)
				wbuf[offset+2*key+1] = byte(device.Key(layer, keyboard, key))
			}
			offset += cnt * 2
		}
	}

	wbuf[offset] = byte(device.GetCurrentRGBMode())
	wbuf[offset+1] = byte(device.GetCurrentRGBMode() >> 8)
	wbuf[offset+2] = device.GetCurrentSpeed()
	wbuf[offset+3] = device.GetCurrentHue()
	wbuf[offset+4] = device.GetCurrentSaturation()
	wbuf[offset+5] = device.GetCurrentValue()
	offset += 6

	copy(wbuf[offset:], device.Macros[:])

	_, err = machine.Flash.WriteAt(wbuf[:], 0)
	if err != nil {
		return err
	}

	return nil
}

func setupHandler(setup usb.Setup) bool {
	ok := false
	if setup.BmRequestType == usb.SET_REPORT_TYPE && setup.BRequest == usb.SET_IDLE {
		machine.SendZlp()
		ok = true
	}
	return ok
}
