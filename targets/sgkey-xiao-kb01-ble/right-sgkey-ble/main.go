package main

import (
	_ "embed"
	"log"
	"machine"
	"machine/usb"
	"time"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
	"tinygo.org/x/bluetooth"
)

func main() {
	usb.Product = "sgkey-0.1.0"

	err := advertise()
	if err != nil {
		log.Fatal(err)
	}
	err = run()
	if err != nil {
		log.Fatal(err)
	}
}

func advertise() error {

	err := adapter.Enable()
	if err != nil {
		return err
	}
	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "corne-right",
	})
	err = adv.Start()
	if err != nil {
		return err
	}

	var buf = make([]byte, 0, 3)

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDNordicUART,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: tx,
				UUID:   bluetooth.CharacteristicUUIDUARTTX,
				Value:  buf[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	return nil
}

type RCS struct {
	row, col int
	state    keyboard.State
}

type KeyEvent struct {
	layer, indx int
	state       keyboard.State
}

func run() error {
	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D8,
		machine.D9,
		machine.D10,
	}

	rowPins := []machine.Pin{
		machine.D1,
		machine.D2,
	}

	keyChan := make(chan KeyEvent, 10)

	matrixKeyCodes := [][]keyboard.Keycode{
		{
			jp.KeyT, jp.KeyI, jp.KeyN,
			jp.KeyY, jp.KeyG, jp.KeyO,
		},
	}
	mk := d.AddMatrixKeyboard(colPins, rowPins, matrixKeyCodes)
	mk.SetCallback(func(layer, index int, state keyboard.State) {
		keyChan <- KeyEvent{layer: layer, indx: index, state: state}
	})

	go func() {
		for {
			select {
			case x := <-keyChan:
				_, err := tx.Write([]byte{
					uint8(x.layer), uint8(x.indx), uint8(x.state),
				})

				if err != nil {
					println("failed to send key:", err.Error())
				}
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
	for cont {
		err := d.Tick()
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
