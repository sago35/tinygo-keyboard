//go:build tinygo && nrf52840

package keyboard

import (
	k "machine/usb/hid/keyboard"

	"tinygo.org/x/bluetooth"
)

type BleTxKeyboard struct {
	RxBleName string

	adapter          *bluetooth.Adapter
	pressed          []k.Keycode
	txCharecteristic bluetooth.DeviceCharacteristic
}

func (k *BleTxKeyboard) Up(c k.Keycode) error {
	for i, p := range k.pressed {
		if c == p {
			k.pressed = append(k.pressed[:i], k.pressed[i+1:]...)
			row := byte(c >> 8)
			col := byte(c)
			_, err := k.txCharecteristic.WriteWithoutResponse([]byte{0x55, byte(row), byte(col)})
			return err
		}
	}
	return nil
}

func (k *BleTxKeyboard) Down(c k.Keycode) error {
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
		_, err := k.txCharecteristic.WriteWithoutResponse([]byte{0xAA, byte(row), byte(col)})
		return err
	}
	return nil
}

func (k *BleTxKeyboard) Connect() error {
	k.adapter = bluetooth.DefaultAdapter
	err := k.adapter.Enable()
	if err != nil {
		return err
	}
	var foundDevice bluetooth.ScanResult
	err = k.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() != k.RxBleName {
			return
		}
		foundDevice = result
		err = adapter.StopScan()
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	device, err := k.adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDNordicUART})
	if err != nil {
		return err
	}
	service := services[0]
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{bluetooth.CharacteristicUUIDUARTTX})
	if err != nil {
		return err
	}
	k.txCharecteristic = chars[0]
	return nil
}

type bleKeyEvent struct {
	Index  int
	IsHigh bool
}

type BleKeyboard struct {
	RxBleName string
	State     []State
	Keys      [][]Keycode
	callback  Callback
	ringbuf   *RingBuffer[bleKeyEvent]
	processed []int

	adapter *bluetooth.Adapter
	buf     []byte
	changed bool
}

func (d *Device) AddBleKeyboard(size int, rxname string, keys [][]Keycode) *BleKeyboard {
	state := make([]State, size)

	keydef := make([][]Keycode, LayerCount)
	for l := 0; l < len(keydef); l++ {
		keydef[l] = make([]Keycode, len(state))
	}
	for l := 0; l < len(keys); l++ {
		for kc := 0; kc < len(keys[l]); kc++ {
			keydef[l][kc] = keys[l][kc]
		}
	}

	var rb [32]bleKeyEvent
	k := &BleKeyboard{
		RxBleName: rxname,
		State:     state,
		Keys:      keydef,
		adapter:   bluetooth.DefaultAdapter,
		callback:  func(layer, index int, state State) {},
		buf:       make([]byte, 3),
		changed:   false,
		ringbuf:   NewRingBuffer(rb[:]),
		processed: make([]int, 0, 8),
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *BleKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *BleKeyboard) Get() []State {

	for i := range d.State {
		switch d.State[i] {
		case NoneToPress:
			d.State[i] = Press
		case PressToRelease:
			d.State[i] = None
		}
	}

	d.processed = d.processed[:0]

	cont := true
	for cont {
		b, ok := d.ringbuf.Peek()
		if !ok {
			return d.State
		}
		index := b.Index
		current := b.IsHigh

		for _, idx := range d.processed {
			if index == idx {
				return d.State
			}
		}
		d.processed = append(d.processed, index)

		d.ringbuf.Get()

		switch d.State[index] {
		case None:
			if current {
				d.State[index] = NoneToPress
				d.callback(0, index, Press)
			} else {
			}
		case NoneToPress:
			if current {
				d.State[index] = Press
			} else {
				d.State[index] = PressToRelease
				d.callback(0, index, PressToRelease)
			}
		case Press:
			if current {
			} else {
				d.State[index] = PressToRelease
				d.callback(0, index, PressToRelease)
			}
		case PressToRelease:
			if current {
				d.State[index] = NoneToPress
				d.callback(0, index, Press)
			} else {
				d.State[index] = None
			}
		}

	}

	return d.State
}

func (d *BleKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *BleKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *BleKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *BleKeyboard) Init() error {
	err := d.adapter.Enable()
	if err != nil {
		return err
	}
	adv := d.adapter.DefaultAdvertisement()
	err = adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: d.RxBleName,
	})
	err = adv.Start()
	if err != nil {
		return err
	}

	d.adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDNordicUART,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDUARTTX,
				Value: d.buf[:],
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 3 {
						return
					}
					index := (int(value[1]) << 8) + int(value[2])
					isHigh := false
					switch value[0] {
					case 0xAA: // press
						isHigh = true
					case 0x55: // release
						isHigh = false
					}
					d.ringbuf.Put(bleKeyEvent{
						Index:  index,
						IsHigh: isHigh,
					})
				},
			},
		},
	})
	return nil

}
