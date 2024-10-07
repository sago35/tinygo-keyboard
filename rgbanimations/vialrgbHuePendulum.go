package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetHuePendulumAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBHuePendulum,
		AnimationType: keyboard.VIALRGB_EFFECT_HUE_PENDULUM,
	}
}

func huePendulumMath(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8) {
	huedelta := uint8(12)
	h := matrix.CurrentHue + Scale8(Abs8(int8(Sin8(uint8(time))+(matrix.LedMatrixMapping[i].PhysicalX)-128))*2, huedelta)
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBHuePendulum(matrix *keyboard.RGBMatrix) {
	effectRunnerI(matrix, huePendulumMath)
}
