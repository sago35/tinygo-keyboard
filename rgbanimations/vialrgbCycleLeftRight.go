package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleLeftRightAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleLeftRight,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_LEFT_RIGHT,
	}
}

func cycleLeftRightMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	h := matrix.LedMatrixMapping[i].PhysicalX - uint8(time)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleLeftRight(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, cycleLeftRightMath)
}
