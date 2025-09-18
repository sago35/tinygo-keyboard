//go:build tinygo

package keyboard

import (
	"machine"
	"machine/usb"
	"machine/usb/descriptor"
)

const (
	EPxIN  = usb.MIDI_ENDPOINT_IN
	EPxOUT = usb.MIDI_ENDPOINT_OUT
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
		0x07, 0x05, 0x80 | EPxIN, 0x03, 0x20, 0x00, 0x01,
		// Length: 7 bytes
		// Descriptor Type: Endpoint (0x05)
		// Endpoint Address: 0x8x (Endpoint X, IN direction)
		// Attributes: 3 (Interrupt transfer type)
		// Maximum Packet Size: 32 bytes (0x0020)
		// Interval: 1 ms

		// Endpoint Descriptor
		0x07, 0x05, 0x00 | EPxOUT, 0x03, 0x20, 0x00, 0x01,
		// Length: 7 bytes
		// Descriptor Type: Endpoint (0x05)
		// Endpoint Address: 0x0x (Endpoint X, OUT direction)
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
				Index:     EPxOUT,
				IsIn:      false,
				Type:      usb.ENDPOINT_TYPE_INTERRUPT,
				RxHandler: rxHandler,
			},
			{
				Index: EPxIN,
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
	switch b[0] {
	//case viaCommandDynamicKeymapGetBuffer, viaCommandDynamicKeymapMacroGetBuffer:
	default:
		//fmt.Printf("RxHandler % X\n", b)
	}

	copy(txb[:32], b)
	switch b[0] {
	case viaCommandGetProtocolVersion:
		txb[2] = 0x09
	case viaCommandDynamicKeymapGetLayerCount:
		txb[1] = 0x06
	case viaCommandDynamicKeymapGetBuffer:
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

	case viaCommandDynamicKeymapMacroGetBufferSize:
		sz := len(device.MacroBuf)
		txb[1] = byte(sz >> 8)
		txb[2] = byte(sz)
	case viaCommandDynamicKeymapMacroGetCount:
		txb[1] = 0x10
	case viaCommandDynamicKeymapMacroGetBuffer:
		offset := (uint16(b[1]) << 8) + uint16(b[2])
		sz := b[3]
		copy(txb[4:4+sz], device.MacroBuf[offset:])
	case viaCommandDynamicKeymapMacroSetBuffer:
		offset := (uint16(b[1]) << 8) + uint16(b[2])
		sz := b[3]
		copy(device.MacroBuf[offset:], txb[4:4+sz])
		device.flashCh <- true
	case viaCommandGetKeyboardValue:
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
	case viaCommandDynamicKeymapSetKeycode:
		//fmt.Printf("XXXXXXXXX % X\n", b)
		//Keys[b[1]][b[2]][b[3]] = Keycode((uint16(b[4]) << 8) + uint16(b[5]))
		device.SetKeycodeVia(int(b[1]), int(b[2]), int(b[3]), Keycode((uint16(b[4])<<8)+uint16(b[5])))
		device.flashCh <- true
		//Changed = true
	case viaCommandLightingGetValue:
		txb[1] = 0x00
		txb[2] = 0x00
	case viaCommandVialPrefix: // vial
		switch b[1] {
		case vialGetKeyboardId:
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
		case vialGetSize:
			// Retrieve keyboard definition size
			size := len(KeyboardDef)
			txb[0] = uint8(size)
			txb[1] = uint8(size >> 8)
			txb[2] = uint8(size >> 16)
			txb[3] = uint8(size >> 24)
		case vialGetDef:
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
		case vialQmkSettingsQuery:
			// not impl
			for i := range txb[:32] {
				txb[i] = 0xFF
			}
		case vialDynamicEntryOp:
			switch b[2] {
			case dynamicVialGetNumberOfEntries:
				txb[0] = 0x00
				txb[1] = 0x20 // combos
				txb[2] = 0x00
			case dynamicVialComboGet:
				txb[0] = 0x00
				idx := b[3]
				txb[1] = byte(keycodeTGKtoVia(device.Combos[idx][0]))
				txb[2] = byte(keycodeTGKtoVia(device.Combos[idx][0]) >> 8)
				txb[3] = byte(keycodeTGKtoVia(device.Combos[idx][1]))
				txb[4] = byte(keycodeTGKtoVia(device.Combos[idx][1]) >> 8)
				txb[5] = byte(keycodeTGKtoVia(device.Combos[idx][2]))
				txb[6] = byte(keycodeTGKtoVia(device.Combos[idx][2]) >> 8)
				txb[7] = byte(keycodeTGKtoVia(device.Combos[idx][3]))
				txb[8] = byte(keycodeTGKtoVia(device.Combos[idx][3]) >> 8)
				txb[9] = byte(keycodeTGKtoVia(device.Combos[idx][4]))
				txb[10] = byte(keycodeTGKtoVia(device.Combos[idx][4]) >> 8)
				// 00 0400 0500 0000 0000 0700 000000000000000000000000000000000000000000
				// 0  1    3    5    7    9
			case dynamicVialComboSet:
				txb[0] = 0x00
				idx := b[3]
				// fe0d04 00 0400 0500 0000 0000 0700 000000000000000000000000000000000000
				// 0 1 2  3  4    6    8    10   12
				device.Combos[idx][0] = keycodeViaToTGK(Keycode(b[4]) + Keycode(b[5])<<8)   // key 1
				device.Combos[idx][1] = keycodeViaToTGK(Keycode(b[6]) + Keycode(b[7])<<8)   // key 2
				device.Combos[idx][2] = keycodeViaToTGK(Keycode(b[8]) + Keycode(b[9])<<8)   // key 3
				device.Combos[idx][3] = keycodeViaToTGK(Keycode(b[10]) + Keycode(b[11])<<8) // key 4
				device.Combos[idx][4] = keycodeViaToTGK(Keycode(b[12]) + Keycode(b[13])<<8) // Output key
				device.flashCh <- true
			default:
				txb[0] = 0x00
				txb[1] = 0x00
				txb[2] = 0x00
			}
		case vialGetUnlockStatus:
			txb[0] = 1 // unlocked
			txb[1] = 0 // unlock_in_progress

		default:
		}
	default:
		return false
	}
	machine.SendUSBInPacket(EPxIN, txb[:32])
	//fmt.Printf("Tx        % X\n", txb[:32])

	return true
}

func Save() error {
	layers := 6
	keyboards := device.GetKeyboardCount()

	cnt := device.GetMaxKeyCount()
	wbuf := make([]byte, 4+layers*keyboards*cnt*2+len(device.MacroBuf)+
		len(device.Combos)*len(device.Combos[0])*2)
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

	macroSize := len(device.MacroBuf)
	copy(wbuf[offset:offset+macroSize], device.MacroBuf[:])
	offset += macroSize

	for _, combo := range device.Combos {
		wbuf[offset+0] = byte(combo[0])
		wbuf[offset+1] = byte(combo[0] >> 8)
		wbuf[offset+2] = byte(combo[1])
		wbuf[offset+3] = byte(combo[1] >> 8)
		wbuf[offset+4] = byte(combo[2])
		wbuf[offset+5] = byte(combo[2] >> 8)
		wbuf[offset+6] = byte(combo[3])
		wbuf[offset+7] = byte(combo[3] >> 8)
		wbuf[offset+8] = byte(combo[4])
		wbuf[offset+9] = byte(combo[4] >> 8)
		offset += len(device.Combos[0]) * 2
	}

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

// https://github.com/vial-kb/vial-qmk/quantum/via.h
const (
	viaCommandGetProtocolVersion              = 0x01 // always 0x01
	viaCommandGetKeyboardValue                = 0x02
	viaCommandSetKeyboardValue                = 0x03
	viaCommandDynamicKeymapGetKeycode         = 0x04
	viaCommandDynamicKeymapSetKeycode         = 0x05
	viaCommandDynamicKeymapReset              = 0x06
	viaCommandCustomSetValue                  = 0x07
	viaCommandCustomGetValue                  = 0x08
	viaCommandCustomSave                      = 0x09
	viaCommandLightingSetValue                = 0x07
	viaCommandLightingGetValue                = 0x08
	viaCommandLightingSave                    = 0x09
	viaCommandEepromReset                     = 0x0A
	viaCommandBootloaderJump                  = 0x0B
	viaCommandDynamicKeymapMacroGetCount      = 0x0C
	viaCommandDynamicKeymapMacroGetBufferSize = 0x0D
	viaCommandDynamicKeymapMacroGetBuffer     = 0x0E
	viaCommandDynamicKeymapMacroSetBuffer     = 0x0F
	viaCommandDynamicKeymapMacroReset         = 0x10
	viaCommandDynamicKeymapGetLayerCount      = 0x11
	viaCommandDynamicKeymapGetBuffer          = 0x12
	viaCommandDynamicKeymapSetBuffer          = 0x13
	viaCommandVialPrefix                      = 0xFE
	viaCommandUnhandled                       = 0xFF
)

// https://github.com/vial-kb/vial-qmk/quantum/vial.h
const (
	vialGetKeyboardId    = 0x00
	vialGetSize          = 0x01
	vialGetDef           = 0x02
	vialGetEncoder       = 0x03
	vialSetEncoder       = 0x04
	vialGetUnlockStatus  = 0x05
	vialUnlockStart      = 0x06
	vialUnlockPoll       = 0x07
	vialLock             = 0x08
	vialQmkSettingsQuery = 0x09
	vialQmkSettingsGet   = 0x0A
	vialQmkSettingsSet   = 0x0B
	vialQmkSettingsReset = 0x0C
	vialDynamicEntryOp   = 0x0D /* operate on tapdance, combos, etc */

	dynamicVialGetNumberOfEntries = 0x00
	dynamicVialTapDanceGet        = 0x01
	dynamicVialTapDanceSet        = 0x02
	dynamicVialComboGet           = 0x03
	dynamicVialComboSet           = 0x04
	dynamicVialKeyOverrideGet     = 0x05
	dynamicVialKeyOverrideSet     = 0x06
)
