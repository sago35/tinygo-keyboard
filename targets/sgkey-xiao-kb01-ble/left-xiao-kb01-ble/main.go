package main

import (
	_ "embed"
	"fmt"
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

	machine.InitSerial()
	var err error
	println("enabling adapter")
	err = adapter.Enable()
	if err != nil {
		log.Fatal(err)
	}
	//niceview.ClearScreen()
	rx, err = connectToSplit()
	if err != nil {
		println("failed to connect to other half:", err.Error())
		log.Fatal(err)
	}
	println("connected to right")

	rx.EnableNotifications(
		func(buf []byte) {
			println("recieved buf len:", len(buf))
			rxEvent = buf
			notified = true
		},
	)
	err = connect()
	if err != nil {
		println("failed to establish LESC connection:", err.Error())
		log.Fatal(err)
	}
	println("esteblished LESC connection")
	registerHID()
	println("registered HID")
	// time.Sleep(10 * time.Second)
	// err = adapter.DefaultAdvertisement().Restart()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //adv.Stop()
	// for {

	// 	time.Sleep(100 * time.Millisecond)
	// }
	err = run()
	if err != nil {
		log.Fatal(err)
	}
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

var (
	white = color.RGBA{0x3F, 0x3F, 0x3F, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

type KeyEvent struct {
	layer, indx int
	state       keyboard.State
}

var rxEvent = make([]byte, 0, 3)
var notified bool

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

	keyChan := make(chan KeyEvent, 5)
	keyChan2 := make(chan KeyEvent, 5)
	keyChan3 := make(chan KeyEvent, 5)
	keyChan4 := make(chan KeyEvent, 5)

	matrixKeyCodes := [][]keyboard.Keycode{
		{
			jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD,
			jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH,
			jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyMod1,
		},
		{
			jp.KeyA, jp.KeyB, jp.KeyC, jp.KeyD,
			jp.KeyE, jp.KeyF, jp.KeyG, jp.KeyH,
			jp.KeyI, jp.KeyJ, jp.KeyK, jp.KeyMod1,
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

	splitKb := [][]keyboard.Keycode{
		{
			jp.KeyT, jp.KeyI, jp.KeyN,
			jp.KeyY, jp.KeyG, jp.KeyO,
		},
		{
			jp.KeyB, jp.KeyL, jp.KeyE,
			jp.KeyK, jp.KeyB, jp.KeyD,
		},
	}

	go func() {
		for {
			var report []byte
			select {
			case x := <-keyChan:
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
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
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
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
				if x.state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
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
			default:
				if !notified {
					time.Sleep(1 * time.Millisecond)
					continue
				}
				println("got key")
				layer := int(rxEvent[0])
				indx := int(rxEvent[1])
				state := keyboard.State(rxEvent[2])
				notified = false

				println("key recieved:", indx)
				if state == keyboard.PressToRelease {
					report = []byte{0x01,
						0x00,
						0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
				} else {
					// TODO: actually move special keys to mods bits in array
					key := splitKb[layer][indx]
					report = []byte{0x01,
						0x00,
						0x00,
						uint8(key), 0x00, 0x00, 0x00, 0x00, 0x00}
				}

				_, err := reportIn.Write(report)

				if err != nil {
					println("failed to send key:", err.Error())
				}
				println("sent key")
			}
		}
	}()

	// for Vial
	loadKeyboardDef()

	err := d.Init()
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
