//go:build xiao_rp2040 || xiao_rp2040_sgkb

package keyboard

import (
	"context"
	"machine"
	k "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/sago35/tinygo-keyboard/keycodes"
)

func (d *Device) LoopUartRx(ctx context.Context) error {
	//buf := make([]byte, 0, 3)
	uart := machine.UART0
	for uart.Buffered() > 0 {
		uart.ReadByte()
	}
	cont := true
	for cont {
		select {
		case <-ctx.Done():
			cont = false
			continue
		default:
		}

		pressToRelease := []Keycode{}

		// read from key matrix
		for _, k := range d.dmk {
			k.Get()
			for row := range k.State {
				for col := range k.State[row] {
					switch k.State[row][col] {
					case None:
						// skip
					case NoneToPress:
						x := k.Keys[d.layer][row][col]
						found := false
						for _, p := range d.pressed {
							if x == p {
								found = true
							}
						}
						if !found {
							d.pressed = append(d.pressed, x)
						}

					case Press:
					case PressToRelease:
						x := k.Keys[d.layer][row][col]

						for i, p := range d.pressed {
							if x == p {
								d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
								pressToRelease = append(pressToRelease, x)
							}
						}
					}
				}
			}
		}

		// // read from uart
		// for uart.Buffered() > 0 {
		// 	data, _ := uart.ReadByte()
		// 	buf = append(buf, data)

		// 	if len(buf) == 3 {
		// 		row, col := buf[1], buf[2]
		// 		current := false
		// 		switch buf[0] {
		// 		case 0xAA: // press
		// 			current = true
		// 		case 0x55: // release
		// 			current = false
		// 		default:
		// 			buf[0], buf[1] = buf[1], buf[2]
		// 			buf = buf[:2]
		// 			continue
		// 		}

		// 		switch d.State2[row][col] {
		// 		case None:
		// 			if current {
		// 				d.State2[row][col] = NoneToPress
		// 			} else {
		// 			}
		// 		case NoneToPress:
		// 			if current {
		// 				d.State2[row][col] = Press
		// 			} else {
		// 				d.State2[row][col] = PressToRelease
		// 			}
		// 		case Press:
		// 			if current {
		// 			} else {
		// 				d.State2[row][col] = PressToRelease
		// 			}
		// 		case PressToRelease:
		// 			if current {
		// 				d.State2[row][col] = NoneToPress
		// 			} else {
		// 				d.State2[row][col] = None
		// 			}
		// 		}

		// 		switch d.State2[row][col] {
		// 		case None:
		// 			// skip
		// 		case NoneToPress:
		// 			x := d.Keys2[d.layer][row][col]
		// 			found := false
		// 			for _, p := range d.pressed {
		// 				if x == p {
		// 					found = true
		// 				}
		// 			}
		// 			if !found {
		// 				d.pressed = append(d.pressed, x)
		// 			}
		// 		case Press:
		// 		case PressToRelease:
		// 			x := d.Keys2[d.layer][row][col]
		// 			for i, p := range d.pressed {
		// 				if x == p {
		// 					d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
		// 					pressToRelease = append(pressToRelease, x)
		// 				}
		// 			}
		// 		}
		// 		buf = buf[:0]
		// 	}
		// }

		for i, x := range d.pressed {
			if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
				d.layer = int(x) & 0x0F
				if d.modKeyCallback != nil {
					d.modKeyCallback(d.layer, true)
				}
			} else if x&0xF000 == 0xD000 {
				switch x & 0x00FF {
				case 0x01, 0x02, 0x03:
					d.Mouse.Press(mouse.Button(x & 0x00FF))
				case 0x04:
					d.Mouse.WheelDown()
					// ここ上手にキーリピートさせたい感じはある
					d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
					pressToRelease = append(pressToRelease, x)
				case 0x05:
					d.Mouse.WheelUp()
					// ここ上手にキーリピートさせたい感じはある
					d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
					pressToRelease = append(pressToRelease, x)
				}
			} else {
				d.Keyboard.Down(k.Keycode(x))
			}
		}

		for _, x := range pressToRelease {
			if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
				if d.modKeyCallback != nil {
					d.modKeyCallback(d.layer, false)
				}
				d.layer = 0

				pressed := []Keycode{}
				for _, p := range d.pressed {
					if p&0xF000 == 0xD000 {
						switch p & 0x00FF {
						case 0x01, 0x02, 0x03:
							d.Mouse.Release(mouse.Button(p & 0x00FF))
						case 0x04:
							//d.Mouse.WheelDown()
						case 0x05:
							//d.Mouse.WheelUp()
						}
					} else {
						switch k.Keycode(p) {
						case keycodes.KeyLeftCtrl, keycodes.KeyRightCtrl:
							pressed = append(pressed, p)
						default:
							d.Keyboard.Up(k.Keycode(p))
						}
					}
				}
				d.pressed = d.pressed[:0]
				d.pressed = append(d.pressed, pressed...)

			} else if x&0xF000 == 0xD000 {
				switch x & 0x00FF {
				case 0x01, 0x02, 0x03:
					d.Mouse.Release(mouse.Button(x & 0x00FF))
				case 0x04:
					//d.Mouse.WheelDown()
				case 0x05:
					//d.Mouse.WheelUp()
				}
			} else {
				d.Keyboard.Up(k.Keycode(x))
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

//func (d *Device) LoopUartTx(ctx context.Context) error {
//	cont := true
//	for cont {
//		select {
//		case <-ctx.Done():
//			cont = false
//			continue
//		default:
//		}
//
//		d.Get()
//
//		for row := range d.State {
//			for col := range d.State[row] {
//				switch d.State[row][col] {
//				case None:
//					// skip
//				case NoneToPress:
//					x := d.Keys[d.layer][row][col]
//					found := false
//					for _, p := range d.pressed {
//						if x == p {
//							found = true
//						}
//					}
//					if !found {
//						d.pressed = append(d.pressed, x)
//					}
//
//					if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
//						d.layer = int(x) & 0x0F
//						if d.modKeyCallback != nil {
//							d.modKeyCallback(d.layer, true)
//						}
//					} else if x&0xF000 == 0xD000 {
//						switch x & 0x00FF {
//						case 0x01, 0x02, 0x03:
//							d.Mouse.Press(mouse.Button(x & 0x00FF))
//						case 0x04:
//							d.Mouse.WheelDown()
//						case 0x05:
//							d.Mouse.WheelUp()
//						}
//					} else {
//						//d.Keyboard.Down(k.Keycode(x))
//						Down(row, col)
//					}
//					if d.Debug {
//						fmt.Printf("%2d %2d %04X down\r\n", row, col, d.Keys[0][row][col])
//					}
//				case Press:
//				case PressToRelease:
//					x := d.Keys[d.layer][row][col]
//
//					for i, p := range d.pressed {
//						if x == p {
//							d.pressed = append(d.pressed[:i], d.pressed[i+1:]...)
//						}
//					}
//
//					if x&keycodes.ModKeyMask == keycodes.ModKeyMask {
//						if d.modKeyCallback != nil {
//							d.modKeyCallback(d.layer, false)
//						}
//						d.layer = 0
//
//						for _, p := range d.pressed {
//							if p&0xF000 == 0xD000 {
//								switch p & 0x00FF {
//								case 0x01, 0x02, 0x03:
//									d.Mouse.Release(mouse.Button(p & 0x00FF))
//								case 0x04:
//									//d.Mouse.WheelDown()
//								case 0x05:
//									//d.Mouse.WheelUp()
//								}
//							} else {
//								Up(row, col)
//							}
//						}
//						d.pressed = d.pressed[:]
//
//					} else if x&0xF000 == 0xD000 {
//						switch x & 0x00FF {
//						case 0x01, 0x02, 0x03:
//							d.Mouse.Release(mouse.Button(x & 0x00FF))
//						case 0x04:
//							//d.Mouse.WheelDown()
//						case 0x05:
//							//d.Mouse.WheelUp()
//						}
//					} else {
//						Up(row, col)
//					}
//					if d.Debug {
//						fmt.Printf("%2d %2d %04X up\r\n", row, col, d.Keys[0][row][col])
//					}
//				}
//			}
//		}
//
//		time.Sleep(5 * time.Millisecond)
//	}
//
//	return nil
//}

func Down(row, col int) {
	machine.UART0.Write([]byte{0xAA, byte(row), byte(col)})
}

func Up(row, col int) {
	machine.UART0.Write([]byte{0x55, byte(row), byte(col)})
}
