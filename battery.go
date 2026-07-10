//go:build tinygo

package keyboard

// BatteryPercent converts a LiPo cell voltage (V) to a remaining charge
// percentage (0-100), linearly interpolating over the 3.3V (empty) to 4.2V
// (full) usable range.
func BatteryPercent(vbat float64) uint8 {
	const (
		vEmpty = 3.3
		vFull  = 4.2
	)
	p := (vbat - vEmpty) / (vFull - vEmpty) * 100
	if p < 0 {
		p = 0
	} else if p > 100 {
		p = 100
	}
	return uint8(p)
}
