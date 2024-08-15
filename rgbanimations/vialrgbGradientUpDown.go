package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetGradientUpDownAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBGradientUpDown,
		AnimationType: keyboard.VIALRGB_EFFECT_GRADIENT_UP_DOWN,
	}
}

func vialRGBGradientUpDown(matrix *keyboard.RGBMatrix) {
	scale := Scale8(64, matrix.CurrentSpeed)
	for i, position := range matrix.LedMatrixMapping {
		h := uint16(matrix.CurrentHue) + uint16(scale*(position.PhysicalY>>4))
		r, g, b, a := HSVToRGB(uint8(h&0xFF), matrix.CurrentSaturation, matrix.CurrentValue)
		matrix.LedMatrixVals[i].R = r
		matrix.LedMatrixVals[i].G = g
		matrix.LedMatrixVals[i].B = b
		matrix.LedMatrixVals[i].A = a
	}
}
