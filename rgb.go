package keyboard

type RGBMatrix struct {
	maximumBrightness  uint8
	ledCount           uint16
	ledMatrixMapping   []LedMatrixPosition
	implementedEffects []uint16
	currentEffect      uint16
	currentSpeed       uint8
	currentHue         uint8
	currentSaturation  uint8
	currentValue       uint8
}

type LedMatrixPosition struct {
	physicalX   uint8
	physicalY   uint8
	kbIndex     uint8
	matrixIndex uint8
	ledFlags    uint8
}

const (
	// LED Flags
	LED_FLAG_NONE      uint8 = 0x00 // If this LED has no flags
	LED_FLAG_ALL       uint8 = 0xFF // if this LED has all flags
	LED_FLAG_MODIFIER  uint8 = 0x01 // if the Key for this LED is a modifier
	LED_FLAG_UNDERGLOW uint8 = 0x02 // if the LED is for underglow
	LED_FLAG_KEYLIGHT  uint8 = 0x04 // if the LED is for key backlight
)

const (
	// RGB Modes
	VIALRGB_EFFECT_OFF = iota
	VIALRGB_EFFECT_DIRECT
	VIALRGB_EFFECT_SOLID_COLOR
	VIALRGB_EFFECT_ALPHAS_MODS
	VIALRGB_EFFECT_GRADIENT_UP_DOWN
	VIALRGB_EFFECT_GRADIENT_LEFT_RIGHT
	VIALRGB_EFFECT_BREATHING
	VIALRGB_EFFECT_BAND_SAT
	VIALRGB_EFFECT_BAND_VAL
	VIALRGB_EFFECT_BAND_PINWHEEL_SAT
	VIALRGB_EFFECT_BAND_PINWHEEL_VAL
	VIALRGB_EFFECT_BAND_SPIRAL_SAT
	VIALRGB_EFFECT_BAND_SPIRAL_VAL
	VIALRGB_EFFECT_CYCLE_ALL
	VIALRGB_EFFECT_CYCLE_LEFT_RIGHT
	VIALRGB_EFFECT_CYCLE_UP_DOWN
	VIALRGB_EFFECT_RAINBOW_MOVING_CHEVRON
	VIALRGB_EFFECT_CYCLE_OUT_IN
	VIALRGB_EFFECT_CYCLE_OUT_IN_DUAL
	VIALRGB_EFFECT_CYCLE_PINWHEEL
	VIALRGB_EFFECT_CYCLE_SPIRAL
	VIALRGB_EFFECT_DUAL_BEACON
	VIALRGB_EFFECT_RAINBOW_BEACON
	VIALRGB_EFFECT_RAINBOW_PINWHEELS
	VIALRGB_EFFECT_RAINDROPS
	VIALRGB_EFFECT_JELLYBEAN_RAINDROPS
	VIALRGB_EFFECT_HUE_BREATHING
	VIALRGB_EFFECT_HUE_PENDULUM
	VIALRGB_EFFECT_HUE_WAVE
	VIALRGB_EFFECT_TYPING_HEATMAP
	VIALRGB_EFFECT_DIGITAL_RAIN
	VIALRGB_EFFECT_SOLID_REACTIVE_SIMPLE
	VIALRGB_EFFECT_SOLID_REACTIVE
	VIALRGB_EFFECT_SOLID_REACTIVE_WIDE
	VIALRGB_EFFECT_SOLID_REACTIVE_MULTIWIDE
	VIALRGB_EFFECT_SOLID_REACTIVE_CROSS
	VIALRGB_EFFECT_SOLID_REACTIVE_MULTICROSS
	VIALRGB_EFFECT_SOLID_REACTIVE_NEXUS
	VIALRGB_EFFECT_SOLID_REACTIVE_MULTINEXUS
	VIALRGB_EFFECT_SPLASH
	VIALRGB_EFFECT_MULTISPLASH
	VIALRGB_EFFECT_SOLID_SPLASH
	VIALRGB_EFFECT_SOLID_MULTISPLASH
	VIALRGB_EFFECT_PIXEL_RAIN
	VIALRGB_EFFECT_PIXEL_FRACTAL
)

func (d *Device) AddRGBMatrix(brightness uint8, ledCount uint16, ledMatrixMapping []LedMatrixPosition) {
	if int(ledCount) != len(ledMatrixMapping) {
		panic("ledMatrixMapping must have length equal to number of ledMatrixMapping")
	}
	rgbMatrix := RGBMatrix{
		maximumBrightness: brightness,
		ledCount:          ledCount,
		ledMatrixMapping:  ledMatrixMapping,
		implementedEffects: []uint16{
			VIALRGB_EFFECT_OFF,
		},
		currentEffect:     VIALRGB_EFFECT_OFF,
		currentSpeed:      0,
		currentHue:        0,
		currentSaturation: 0,
		currentValue:      0,
	}
	d.rgbMat = append(d.rgbMat, rgbMatrix)
}

func (d *Device) GetRGBMatrixMaximumBrightness() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].maximumBrightness
}

func (d *Device) GetRGBMatrixLEDCount() uint16 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[1].ledCount
}

func (d *Device) GetRGBMatrixLEDMapping(ledIndex uint16) LedMatrixPosition {
	invalidPosition := LedMatrixPosition{
		kbIndex:     0xFF,
		matrixIndex: 0xFF,
	}
	if !d.IsRGBMatrixEnabled() || ledIndex >= d.rgbMat[0].ledCount {
		return invalidPosition
	}
	return d.rgbMat[0].ledMatrixMapping[ledIndex]
}

func (d *Device) GetSupportedRGBModes() []uint16 {
	if !d.IsRGBMatrixEnabled() {
		return []uint16{}
	}
	return d.rgbMat[0].implementedEffects
}

func (d *Device) GetCurrentRGBMode() uint16 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].currentEffect
}

func (d *Device) SetCurrentRGBMode(mode uint16) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	d.rgbMat[0].currentEffect = mode
}

func (d *Device) GetCurrentSpeed() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].currentSpeed
}

func (d *Device) SetCurrentSpeed(speed uint8) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	d.rgbMat[0].currentSpeed = speed
}

func (d *Device) GetCurrentHue() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].currentHue
}

func (d *Device) GetCurrentSaturation() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].currentSaturation
}

func (d *Device) GetCurrentValue() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat[0].currentValue
}

func (d *Device) SetCurrentHSV(hue uint8, saturation uint8, value uint8) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	d.rgbMat[0].currentHue = hue
	d.rgbMat[0].currentSaturation = saturation
	d.rgbMat[0].currentValue = value
}

func (d *Device) IsRGBMatrixEnabled() bool {
	return len(d.rgbMat) > 0
}
