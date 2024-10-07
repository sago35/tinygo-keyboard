package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBansSpiralSatAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandSpiralSat,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_SPIRAL_SAT,
	}
}

func bandSpiralSatMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, dist uint8, time uint16) (uint8, uint8, uint8) {
	s := Scale8(matrix.CurrentSaturation+dist-uint8(time)-Atan28(dy, dx), matrix.CurrentSaturation)
	return matrix.CurrentHue, s, matrix.CurrentValue
}

func vialRGBBandSpiralSat(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDYDist(matrix, bandSpiralSatMath)
}
