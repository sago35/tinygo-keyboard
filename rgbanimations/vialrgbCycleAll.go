package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleAllAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleAll,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_ALL,
	}
}

func cycleAllMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	return uint8(time), matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleAll(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, cycleAllMath)
}
