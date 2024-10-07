package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetRainbowPinwheelsAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBRainbowPinwheels,
		AnimationType: keyboard.VIALRGB_EFFECT_RAINBOW_PINWHEELS,
	}
}

func rainbowPinwheelsMath(matrix *keyboard.RGBMatrix, sin int8, cos int8, i int, time uint16) (uint8, uint8, uint8) {
	h := matrix.CurrentHue +
		uint8(
			int8(matrix.LedMatrixMapping[i].PhysicalY-matrix.CenterYPhysical)*3*cos+
				int8(56-Abs8(int8(matrix.LedMatrixMapping[i].PhysicalX-matrix.CenterXPhysical)))*2*sin,
		)/128
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBRainbowPinwheels(matrix *keyboard.RGBMatrix) {
	effectRunnerSinCosI(matrix, rainbowPinwheelsMath)
}
