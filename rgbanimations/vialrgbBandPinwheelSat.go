package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBandPinwheelSatAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandPinwheelSat,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_PINWHEEL_SAT,
	}
}

func bandPinwheelSatMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, time uint16) (uint8, uint8, uint8) {
	s := Scale8(matrix.CurrentSaturation-uint8(time)-Atan28(dy, dx), matrix.CurrentSaturation)
	return matrix.CurrentHue, s, matrix.CurrentValue
}

func vialRGBBandPinwheelSat(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDY(matrix, bandPinwheelSatMath)
}
