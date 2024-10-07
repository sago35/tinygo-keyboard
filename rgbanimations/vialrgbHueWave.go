package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetHueWave() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBHueWave,
		AnimationType: keyboard.VIALRGB_EFFECT_HUE_WAVE,
	}
}

func hueWaveMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	huedelta := uint8(24)
	h := matrix.CurrentHue + Scale8(Abs8(int8(matrix.LedMatrixMapping[i].PhysicalX-uint8(time))), huedelta)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBHueWave(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, hueWaveMath)
}
