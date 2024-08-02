package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetAlphasModsAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBAlphasModsAnim,
		AnimationType: keyboard.VIALRGB_EFFECT_ALPHAS_MODS,
	}
}

func vialRGBAlphasModsAnim(matrix *keyboard.RGBMatrix) {
	r1, g1, b1, a1 := HSVToRGB(matrix.CurrentHue, matrix.CurrentSaturation, matrix.CurrentValue)
	r2, g2, b2, a2 := HSVToRGB((matrix.CurrentHue+matrix.CurrentSpeed)%0xFF, matrix.CurrentSaturation, matrix.CurrentValue)
	for i := range matrix.LedMatrixVals {
		if matrix.LedMatrixMapping[i].LedFlags&keyboard.LED_FLAG_MODIFIER == keyboard.LED_FLAG_MODIFIER {
			matrix.LedMatrixVals[i].R = r2
			matrix.LedMatrixVals[i].G = g2
			matrix.LedMatrixVals[i].B = b2
			matrix.LedMatrixVals[i].A = a2
		} else {
			matrix.LedMatrixVals[i].R = r1
			matrix.LedMatrixVals[i].G = g1
			matrix.LedMatrixVals[i].B = b1
			matrix.LedMatrixVals[i].A = a1
		}
	}
}
