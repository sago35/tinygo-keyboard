package main

import (
	"machine"
)

type ADCDevice struct {
	adc      machine.ADC
	min      int
	max      int
	RawValue uint16
	Value    int16
	invert   bool
	Prev10   int16
}

func NewADCDevice(pin machine.Pin, min, max int, invert bool) *ADCDevice {
	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})
	return &ADCDevice{
		adc:    adc,
		min:    min,
		max:    max,
		invert: invert,
	}
}

func (a *ADCDevice) Get() int16 {
	a.RawValue = a.adc.Get()
	ave := a.RawValue

	ret := 32767 * (int(ave) - a.min) / (a.max - a.min)
	if ret < 0 {
		ret = 0
	}
	if 32767 < ret {
		ret = 32767
	}

	ret -= 0x4000
	if a.invert {
		ret *= -1
	}
	a.Value = int16(ret)
	return a.Value
}

var mapx = map[int]int16{
	0: 0,
	1: 0,
	2: 0,
	3: 10,
	4: 20,
	5: 30,
	6: 40,
	7: 50,
	8: 50,
}

func (a *ADCDevice) Get2() int16 {
	a.Get()

	ret := int(a.Value)
	ret += 0x4000
	ret >>= 11
	ret -= 8

	v := mapx[abs(ret)]
	if abs(ret) >= 7 {
		if v <= abs(a.Prev10) {
			v = abs(a.Prev10)
			if abs(a.Prev10) < v*3 {
				v += 2
			}
		}
	}

	if ret < 0 {
		v *= -1
	}

	a.Prev10 = v
	return v / 10
}

func abs[T int | int16](x T) T {
	if x < 0 {
		return -1 * x
	}
	return x
}
