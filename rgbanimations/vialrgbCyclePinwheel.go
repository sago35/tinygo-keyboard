package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetCyclePinwheelAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBCyclePinwheel,
		AnimationType: keyboard.VIALRGB_EFFECT_CYCLE_PINWHEEL,
	}
}

func cyclePinwheelMath(matrix *keyboard.RGBMatrix, dx int16, dy int16, time uint16) (uint8, uint8, uint8) {
	h := Atan28(dy, dx) + uint8(time)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBCyclePinwheel(matrix *keyboard.RGBMatrix) {
	effectRunnerDXDY(matrix, cyclePinwheelMath)
}
