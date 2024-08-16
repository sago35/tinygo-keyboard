package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetBreathingAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBBreathing,
		AnimationType: keyboard.VIALRGB_EFFECT_BREATHING,
	}
}

var time = uint16(0)

func vialRGBBreathing(matrix *keyboard.RGBMatrix) {
	v := Scale8(Abs8(int8(Sin8(uint8(time))-128))*8, matrix.CurrentValue)
	r, g, b, a := HSVToRGB(matrix.CurrentHue, matrix.CurrentSaturation, v)
	for _, val := range matrix.LedMatrixVals {
		val.R = r
		val.G = g
		val.B = b
		val.A = a
	}
	time++
}
