//go:build tinygo

package keyboard

import (
	"fmt"
	"machine"
	"machine/usb"
	"machine/usb/descriptor"
)

func init() {
	descriptor.CDCHID.Configuration[2] = 0x84
	descriptor.CDCHID.Configuration[3] = 0x00
	descriptor.CDCHID.Configuration[4] = 0x04

	descriptor.CDCHID.Configuration = append(descriptor.CDCHID.Configuration, []byte{
		// 32 byte
		0x09, 0x04, 0x03, 0x00, 0x02, 0x03, 0x00, 0x00, 0x00,
		0x09, 0x21, 0x11, 0x01, 0x00, 0x01, 0x22, 0x22, 0x00,
		0x07, 0x05, 0x86, 0x03, 0x20, 0x00, 0x01,
		0x07, 0x05, 0x07, 0x03, 0x20, 0x00, 0x01,
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
	txb      [256]byte
	Keys     [][][]Keycode // [row][col]Keycode
	Changed  bool
	Changed2 bool
	wbuf     []byte
)

func rxHandler(b []byte) {
	rxHandler2(b)
}

var keyboardDef = []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00, 0x00, 0x04, 0xE6, 0xD6, 0xB4, 0x46, 0x02, 0x00, 0x21, 0x01, 0x16, 0x00, 0x00, 0x00, 0x74, 0x2F, 0xE5, 0xA3, 0xE0, 0x01, 0x3A, 0x00, 0xB9, 0x5D, 0x00, 0x3D, 0x88, 0x89, 0xC6, 0x54, 0x36, 0xC3, 0x17, 0x4F, 0xE4, 0xFA, 0x84, 0x23, 0x76, 0xB7, 0xFC, 0x71, 0xB0, 0xCD, 0xA7, 0xB6, 0x4B, 0x88, 0x71, 0x53, 0x05, 0x81, 0x49, 0xAB, 0x2B, 0xEF, 0xF4, 0x9E, 0x19, 0x5A, 0x50, 0x3F, 0x4B, 0x6D, 0x55, 0x1E, 0x51, 0xB1, 0xF5, 0x89, 0x00, 0x5C, 0x0D, 0x02, 0xFA, 0xC7, 0x31, 0xFE, 0xC0, 0xB2, 0xA2, 0x91, 0x7D, 0x62, 0xE6, 0xEA, 0xD2, 0x35, 0xA7, 0xE7, 0xA1, 0x84, 0x9D, 0x6E, 0x17, 0x6E, 0x0D, 0xD5, 0xA7, 0x2E, 0xBD, 0xD3, 0x18, 0xC6, 0xB5, 0xD9, 0xE8, 0x88, 0x5D, 0x4D, 0xA5, 0x57, 0x7C, 0x15, 0xB2, 0x12, 0xC5, 0x54, 0xF7, 0x9C, 0x5A, 0x07, 0x21, 0x43, 0x93, 0xF3, 0x7D, 0x55, 0xFB, 0x2F, 0xE8, 0x67, 0x46, 0x1B, 0x34, 0x76, 0x9B, 0x1C, 0xDC, 0x6F, 0x0D, 0x2C, 0x09, 0x80, 0x1A, 0x3F, 0xEF, 0x05, 0x9B, 0x03, 0x7D, 0xAA, 0x10, 0x30, 0x1D, 0x89, 0x87, 0x11, 0x8D, 0x2A, 0xED, 0x53, 0x79, 0xFF, 0x93, 0x25, 0x6C, 0x3C, 0xAC, 0xE3, 0x2C, 0x75, 0x3C, 0x98, 0xA7, 0x12, 0x32, 0xA8, 0x42, 0x4E, 0xC6, 0x3A, 0x9A, 0x54, 0xD7, 0x72, 0xBA, 0x2B, 0x87, 0x65, 0x78, 0x6D, 0x65, 0x00, 0xCC, 0xE9, 0xF6, 0xF3, 0xE9, 0xBF, 0x5E, 0xC8, 0xD3, 0x02, 0x2C, 0x58, 0x5E, 0x81, 0xA0, 0xE6, 0x00, 0x00, 0x00, 0x00, 0x9F, 0x18, 0xB7, 0xA0, 0x13, 0x34, 0x46, 0xD3, 0x00, 0x01, 0xD5, 0x01, 0xBB, 0x02, 0x00, 0x00, 0x04, 0x6D, 0x31, 0x75, 0xB1, 0xC4, 0x67, 0xFB, 0x02, 0x00, 0x00, 0x00, 0x00, 0x04, 0x59, 0x5A}

func rxHandler2(b []byte) bool {
	switch b[0] {
	//case 0x12, 0x0E:
	default:
		fmt.Printf("RxHandler % X\n", b)
	}

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
		//offset := (uint16(b[1]) << 8) + uint16(b[2])
		//sz := b[3]
		//// offset + 0 〜 offset + sz/2 までのデータを返す必要あり
		//// このとき、 JSON で指定した row * col を超えると次の layer のデータとなる
		//// なので、 offset + 0 〜 offset + sz/2 を Keys[layer][row][col] に変換する関数が必要
		//layer := len(Keys)
		//row := len(Keys[0])
		//col := len(Keys[0][0])
		//for i := 0; i < int(sz/2); i++ {
		//	idx := i + int(offset)/2
		//	ll := uint8(idx / (row * col))
		//	rr := uint8((idx % (row * col)) / col)
		//	cc := uint8((idx % (row * col)) % col)
		//	//fmt.Printf("%02X %02X\n", (i+offset/2)/col, (i+offset/2)%col)
		//	if (i+int(offset)/2)/col < row*layer {
		//		txb[i*2+4+0] = byte(Keys[ll][rr][cc] >> 8)
		//		txb[i*2+4+1] = byte(Keys[ll][rr][cc])
		//	} else if offset == 0x0038 {
		//		txb[i*2+4+0] = 0
		//		txb[i*2+4+1] = 5 // B
		//	} else {
		//		//txb[i*2+4+0] = byte((offset + i) >> 8)
		//		//txb[i*2+4+1] = byte(offset + i)
		//		txb[i*2+4+0] = 0
		//		txb[i*2+4+1] = 4 // A
		//	}
		//}
		//if offset+uint16(sz) >= uint16(layer*col*row*2) {
		//	Changed2 = true
		//}

	case 0x0D:
		// DynamicKeymapMacroGetBufferSizeCommand
		txb[1] = 0x07
		txb[2] = 0x9B
	case 0x0C:
		// DynamicKeymapMacroGetCountCommand
		txb[1] = 0x10
	case 0x0E:
		// DynamicKeymapMacroGetBufferCommand
	case 0x02:
		// id_get_keyboard_value
		Changed = false
		Changed2 = false
	case 0x05:
		fmt.Printf("XXXXXXXXX % X\n", b)
		//Keys[b[1]][b[2]][b[3]] = Keycode((uint16(b[4]) << 8) + uint16(b[5]))
		//Changed = true
	case 0x08:
		// id_lighting_get_value
		txb[1] = 0x00
		txb[2] = 0x00
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
		case 0x01:
			// Retrieve keyboard definition size
			size := len(keyboardDef)
			txb[0] = uint8(size)
			txb[1] = uint8(size >> 8)
			txb[2] = uint8(size >> 16)
			txb[3] = uint8(size >> 24)
		case 0x02:
			// Retrieve 32-bytes block of the definition, page ID encoded within 2 bytes
			page := uint16(b[2]) + (uint16(b[3]) << 8)
			start := page * 32
			end := start + 32
			if end < start || int(start) >= len(keyboardDef) {
				return false
			}
			if int(end) > len(keyboardDef) {
				end = uint16(len(keyboardDef))
			}
			fmt.Printf("vial_get_def : page=%04X start=%04X end=%04X\n", page, start, end)
			copy(txb[:32], keyboardDef[start:end])
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
	fmt.Printf("Tx        % X\n", txb[:32])

	return true
}

func Load(l, r, c int) error {
	Keys = make([][][]Keycode, l)
	for ll := range Keys {
		Keys[ll] = make([][]Keycode, r)
		for rr := range Keys[ll] {
			Keys[ll][rr] = make([]Keycode, c)
		}
	}
	wbuf = make([]byte, 4+l*r*c*2)

	keyword := [4]byte{}
	_, err := machine.Flash.ReadAt(keyword[:4], 0)
	if err != nil {
		return err
	}
	fmt.Printf("keyword: % X\n", keyword[:])
	size := (int64(keyword[0]) << 24) +
		(int64(keyword[1]) << 16) +
		(int64(keyword[2]) << 8) +
		(int64(keyword[3]) << 0)
	if false {
		// debug
		size++
	}

	if size != machine.Flash.Size() {
		fmt.Printf("config: not found : %08X vs %08X\n", size, machine.Flash.Size())
	} else {
		fmt.Printf("config: found : %08X vs %08X\n", size, machine.Flash.Size())

		layer := len(Keys)
		row := len(Keys[0])
		col := len(Keys[0][0])
		buf := make([]byte, 2*col)
		offset := int64(0)
		for l := 0; l < layer; l++ {
			for r := 0; r < row; r++ {
				_, err := machine.Flash.ReadAt(buf[:], 4+offset)
				if err != nil {
					return err
				}
				fmt.Printf("% X\n", buf[:])

				for c := 0; c < col; c++ {
					Keys[uint8(l)][r][c] = Keycode((uint16(buf[2*c]) << 8) + uint16(buf[2*c+1]))
					fmt.Printf("Keys[%d][%d][%d] = %04X\n", l, r, c, Keys[uint8(l)][r][c])
				}

				offset += int64(len(buf))
			}
		}
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
