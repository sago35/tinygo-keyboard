package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBandPinwheelValAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBandPinwheelVal,
		AnimationType: keyboard.VIALRGB_EFFECT_BAND_PINWHEEL_VAL,
	}
}

func bandPinwheelValMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, time uint16) (uint8, uint8, uint8) {
	v := Scale8(matrix.CurrentValue-uint8(time)-Atan28(dy, dx), matrix.CurrentValue)
	return matrix.CurrentHue, matrix.CurrentSaturation, v
}

func vialRGBBandPinwheelVal(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDY(matrix, bandPinwheelValMath)
}
