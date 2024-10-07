package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetRainbowMovingChevronAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBRainbowMovingChevron,
		AnimationType: keyboard.VIALRGB_EFFECT_RAINBOW_MOVING_CHEVRON,
	}
}

func rainbowMovingChevronMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	h := matrix.CurrentHue + Abs8(int8(matrix.LedMatrixMapping[i].PhysicalY-matrix.CenterYPhysical)+int8(matrix.CenterXPhysical-uint8(time)))
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBRainbowMovingChevron(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, rainbowMovingChevronMath)
}
