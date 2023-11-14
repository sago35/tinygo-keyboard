package main

import (
	_ "embed"
	"image/color"
	"log"
	"machine"
	"machine/usb"
	"runtime/volatile"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/bluetooth"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	usb.Product = "xiao-kb01-0.1.0"

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

var (
	white = color.RGBA{0x3F, 0x3F, 0x3F, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

type KeyEvent struct {
	layer, indx int
	state       keyboard.State
}

func run() error {
	var changed volatile.Register8
	changed.Set(0)

	neo := machine.D4
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)
	wsLeds := [12]color.RGBA{}
	for i := range wsLeds {
		wsLeds[i] = black
	}

	d := keyboard.New()

	pins := []machine.Pin{
		machine.D0,
		machine.D1,
		machine.D2,
		machine.D3,
	}

	// --------------------------------------------------
	bluetooth.SetSecParamsBonding()
	bluetooth.SetSecCapabilities(bluetooth.NoneGapIOCapability)

	err := adapter.Enable()
	if err != nil {
		return err
	}
	adv := adapter.DefaultAdvertisement()
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinygo-corne",
		ServiceUUIDs: []bluetooth.UUID{

			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
		/*
		   gatt
		   gacc
		   parameters ?
		   battery service
		   device information
		   hid

		*/
	})
	adv.Start()
	// device information
	// model number string r 0x2A24
	// manufacture name string r  2A29
	// pnp id r 2A50
	//
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
				Value: reportMap, //hidReporMap, // make([]byte, 0, len(HID_REPORT_MAP)),
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("report map")
				},
			},
			{

				Handle: &reportIn,
				UUID:   bluetooth.CharacteristicUUIDReport,
				Value:  reportValue[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("report in")
				},
			},
			{
				// protocl mode
				UUID:  bluetooth.New16BitUUID(0x2A4E),
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicReadPermission,
				Value: []byte{uint8(1)},
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					print("protocol mode")
				},
			},
			{
				UUID:  bluetooth.CharacteristicUUIDHIDControlPoint,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				//	Value: []byte{0x02},
			},
		},
	})
	/*
		d.AddBleKeyboard(24, "Corne-Left", [][]keyboard.Keycode{
			{
				keyboard.Keycode(kb.KeyY), keyboard.Keycode(kb.KeyU), keyboard.Keycode(kb.KeyI), keyboard.Keycode(kb.KeyO), keyboard.Keycode(kb.KeyP), keyboard.Keycode(kb.KeyBackspace),
				keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
				keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
				keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
			},
			{
				keyboard.Keycode(kb.KeyY), keyboard.Keycode(kb.KeyU), keyboard.Keycode(kb.KeyI), keyboard.Keycode(kb.KeyO), keyboard.Keycode(kb.KeyP), keyboard.Keycode(kb.KeyBackspace),
				keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
				keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
				keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
			},
			{
				keyboard.Keycode(kb.Key6), keyboard.Keycode(kb.Key7), keyboard.Keycode(kb.Key8), keyboard.Keycode(kb.Key9), keyboard.Keycode(kb.Key0), keyboard.Keycode(kb.KeyBackspace),
				keyboard.Keycode(kb.KeyH), keyboard.Keycode(kb.KeyJ), keyboard.Keycode(kb.KeyK), keyboard.Keycode(kb.KeyL), keyboard.Keycode(kb.KeySemicolon), jp.KeyColon,
				keyboard.Keycode(kb.KeyN), keyboard.Keycode(kb.KeyM), keyboard.Keycode(kb.KeyComma), keyboard.Keycode(kb.KeyPeriod), keyboard.Keycode(kb.KeyBackslash), keyboard.Keycode(kb.KeyEsc),
				keyboard.Keycode(kb.KeyEnter), kc.KeyMod2, keyboard.Keycode(kb.KeyModifierRightAlt), 0, 0, 0,
			},
		})
	*/
	keyChan := make(chan KeyEvent, 5)
	keyChan2 := make(chan KeyEvent, 5)
	keyChan3 := make(chan KeyEvent, 5)
	keyChan4 := make(chan KeyEvent, 5)
	// --------------------------------------------------

	matrixKeyCodes := [][]keyboard.Keycode{
		{
			jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD,
			jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH,
			jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyL,
		},
	}
	sm := d.AddSquaredMatrixKeyboard(pins, matrixKeyCodes)
	sm.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan <- KeyEvent{layer: layer, indx: index, state: state}
	})

	matrixKeyCodes2 := [][]keyboard.Keycode{
		{
			jp.Key1, jp.Key2,
		},
	}
	rk1 := d.AddRotaryKeyboard(machine.D5, machine.D10, matrixKeyCodes2)
	rk1.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan2 <- KeyEvent{layer: layer, indx: index, state: state}
	})

	matrixKeyCodes3 := [][]keyboard.Keycode{
		{
			jp.Key3, jp.Key4,
		},
	}
	rk2 := d.AddRotaryKeyboard(machine.D9, machine.D8, matrixKeyCodes3)
	rk2.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan3 <- KeyEvent{layer: layer, indx: index, state: state}
	})

	gpioPins := []machine.Pin{machine.D7, machine.D6}
	for c := range gpioPins {
		gpioPins[c].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	matrixKeyCodes4 := [][]keyboard.Keycode{
		{
			jp.Key5, jp.Key6,
		},
	}
	gk := d.AddGpioKeyboard(gpioPins, matrixKeyCodes4)
	gk.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan4 <- KeyEvent{layer: layer, indx: index, state: state}
	})

	go func() {
		for {
			select {
			case x := <-keyChan:
				var report []byte
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mod bits in array
					key := matrixKeyCodes[x.layer][x.indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
			case x := <-keyChan2:
				var report []byte
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mod bits in array
					key := matrixKeyCodes2[x.layer][x.indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
			case x := <-keyChan3:
				var report []byte
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mod bits in array
					key := matrixKeyCodes3[x.layer][x.indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
			case x := <-keyChan4:
				var report []byte
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mod bits in array
					key := matrixKeyCodes4[x.layer][x.indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
			}
		}
	}()

	// for Vial
	loadKeyboardDef()

	err = d.Init()
	if err != nil {
		return err
	}

	cont := true
	ticker := time.Tick(4 * time.Millisecond)
	for cont {
		<-ticker
		err := d.Tick()
		if err != nil {
			return err
		}
		if changed.Get() != 0 {
			ws.WriteColors(wsLeds[:])
			changed.Set(0)
		}
	}

	return nil
}
