//go:build tinygo && nrf52840

package keyboard

import (
	"tinygo.org/x/bluetooth"
)

type BleSplitKeyboard struct {
	State    []State
	Keys     [][]Keycode
	callback Callback

	connectTo string
	adapter   *bluetooth.Adapter
	ringbuf   *RingBuffer[bleKeyEvent]
	buf       []byte
	processed []int
}

type bleKeyEvent struct {
	index  int
	isHigh bool
}

func (d *Device) AddBleSplitKeyboard(size int, adapter *bluetooth.Adapter, connectTo string, keys [][]Keycode, opt ...Option) *BleSplitKeyboard {
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

	var vb [32]bleKeyEvent
	k := &BleSplitKeyboard{
		State:    state,
		Keys:     keydef,
		callback: func(layer, index int, state State) {},

		adapter:   adapter,
		connectTo: connectTo,
		ringbuf:   NewRingBuffer(vb[:]),
		buf:       make([]byte, 3),
		processed: make([]int, 0, 8),
	}

	d.kb = append(d.kb, k)
	return k
}

func (d *BleSplitKeyboard) SetCallback(fn Callback) {
	d.callback = fn
}

func (d *BleSplitKeyboard) Callback(layer, index int, state State) {
	if d.callback != nil {
		d.callback(layer, index, state)
	}
}

func (d *BleSplitKeyboard) Get() []State {
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
		index := b.index
		current := b.isHigh

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

func (d *BleSplitKeyboard) Key(layer, index int) Keycode {
	if layer >= LayerCount {
		return 0
	}
	if index >= len(d.Keys[layer]) {
		return 0
	}
	return d.Keys[layer][index]
}

func (d *BleSplitKeyboard) SetKeycode(layer, index int, key Keycode) {
	if layer >= LayerCount {
		return
	}
	if index >= len(d.Keys[layer]) {
		return
	}
	d.Keys[layer][index] = key
}

func (d *BleSplitKeyboard) GetKeyCount() int {
	return len(d.State)
}

func (d *BleSplitKeyboard) Init() error {
	d.adapter.SetConnectHandler(func(device bluetooth.Address, connected bool) {
		println("connected:", connected)
	})

	return d.connectToPeriph()
}

func (d *BleSplitKeyboard) connectToPeriph() error {
	var foundDevice bluetooth.ScanResult
	name := d.connectTo
	if len(name) > 14 {
		name = name[:14]
	}

	err := d.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		println(result.Address.String(), result.LocalName())
		if result.LocalName() != name {
			return
		}
		foundDevice = result

		// Stop the scan.
		err := d.adapter.StopScan()
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	println("connected")
	device, err := d.adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
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

	rx := chars[0]
	rx.EnableNotifications(
		func(buf []byte) {
			//println("received buf len:", len(buf), ":", buf[0], buf[1], buf[2])
			current := false
			if buf[0] == 0xAA {
				current = true
			}
			index := (int(buf[1]) << 8) + int(buf[2])

			d.ringbuf.Put(bleKeyEvent{
				index:  index,
				isHigh: current,
			})
		},
	)

	return nil
}
