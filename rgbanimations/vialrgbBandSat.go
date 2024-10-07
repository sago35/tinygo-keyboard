package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBandSatAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandSat,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_SAT,
	}
}

func bandSatMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	s16 := int16(matrix.CurrentSaturation) -
		Abs16(int16(Scale8(matrix.LedMatrixMapping[i].PhysicalX, 228))+28-int16(time))*8
	var s8 uint8
	if s16 < 0 {
		s8 = 0
	} else {
		s8 = Scale8(uint8(s16), matrix.CurrentSaturation)
	}
	return matrix.CurrentHue, s8, matrix.CurrentValue
}

func vialRGBBandSat(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, bandSatMath)
}
