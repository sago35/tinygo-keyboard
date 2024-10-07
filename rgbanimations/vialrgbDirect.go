package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetDirectAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBDirectAnim,
		AnimationType: keyboard.VIALRGB_EFFECT_DIRECT,
	}
}

func vialRGBDirectAnim(matrix *keyboard.RGBMatrix) {
	for i, hsv := range matrix.LedMatrixDirectVals {
		r, g, b, a := HSVToRGB(hsv.H, hsv.S, hsv.V)
		matrix.LedMatrixVals[i].R = r
		matrix.LedMatrixVals[i].G = g
		matrix.LedMatrixVals[i].B = b
		matrix.LedMatrixVals[i].A = a
	}
}
