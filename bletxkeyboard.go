//go:build tinygo && softdevice

package keyboard

import (
	"fmt"
	"machine/usb"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"runtime/volatile"

	"github.com/sago35/tinygo-keyboard/keycodes"
	"tinygo.org/x/bluetooth"
)

// bleReportMap contains three reports distinguished by report IDs:
//
//   - Report ID 1: keyboard. 8 modifier bits, 1 reserved byte, 5 LED output
//     bits and 6 key codes (the same layout as the USB boot protocol). The
//     key code range is 0-255 so that keys such as the japanese ones
//     (0x87-0x90) can be sent.
//   - Report ID 2: consumer control (media keys). A single 16-bit usage.
//   - Report ID 3: mouse. 5 button bits + 3 padding bits, relative X/Y
//     (8-bit signed each) and relative wheel (8-bit signed). Same 4-byte
//     payload layout as machine/usb/hid/mouse's tx().
//
// Note: over HID over GATT, the report ID is conveyed by the Report
// Reference descriptor of each report characteristic and is NOT prefixed to
// the notification payload.
var bleReportMap = []byte{
	0x05, 0x01, // Usage Page (Generic Desktop)
	0x09, 0x06, // Usage (Keyboard)
	0xA1, 0x01, // Collection (Application)
	0x85, 0x01, //   Report ID (1)
	0x05, 0x07, //   Usage Page (Key Codes)
	0x19, 0xE0, //   Usage Minimum (224)
	0x29, 0xE7, //   Usage Maximum (231)
	0x15, 0x00, //   Logical Minimum (0)
	0x25, 0x01, //   Logical Maximum (1)
	0x75, 0x01, //   Report Size (1)
	0x95, 0x08, //   Report Count (8)
	0x81, 0x02, //   Input (Data, Variable, Absolute): modifier keys
	0x95, 0x01, //   Report Count (1)
	0x75, 0x08, //   Report Size (8)
	0x81, 0x01, //   Input (Constant): reserved byte
	0x95, 0x05, //   Report Count (5)
	0x75, 0x01, //   Report Size (1)
	0x05, 0x08, //   Usage Page (LEDs)
	0x19, 0x01, //   Usage Minimum (1)
	0x29, 0x05, //   Usage Maximum (5)
	0x91, 0x02, //   Output (Data, Variable, Absolute): LEDs
	0x95, 0x01, //   Report Count (1)
	0x75, 0x03, //   Report Size (3)
	0x91, 0x01, //   Output (Constant): padding
	0x95, 0x06, //   Report Count (6)
	0x75, 0x08, //   Report Size (8)
	0x15, 0x00, //   Logical Minimum (0)
	0x26, 0xFF, 0x00, //   Logical Maximum (255)
	0x05, 0x07, //   Usage Page (Key Codes)
	0x19, 0x00, //   Usage Minimum (0)
	0x29, 0xFF, //   Usage Maximum (255)
	0x81, 0x00, //   Input (Data, Array): key codes
	0xC0, // End Collection

	0x05, 0x0C, // Usage Page (Consumer)
	0x09, 0x01, // Usage (Consumer Control)
	0xA1, 0x01, // Collection (Application)
	0x85, 0x02, //   Report ID (2)
	0x15, 0x00, //   Logical Minimum (0)
	0x26, 0xFF, 0x03, //   Logical Maximum (1023)
	0x19, 0x00, //   Usage Minimum (0)
	0x2A, 0xFF, 0x03, //   Usage Maximum (1023)
	0x75, 0x10, //   Report Size (16)
	0x95, 0x01, //   Report Count (1)
	0x81, 0x00, //   Input (Data, Array): media key
	0xC0, // End Collection

	0x05, 0x01, // Usage Page (Generic Desktop)
	0x09, 0x02, // Usage (Mouse)
	0xA1, 0x01, // Collection (Application)
	0x85, 0x03, //   Report ID (3)
	0x09, 0x01, //   Usage (Pointer)
	0xA1, 0x00, //   Collection (Physical)
	0x05, 0x09, //     Usage Page (Button)
	0x19, 0x01, //     Usage Minimum (Button 1)
	0x29, 0x05, //     Usage Maximum (Button 5)
	0x15, 0x00, //     Logical Minimum (0)
	0x25, 0x01, //     Logical Maximum (1)
	0x95, 0x05, //     Report Count (5)
	0x75, 0x01, //     Report Size (1)
	0x81, 0x02, //     Input (Data, Variable, Absolute): buttons
	0x95, 0x01, //     Report Count (1)
	0x75, 0x03, //     Report Size (3)
	0x81, 0x01, //     Input (Constant): padding
	0x05, 0x01, //     Usage Page (Generic Desktop)
	0x09, 0x30, //     Usage (X)
	0x09, 0x31, //     Usage (Y)
	0x15, 0x81, //     Logical Minimum (-127)
	0x25, 0x7F, //     Logical Maximum (127)
	0x75, 0x08, //     Report Size (8)
	0x95, 0x02, //     Report Count (2)
	0x81, 0x06, //     Input (Data, Variable, Relative): X, Y
	0x09, 0x38, //     Usage (Wheel)
	0x15, 0x81, //     Logical Minimum (-127)
	0x25, 0x7F, //     Logical Maximum (127)
	0x75, 0x08, //     Report Size (8)
	0x95, 0x01, //     Report Count (1)
	0x81, 0x06, //     Input (Data, Variable, Relative): wheel
	0xC0, //   End Collection
	0xC0, // End Collection
}

// BLETxKeyboard is a keyboard that sends key events to a BLE central via the
// HID over GATT profile (HOGP). Set Name (optional) and call Init before
// use. For instructions on how to set it up, see below.
//
//	./_x/sgkey-xiao-ble
type BLETxKeyboard struct {
	InputReport    bluetooth.Characteristic
	ConsumerReport bluetooth.Characteristic
	MouseReport    bluetooth.Characteristic
	BatteryLevel   bluetooth.Characteristic

	// Name is the advertised device name. If empty, usb.Product is used
	// (falling back to "tinygo-keyboard" when that is empty too), so the
	// keyboard shows up under the same name over USB and BLE.
	Name string

	adv             *bluetooth.Advertisement
	connected       bool
	connectedDevice bluetooth.Device

	report       [8]byte
	consumer     uint16
	mouseButtons mouse.Button
	battery      uint8

	// The SoftDevice's per-connection notification queue only drains as
	// connection events pass, so sending two reports in the same tick makes
	// Characteristic.Write fail with bluetooth.ErrNotEnoughResources. Each
	// report is a full snapshot of the current state, so nothing is queued:
	// the pending flags mark which reports still have to go out and Flush
	// (called by Device.Tick every tick) resends the then-current state.
	// Mouse X/Y/wheel are relative rather than snapshots and accumulate in
	// pendingX/pendingY/pendingWheel until sent.
	pendingInput    bool
	pendingConsumer bool
	pendingMouse    bool
	pendingX        int
	pendingY        int
	pendingWheel    int

	// Bond profile switching runs on its own goroutine because
	// bluetooth.SelectBondSlot blocks - it waits out the disconnect of the
	// current central and rewrites a flash page, up to a few seconds - and
	// stalling Device.Tick that long would starve the key scan (and trip
	// the watchdog on targets that have one). profileTarget always holds
	// the latest request ("latest wins"); profileSwitching marks the one
	// worker goroutine as alive.
	profileTarget    volatile.Register8
	profileSwitching volatile.Register8
}

// NewBLEKeyboard returns a BLE (HID over GATT) keyboard that can always be
// switched to USB at runtime (see SwitchableKeyboard); it starts on BLE
// unless a different output was saved to flash. The advertised name defaults
// to usb.Product. To customize (advertised name, default output), construct
// a SwitchableKeyboard with a BLETxKeyboard directly.
func NewBLEKeyboard() *SwitchableKeyboard {
	return &SwitchableKeyboard{
		USB:     NewUSBKeyboard(),
		BLE:     &BLETxKeyboard{},
		Default: OutputBLE,
	}
}

// Init enables the BLE stack and pairing, registers the services required by
// the HID over GATT profile (device information, battery and HID), starts
// advertising and fills in InputReport.
func (t *BLETxKeyboard) Init() error {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return fmt.Errorf("failed to enable BLE stack: %w", err)
	}

	// From here on, direct NVMC access (machine.Flash) would hang: keymap
	// persistence has to go through the SoftDevice.
	FlashDevice = bluetooth.SDFlash{}

	// Track the current connection so MuteRadio can disconnect it on demand.
	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		t.connected = connected
		if connected {
			t.connectedDevice = device
		}
	})

	// HID hosts expect bonding; Just Works needs no user interaction.
	err := adapter.EnablePairing(bluetooth.PairingParams{
		IOCapabilities: bluetooth.IOCapsNone,
		LESC:           true,
		PairingCompleteHandler: func(device bluetooth.Device, err error) {
			if err != nil {
				println("pairing failed:", err.Error())
			} else {
				println("pairing complete")
			}
		},
	})
	if err != nil {
		return fmt.Errorf("failed to enable pairing: %w", err)
	}

	// No AllowNewPairing here: while no bond exists, pairing is accepted
	// anyway, and once bonded, a new central may only take over after an
	// explicit Unpair (keycodes.KeyBluetoothUnpair).

	// The device information service with the PnP ID characteristic is
	// required by the HID over GATT profile.
	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDDeviceInformation,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDManufacturerNameString,
				Value: []byte("TinyGo"),
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID: bluetooth.CharacteristicUUIDPnPID,
				// Vendor ID source (0x02 = USB), vendor 0x1915, product
				// 0x0001, version 0x0001 (all little-endian).
				Value: []byte{0x02, 0x15, 0x19, 0x01, 0x00, 0x01, 0x00},
				Flags: bluetooth.CharacteristicReadPermission,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add device information service: %w", err)
	}

	// The battery service is also part of the HID over GATT profile. The
	// level starts at 100% and is updated via SetBatteryLevel.
	t.battery = 100
	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBattery,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &t.BatteryLevel,
				UUID:   bluetooth.CharacteristicUUIDBatteryLevel,
				Value:  []byte{t.battery},
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add battery service: %w", err)
	}

	// The HID service itself. Its characteristics may only be accessed over
	// an encrypted link (Just Works pairing is enough).
	err = adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID: bluetooth.CharacteristicUUIDHIDInformation,
				// bcdHID 1.11, no country code, normally connectable.
				Value: []byte{0x11, 0x01, 0x00, 0x02},
				Flags: bluetooth.CharacteristicReadPermission,
			},
			{
				UUID:         bluetooth.CharacteristicUUIDReportMap,
				Value:        bleReportMap,
				Flags:        bluetooth.CharacteristicReadPermission,
				ReadSecurity: bluetooth.SecurityEncrypted,
			},
			{
				UUID:         bluetooth.CharacteristicUUIDProtocolMode,
				Value:        []byte{1}, // report protocol
				Flags:        bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				ReadSecurity: bluetooth.SecurityEncrypted,
			},
			{
				Handle:       &t.InputReport,
				UUID:         bluetooth.CharacteristicUUIDReport,
				Value:        make([]byte, 8),
				Flags:        bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
				ReadSecurity: bluetooth.SecurityEncrypted,
				Descriptors: []bluetooth.DescriptorConfig{
					{
						UUID:         bluetooth.New16BitUUID(0x2908), // Report Reference
						Value:        []byte{1, 1},                   // report ID 1, input report
						ReadSecurity: bluetooth.SecurityEncrypted,
					},
				},
			},
			{
				UUID:          bluetooth.CharacteristicUUIDReport,
				Value:         []byte{0},
				Flags:         bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				ReadSecurity:  bluetooth.SecurityEncrypted,
				WriteSecurity: bluetooth.SecurityEncrypted,
				// NOTE: WriteEvent is called from the SoftDevice interrupt
				// context. CDC output such as println is forbidden here:
				// when the ring buffer is full it yields inside the
				// interrupt and breaks the scheduler.
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					// Keyboard LED state (NumLock etc.) is not used yet.
				},
				Descriptors: []bluetooth.DescriptorConfig{
					{
						UUID:         bluetooth.New16BitUUID(0x2908), // Report Reference
						Value:        []byte{1, 2},                   // report ID 1, output report
						ReadSecurity: bluetooth.SecurityEncrypted,
					},
				},
			},
			{
				Handle:       &t.ConsumerReport,
				UUID:         bluetooth.CharacteristicUUIDReport,
				Value:        make([]byte, 2),
				Flags:        bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
				ReadSecurity: bluetooth.SecurityEncrypted,
				Descriptors: []bluetooth.DescriptorConfig{
					{
						UUID:         bluetooth.New16BitUUID(0x2908), // Report Reference
						Value:        []byte{2, 1},                   // report ID 2, input report
						ReadSecurity: bluetooth.SecurityEncrypted,
					},
				},
			},
			{
				Handle:       &t.MouseReport,
				UUID:         bluetooth.CharacteristicUUIDReport,
				Value:        make([]byte, 4),
				Flags:        bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
				ReadSecurity: bluetooth.SecurityEncrypted,
				Descriptors: []bluetooth.DescriptorConfig{
					{
						UUID:         bluetooth.New16BitUUID(0x2908), // Report Reference
						Value:        []byte{3, 1},                   // report ID 3, input report
						ReadSecurity: bluetooth.SecurityEncrypted,
					},
				},
			},
			{
				UUID:          bluetooth.CharacteristicUUIDHIDControlPoint,
				Value:         []byte{0},
				Flags:         bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteSecurity: bluetooth.SecurityEncrypted,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add HID service: %w", err)
	}

	name := t.Name
	if name == "" {
		name = usb.Product + " ble"
	}
	if name == "" {
		name = "tinygo-keyboard-ble"
	}
	t.adv = adapter.DefaultAdvertisement()
	err = t.adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    name,
		ServiceUUIDs: []bluetooth.UUID{bluetooth.ServiceUUIDHumanInterfaceDevice},
		Appearance:   961, // keyboard
	})
	if err != nil {
		return fmt.Errorf("failed to config adv: %w", err)
	}
	if err := t.adv.Start(); err != nil {
		return fmt.Errorf("failed to start adv: %w", err)
	}
	println("advertising as", name)

	return nil
}

// MuteRadio stops advertising and disconnects the current central, if any,
// silencing the BLE radio. Used by SwitchableKeyboard.MuteBLEOnUSB to go
// quiet on BLE while USB is the active output. Undo with UnmuteRadio.
func (t *BLETxKeyboard) MuteRadio() error {
	if t.adv == nil {
		return nil
	}
	// Stop advertising first: on disconnect, the SoftDevice auto-restarts
	// advertising if it thinks it should still be advertising, so Stop must
	// happen before Disconnect to prevent that.
	t.adv.Stop() // best-effort: may already be stopped (e.g. while connected)
	if t.connected {
		return t.connectedDevice.Disconnect()
	}
	return nil
}

// UnmuteRadio resumes advertising after MuteRadio.
func (t *BLETxKeyboard) UnmuteRadio() error {
	if t.adv == nil {
		return nil
	}
	return t.adv.Start()
}

func (t *BLETxKeyboard) Down(c k.Keycode) error {
	// The same dispatch as machine/usb/hid/keyboard's Down. Layer keys
	// (MO(x)/TO(x)), mouse keys and macros are handled by Device and never
	// reach here.
	msb := c >> 8
	if 0xE4 <= msb && msb <= 0xE7 {
		// Media key (Consumer Page). The report has a single slot; pressing
		// a second media key while one is held replaces it.
		usage := uint16(c & 0x03FF)
		if t.consumer == usage {
			return nil
		}
		t.consumer = usage
		return t.sendConsumer()
	}
	if msb != keycodes.TypeNormal>>8 {
		return nil
	}
	usage := byte(c)
	if 0xE0 <= usage && usage <= 0xE7 {
		bit := byte(1) << (usage - 0xE0)
		if t.report[0]&bit != 0 {
			return nil
		}
		t.report[0] |= bit
	} else {
		for _, u := range t.report[2:] {
			if u == usage {
				return nil
			}
		}
		found := false
		for i := 2; i < len(t.report); i++ {
			if t.report[i] == 0 {
				t.report[i] = usage
				found = true
				break
			}
		}
		if !found {
			// More than 6 keys are already pressed.
			return nil
		}
	}
	return t.sendInput()
}

func (t *BLETxKeyboard) Up(c k.Keycode) error {
	msb := c >> 8
	if 0xE4 <= msb && msb <= 0xE7 {
		if t.consumer != uint16(c&0x03FF) {
			return nil
		}
		t.consumer = 0
		return t.sendConsumer()
	}
	if msb != keycodes.TypeNormal>>8 {
		return nil
	}
	usage := byte(c)
	changed := false
	if 0xE0 <= usage && usage <= 0xE7 {
		bit := byte(1) << (usage - 0xE0)
		if t.report[0]&bit != 0 {
			t.report[0] &^= bit
			changed = true
		}
	} else {
		for i := 2; i < len(t.report); i++ {
			if t.report[i] == usage {
				t.report[i] = 0
				changed = true
			}
		}
	}
	if !changed {
		return nil
	}
	return t.sendInput()
}

func (t *BLETxKeyboard) sendInput() error {
	_, err := t.InputReport.Write(t.report[:])
	if err == bluetooth.ErrNotEnoughResources {
		t.pendingInput = true
		return nil
	}
	t.pendingInput = false
	return err
}

func (t *BLETxKeyboard) sendConsumer() error {
	_, err := t.ConsumerReport.Write([]byte{byte(t.consumer), byte(t.consumer >> 8)})
	if err == bluetooth.ErrNotEnoughResources {
		t.pendingConsumer = true
		return nil
	}
	t.pendingConsumer = false
	return err
}

// Pending reports whether Flush still has reports to resend.
func (t *BLETxKeyboard) Pending() bool {
	return t.pendingInput || t.pendingConsumer || t.pendingMouse
}

// Flush resends reports whose notification could not be sent because the
// SoftDevice's notification queue was full. The queue drains as connection
// events pass, so calling this every tick (Device.Tick does) converges
// without ever blocking the key scan.
func (t *BLETxKeyboard) Flush() error {
	if t.pendingInput {
		if err := t.sendInput(); err != nil {
			return err
		}
	}
	if t.pendingConsumer {
		if err := t.sendConsumer(); err != nil {
			return err
		}
	}
	if t.pendingMouse {
		if err := t.sendMouse(0, 0, 0); err != nil {
			return err
		}
	}
	return nil
}

func (t *BLETxKeyboard) Write(b []byte) (n int, err error) {
	return len(b), nil
}

// SetBatteryLevel updates the battery service with the remaining charge in
// percent (0-100, larger values are clamped to 100). Subscribed centrals are
// notified of the change. Must be called after Init.
func (t *BLETxKeyboard) SetBatteryLevel(percent uint8) error {
	if percent > 100 {
		percent = 100
	}
	if percent == t.battery {
		return nil
	}
	t.battery = percent
	_, err := t.BatteryLevel.Write([]byte{percent})
	return err
}

// Unpair deletes the stored bond of the active profile and disconnects the
// connected central, if any, so that a new central can pair. It is triggered
// by the keycodes.KeyBluetoothUnpair (BT_UNPR) keycode.
func (t *BLETxKeyboard) Unpair() {
	println("ble: unpair")
	err := bluetooth.DefaultAdapter.RemoveBond()
	if err != nil {
		println("ble: remove bond failed:", err.Error())
	}
}

// SelectProfile switches to the given bond profile (0 to
// bluetooth.BondSlotCount-1): the central bonded to that profile can then
// connect, while the one bonded to the previous profile is disconnected. It
// returns immediately; the switch itself runs on a worker goroutine (see the
// profileTarget field). When called repeatedly the latest request wins.
func (t *BLETxKeyboard) SelectProfile(n int) {
	if n < 0 || n >= bluetooth.BondSlotCount {
		return
	}
	t.profileTarget.Set(uint8(n))
	if t.profileSwitching.Get() != 0 {
		// The worker is alive and rechecks profileTarget after every
		// switch; it picks this request up. (With TinyGo's cooperative
		// scheduler this check cannot race with the worker's exit: the
		// worker only yields inside SelectBondSlot, and rechecks the
		// target after it.)
		return
	}
	t.profileSwitching.Set(1)
	go func() {
		for {
			target := int(t.profileTarget.Get())
			if err := bluetooth.DefaultAdapter.SelectBondSlot(target); err != nil {
				println("ble: select profile:", err.Error())
			}
			if int(t.profileTarget.Get()) == target {
				t.profileSwitching.Set(0)
				return
			}
		}
	}()
}

// BLETxKeyboard also implements Mouser (mouse.go), notifying MouseReport
// (report ID 3, see bleReportMap) instead of sending a USB packet. The report
// layout mirrors machine/usb/hid/mouse's tx(): [buttons, x, y, wheel].

// Move moves the mouse cursor by (vx, vy), clamped to the -127..127 range
// that the report descriptor declares.
func (t *BLETxKeyboard) Move(vx, vy int) {
	if vx == 0 && vy == 0 {
		return
	}
	t.sendMouse(vx, vy, 0)
}

// Click presses and releases the given mouse buttons.
func (t *BLETxKeyboard) Click(btn mouse.Button) {
	t.Press(btn)
	t.Release(btn)
}

// Press presses the given mouse buttons.
func (t *BLETxKeyboard) Press(btn mouse.Button) {
	t.mouseButtons |= btn
	t.sendMouse(0, 0, 0)
}

// Release releases the given mouse buttons.
func (t *BLETxKeyboard) Release(btn mouse.Button) {
	t.mouseButtons &^= btn
	t.sendMouse(0, 0, 0)
}

// Wheel scrolls the mouse wheel by v, clamped to the -127..127 range that
// the report descriptor declares.
func (t *BLETxKeyboard) Wheel(v int) {
	if v == 0 {
		return
	}
	t.sendMouse(0, 0, v)
}

// WheelDown scrolls the mouse wheel down.
func (t *BLETxKeyboard) WheelDown() {
	t.Wheel(-1)
}

// WheelUp scrolls the mouse wheel up.
func (t *BLETxKeyboard) WheelUp() {
	t.Wheel(1)
}

func (t *BLETxKeyboard) sendMouse(x, y, wheel int) error {
	t.pendingX += x
	t.pendingY += y
	t.pendingWheel += wheel
	// Send at most one report's worth of movement; anything beyond the
	// -127..127 range stays pending and goes out on a following tick.
	sx := clampMouseAxis(t.pendingX)
	sy := clampMouseAxis(t.pendingY)
	sw := clampMouseAxis(t.pendingWheel)
	_, err := t.MouseReport.Write([]byte{byte(t.mouseButtons), byte(sx), byte(sy), byte(sw)})
	if err == bluetooth.ErrNotEnoughResources {
		t.pendingMouse = true
		return nil
	}
	if err != nil {
		// Not transient (e.g. disconnected): drop the movement.
		t.pendingX, t.pendingY, t.pendingWheel = 0, 0, 0
		t.pendingMouse = false
		return err
	}
	t.pendingX -= int(sx)
	t.pendingY -= int(sy)
	t.pendingWheel -= int(sw)
	t.pendingMouse = t.pendingX != 0 || t.pendingY != 0 || t.pendingWheel != 0
	return nil
}

func clampMouseAxis(v int) int8 {
	if v < -127 {
		v = -127
	}
	if v > 127 {
		v = 127
	}
	return int8(v)
}
