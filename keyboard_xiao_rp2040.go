//go:build xiao_rp2040 || xiao_rp2040_sgkb

package keyboard

import (
	"machine"
)

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
