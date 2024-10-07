package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetSolidColorAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialrgbSolidColorAnim,
		AnimationType: keyboard.VIALRGB_EFFECT_SOLID_COLOR,
	}
}

func vialrgbSolidColorAnim(matrix *keyboard.RGBMatrix) {
	r, g, b, a := HSVToRGB(matrix.CurrentHue, matrix.CurrentSaturation, matrix.CurrentValue)
	for i := range matrix.LedMatrixVals {
		matrix.LedMatrixVals[i].R = r
		matrix.LedMatrixVals[i].G = g
		matrix.LedMatrixVals[i].B = b
		matrix.LedMatrixVals[i].A = a
	}
}
