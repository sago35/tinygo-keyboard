package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleSpiralAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleSpiral,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_SPIRAL,
	}
}

func cycleSpiralMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, dist uint8, time uint16) (uint8, uint8, uint8) {
	h := dist - uint8(time) - Atan28(dy, dx)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleSpiral(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDYDist(matrix, cycleSpiralMath)
}
