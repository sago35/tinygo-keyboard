//go:build nrf52840

package ble

import (
	"fmt"
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

	adv := adapter.DefaultAdvertisement()
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo", //k.Name,
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

	registerHID()

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

func connect() error {

	// peerPKey := make([]byte, 0, 64)
	// privLesc, err := ecdh.P256().GenerateKey(rand.Reader)
	// if err != nil {
	// 	return err
	// }
	// lescChan := make(chan struct{})
	bluetooth.SetSecParamsBonding()
	//bluetooth.SetSecParamsLesc()
	bluetooth.SetSecCapabilities(bluetooth.NoneGapIOCapability)
	// time.Sleep(4 * time.Second)
	// println("getting own pub key")
	// var key []byte

	// pk := privLesc.PublicKey().Bytes()
	// pubKey := swapEndinan(pk[1:])
	//bluetooth.SetLesPublicKey(swapEndinan(privLesc.PublicKey().Bytes()[1:]))
	// pubKey = nil
	//println(" key is set")

	// println("register lesc callback")
	// adapter.SetLescRequestHandler(
	// 	func(pubKey []byte) {
	// 		peerPKey = pubKey
	// 		close(lescChan)
	// 	},
	// )

	println("def adv")
	adv := adapter.DefaultAdvertisement()
	println("adv config")
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo-corne",
		ServiceUUIDs: []bluetooth.UUID{
			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
	})
	println("adv start")
	return adv.Start()

	// select {
	// case <-lescChan:
	// 	peerPKey = append([]byte{0x04}, swapEndinan(peerPKey)...)
	// 	p, err := ecdh.P256().NewPublicKey(peerPKey)
	// 	if err != nil {
	// 		println("failed on parsing pub:", err.Error())
	// 		return err
	// 	}
	// 	println("calculating ecdh")
	// 	key, err = privLesc.ECDH(p)
	// 	if err != nil {
	// 		println("failed on curving:", err.Error())
	// 		return errfffffff
	// 	}
	// 	println("key len:", len(key))
	// 	return bluetooth.ReplyLesc(swapEndinan(key))
	// }

}

func swapEndinan(in []byte) []byte {
	var reverse = make([]byte, len(in))
	for i, b := range in[:32] {

		reverse[31-i] = b
	}
	if len(in) > 32 {
		for i, b := range in[32:] {
			reverse[63-i] = b
		}
	}

	return reverse
}

func registerHID() {
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

	//v := []byte{0x85, 0x02} // 0x85, 0x02
	reportValue := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	//var reportmap bluetooth.Characteristic

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
				Value:  reportValue[:],
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
}

func connectToSplit() (bluetooth.DeviceCharacteristic, error) {
	var tx bluetooth.DeviceCharacteristic
	var foundDevice bluetooth.ScanResult
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		fmt.Printf("%q:%#v\n", result.LocalName(), result.Address.String())
		if result.LocalName() != "corne-right" {
			return
		}
		foundDevice = result

		// Stop the scan.
		err := adapter.StopScan()
		if err != nil {
			return
		}
	})
	if err != nil {
		return tx, err
	}
	device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return tx, err
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDNordicUART})
	if err != nil {
		return tx, err
	}
	service := services[0]
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDUARTTX})
	if err != nil {
		return tx, err
	}
	tx = chars[0]
	return tx, nil

}
