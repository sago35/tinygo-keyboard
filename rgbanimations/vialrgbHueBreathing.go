package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetHueBreathingAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBHueBreathing,
		AnimationType: keyboard.VIALRGB_EFFECT_HUE_BREATHING,
	}
}

func vialRGBHueBreathing(matrix *keyboard.RGBMatrix) {
	huedelta := uint8(12)
	h, s, v := matrix.CurrentHue, matrix.CurrentSaturation, matrix.CurrentValue
	timeScaled := Scale16by8(time, matrix.CurrentSpeed/8)
	h = h + Scale8(Abs8(int8(Sin8(uint8(timeScaled)))-127)*2, huedelta)
	r, g, b, a := HSVToRGB(h, s, v)
	for _, val := range matrix.LedMatrixVals {
		val.R = r
		val.G = g
		val.B = b
		val.A = a
	}
	time++
}
