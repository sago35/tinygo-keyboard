package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleUpDownAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleUpDown,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_UP_DOWN,
	}
}

func cycleUpDownMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	h := matrix.LedMatrixMapping[i].PhysicalY - uint8(time)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleUpDown(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, cycleUpDownMath)
}
