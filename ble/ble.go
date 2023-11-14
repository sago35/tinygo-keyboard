//go:build nrf52840

package ble

import (
	k "machine/usb/hid/keyboard"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var reportIn bluetooth.Characteristic
var rx bluetooth.DeviceCharacteristic

var reportMap = []byte{
	0x05, 0x01, // USAGE_PAGE (Generic Desktop)
	0x09, 0x06, // USAGE (Keyboard)
	0xa1, 0x01, // COLLECTION (Application)
	0x85, 0x01, //   REPORT_ID (1)
	0x05, 0x07, //   USAGE_PAGE (Keyboard)
	0x19, 0x01, //   USAGE_MINIMUM
	0x29, 0x7f, //   USAGE_MAXIMUM
	0x15, 0x00, //   LOGICAL_MINIMUM (0)
	0x25, 0x01, //   LOGICAL_MAXIMUM (1)
	0x75, 0x01, //   REPORT_SIZE (1)
	0x95, 0x08, //   REPORT_COUNT (8)
	0x81, 0x02, //   INPUT (Data,Var,Abs)
	0x95, 0x01, //   REPORT_COUNT (1)
	0x75, 0x08, //   REPORT_SIZE (8)
	0x81, 0x01, //   INPUT (Cnst,Ary,Abs)
	0x95, 0x06, //   REPORT_COUNT (6)
	0x75, 0x08, //   REPORT_SIZE (8)
	0x15, 0x00, //   LOGICAL_MINIMUM (0)
	0x25, 0x65, //   LOGICAL_MAXIMUM (101)
	0x05, 0x07, //   USAGE_PAGE (Keyboard)
	0x19, 0x00, //   USAGE_MINIMUM (Reserved (no event indicated))
	0x29, 0x65, //   USAGE_MAXIMUM (Keyboard Application)
	0x81, 0x00, //   INPUT (Data,Ary,Abs)
	0xc0, // END_COLLECTION
}

type Keyboard struct {
	Name   string
	report [9]byte
}

func NewKeyboard(name string) *Keyboard {
	return &Keyboard{
		Name: name,
	}
}

func (k *Keyboard) Connect() error {
	err := adapter.Enable()
	if err != nil {
		return err
	}

	bluetooth.SetSecParamsBonding()
	bluetooth.SetSecCapabilities(bluetooth.NoneGapIOCapability)

	name := k.Name
	if len(name) > 14 {
		name = name[:14]
	}

	adv := adapter.DefaultAdvertisement()
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: name,
		ServiceUUIDs: []bluetooth.UUID{
			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
	})

	err = adv.Start()
	if err != nil {
		return err
	}

	k.registerHID()

	return nil
}

func (k *Keyboard) registerHID() error {
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDDeviceInformation,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDManufacturerNameString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("Nice Keyboards"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDModelNumberString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("nice!nano"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDPnPID,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{0x02, 0x8a, 0x24, 0x66, 0x82, 0x34, 0x36},
				//Value: []byte{0x02, uint8(0x10C4 >> 8), uint8(0x10C4 & 0xff), uint8(0x0001 >> 8), uint8(0x0001 & 0xff)},
			},
		},
	})
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBattery,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDBatteryLevel,
				Value: []byte{80},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	// gacc
	/*
	   device name r
	   apperance r
	   peripheral prefreed connection

	*/

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDGenericAccess,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDDeviceName,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("tinygo-corne"),
			},
			{

				UUID:  bluetooth.New16BitUUID(0x2A01),
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{uint8(0x03c4 >> 8), uint8(0x03c4 & 0xff)}, /// []byte(strconv.Itoa(961)),
			},
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDPeripheralPreferredConnectionParameters,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{0x02},
			// },

			// // 		//
		},
	})

	// hid
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		/*
			 - hid information r
			 - report map r
			 - report nr
			   - client charecteristic configuration
			   - report reference
			- report nr
			   - client charecteristic configuration
			   - report reference
			- hid control point wnr
		*/
		Characteristics: []bluetooth.CharacteristicConfig{
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDHIDInformation,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{uint8(0x0111 >> 8), uint8(0x0111 & 0xff), uint8(0x0002 >> 8), uint8(0x0002 & 0xff)},
			// },
			{
				//Handle: &reportmap,
				UUID:  bluetooth.CharacteristicUUIDReportMap,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: reportMap,
			},
			{

				Handle: &reportIn,
				UUID:   bluetooth.CharacteristicUUIDReport,
				Value:  k.report[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				// protocl mode
				UUID:  bluetooth.New16BitUUID(0x2A4E),
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicReadPermission,
				// Value: []byte{uint8(1)},
				// WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
				// 	print("protocol mode")
				// },
			},
			{
				UUID:  bluetooth.CharacteristicUUIDHIDControlPoint,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				//	Value: []byte{0x02},
			},
		},
	})

	return nil
}

func (k *Keyboard) Up(c k.Keycode) error {
	k.report[0] = 0x01
	k.report[1] = 0x00
	k.report[2] = 0x00
	k.report[3] = 0x00
	k.report[4] = 0x00
	k.report[5] = 0x00
	k.report[6] = 0x00
	k.report[7] = 0x00
	k.report[8] = 0x00
	_, err := reportIn.Write(k.report[:9])
	return err
}

func (k *Keyboard) Down(c k.Keycode) error {
	k.report[0] = 0x01
	k.report[1] = 0x00
	k.report[2] = 0x00
	k.report[3] = uint8(c)
	k.report[4] = 0x00
	k.report[5] = 0x00
	k.report[6] = 0x00
	k.report[7] = 0x00
	k.report[8] = 0x00
	_, err := reportIn.Write(k.report[:9])
	return err
}
