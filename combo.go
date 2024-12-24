//go:build tinygo

package keyboard

type Combo struct {
	Keys      [4]Keycode
	OutputKey Keycode
}

func (d *Device) SetCombo(index int, c Combo) {
	d.Combos[index][0] = c.Keys[0]
	d.Combos[index][1] = c.Keys[1]
	d.Combos[index][2] = c.Keys[2]
	d.Combos[index][3] = c.Keys[3]
	d.Combos[index][4] = c.OutputKey
}
