package rgbanimations

import keyboard "github.com/percyjw-2/tinygo-keyboard"

var time = uint16(0)

type effectRunnerIFunc func(matrix *keyboard.RGBMatrix, i int, time uint16) (uint8, uint8, uint8)
type effectRunnerDXDYFunc func(matrix *keyboard.RGBMatrix, dx int16, dy int16, time uint16) (uint8, uint8, uint8)
type effectRunnerDXDYDistFunc func(matrix *keyboard.RGBMatrix, dx int16, dy int16, dist uint8, time uint16) (uint8, uint8, uint8)
type effectRunnerSinCosIFunc func(matrix *keyboard.RGBMatrix, sin int8, cos int8, i int, time uint16) (uint8, uint8, uint8)

func effectRunnerI(matrix *keyboard.RGBMatrix, mathFunc effectRunnerIFunc) {
	timeScaled := Scale16by8(time, matrix.CurrentSpeed)
	for i, val := range matrix.LedMatrixVals {
		h, s, v := mathFunc(matrix, i, timeScaled)
		val.R, val.G, val.B, val.A = HSVToRGB(h, s, v)
	}
	time++
}

func effectRunnerDXDY(matrix *keyboard.RGBMatrix, mathFunc effectRunnerDXDYFunc) {
	timeScaled := Scale16by8(time, matrix.CurrentSpeed)
	for i, val := range matrix.LedMatrixVals {
		dx := int16(matrix.LedMatrixMapping[i].PhysicalX) - int16(matrix.CenterXPhysical)
		dy := int16(matrix.LedMatrixMapping[i].PhysicalY) - int16(matrix.CenterYPhysical)
		h, s, v := mathFunc(matrix, dx, dy, timeScaled)
		val.R, val.G, val.B, val.A = HSVToRGB(h, s, v)
	}
	time++
}

func effectRunnerDXDYDist(matrix *keyboard.RGBMatrix, mathFunc effectRunnerDXDYDistFunc) {
	timeScaled := Scale16by8(time, matrix.CurrentSpeed)
	for i, val := range matrix.LedMatrixVals {
		dx := int16(matrix.LedMatrixMapping[i].PhysicalX) - int16(matrix.CenterXPhysical)
		dy := int16(matrix.LedMatrixMapping[i].PhysicalY) - int16(matrix.CenterYPhysical)
		dist := Sqrt16(uint16(dx*dx + dy*dy))
		h, s, v := mathFunc(matrix, dx, dy, dist, timeScaled)
		val.R, val.G, val.B, val.A = HSVToRGB(h, s, v)
	}
	time++
}

func effectRunnerSinCosI(matrix *keyboard.RGBMatrix, mathFunc effectRunnerSinCosIFunc) {
	timeScaled := Scale16by8(time, matrix.CurrentSpeed)
	cosValue := int8(Cos8(uint8(time)) - 128)
	sinValue := int8(Sin8(uint8(time)) - 128)
	for i, val := range matrix.LedMatrixVals {
		h, s, v := mathFunc(matrix, sinValue, cosValue, i, timeScaled)
		val.R, val.G, val.B, val.A = HSVToRGB(h, s, v)
	}
	time++
}
