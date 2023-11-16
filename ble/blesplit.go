//go:build tinygo && nrf52840

package ble

import (
	k "machine/usb/hid/keyboard"

	"tinygo.org/x/bluetooth"
)

var tx = &bluetooth.Characteristic{}

type bleSplitKeyboard struct {
	keyboard
	Name      string
	report    [9]byte
	pressed   []k.Keycode
	connected bool
}

func NewSplitKeyboard(name string) *bleSplitKeyboard {
	return &bleSplitKeyboard{
		Name: name,
	}
}

func (k *bleSplitKeyboard) Connect() error {
	var err error

	name := k.Name
	if len(name) > 14 {
		name = name[:14]
	}
	adapter.SetConnectHandler(func(device bluetooth.Address, connected bool) {
		println("connected:", connected)
	})

	adv := adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: name,
	})
	err = adv.Start()
	if err != nil {
		return err
	}

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDNordicUART,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: tx,
				UUID:   bluetooth.CharacteristicUUIDUARTTX,
				Value:  k.report[:3],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})
	return nil
}

func (k *bleSplitKeyboard) Up(c k.Keycode) error {
	for i, p := range k.pressed {
		if c == p {
			k.pressed = append(k.pressed[:i], k.pressed[i+1:]...)
			row := byte(c >> 8)
			col := byte(c)
			_, err := tx.Write([]byte{0x55, byte(row), byte(col)})
			return err
		}
	}
	return nil
}

func (k *bleSplitKeyboard) Down(c k.Keycode) error {
	found := false
	for _, p := range k.pressed {
		if c == p {
			found = true
		}
	}
	if !found {
		k.pressed = append(k.pressed, c)
		row := byte(c >> 8)
		col := byte(c)
		_, err := tx.Write([]byte{0xAA, byte(row), byte(col)})
		return err
	}
	return nil
}
