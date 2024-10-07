package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

func GetDualBeaconAnim() keyboard.RgbAnimation {
	return keyboard.RgbAnimation{
		AnimationFunc: vialRGBDualBeacon,
		AnimationType: keyboard.VIALRGB_EFFECT_DUAL_BEACON,
	}
}

func dualBeaconMath(matrix *keyboard.RGBMatrix, sin int8, cos int8, i int, _ uint16) (uint8, uint8, uint8) {
	h := matrix.CurrentHue + uint8(int8(matrix.LedMatrixMapping[i].PhysicalY-matrix.CenterYPhysical)*cos+int8(matrix.LedMatrixMapping[i].PhysicalX-matrix.CenterXPhysical)*sin)/128
	return h, matrix.CurrentSaturation, matrix.CurrentValue
}

func vialRGBDualBeacon(matrix *keyboard.RGBMatrix) {
	effectRunnerSinCosI(matrix, dualBeaconMath)
}
