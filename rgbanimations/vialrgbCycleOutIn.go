package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleOutInAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleOutIn,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_OUT_IN,
	}
}

func cycleOutInMath(matrix *keyboard.RGBMatrix, _ int16, _ int16, dist uint8, time uint16) (uint8, uint8, uint8) {
	h := 3*dist/2 + uint8(time)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleOutIn(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDYDist(matrix, cycleOutInMath)
}
