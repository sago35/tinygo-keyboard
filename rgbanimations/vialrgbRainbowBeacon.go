package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetRainbowBeaconAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBRainbowBeacon,
		AnimationType: keyboard.VIALRGB_EFFECT_RAINBOW_BEACON,
	}
}

func rainbowBeaconMath(matrix *keyboard.RGBMatrix, sin int8, cos int8, i int, _ uint16) (uint8, uint8, uint8) {
	h := matrix.CurrentHue + uint8(int8(matrix.LedMatrixMapping[i].PhysicalY-matrix.CenterYPhysical)*2*cos+int8(matrix.LedMatrixMapping[i].PhysicalX-matrix.CenterXPhysical)*2*sin)/128
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBRainbowBeacon(matrix *keyboard.RGBMatrix) {
	effectRunnerSinCosI(matrix, rainbowBeaconMath)
}
