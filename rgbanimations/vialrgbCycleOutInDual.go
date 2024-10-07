package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCycleOutInDualAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCycleOutInDual,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_OUT_IN_DUAL,
	}
}

func cycleOutInDualMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, time uint16) (uint8, uint8, uint8) {
	dx = int16(matrix.CenterXPhysical/2) - Abs16(dx)
	dist := Sqrt16(uint16(dx*dx) + uint16(dy*dy))
	h := 3*dist + uint8(time)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCycleOutInDual(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDY(matrix, cycleOutInDualMath)
}
