package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBandValAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandVal,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_VAL,
	}
}

func bandValMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	v16 := int16(matrix.CurrentValue) -
		Abs16(int16(Scale8(matrix.LedMatrixMapping[i].PhysicalX, 228))+28-int16(time))*8
	var v8 uint8
	if v16 < 0 {
		v8 = 0
	} else {
		v8 = Scale8(uint8(v16), matrix.CurrentValue)
	}
	return matrix.CurrentHue, matrix.CurrentSaturation, v8
}

func vialRGBBandVal(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, bandValMath)
}
