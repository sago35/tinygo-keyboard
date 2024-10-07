package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBansSpiralValAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandSpiralVal,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_SPIRAL_VAL,
	}
}

func bandSpiralValMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, dist uint8, time uint16) (uint8, uint8, uint8) {
	v := Scale8(matrix.CurrentValue+dist-uint8(time)-Atan28(dy, dx), matrix.CurrentValue)
	return matrix.CurrentHue, matrix.CurrentSaturation, v
}

func vialRGBBandSpiralVal(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDYDist(matrix, bandSpiralValMath)
}
