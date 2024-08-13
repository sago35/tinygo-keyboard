package keyboard

import (
	"image/color"
	"time"
	"tinygo.org/x/drivers/ws2812"
)

type RGBMatrix struct {
	maximumBrightness   uint8
	ledCount            uint16
	LedMatrixMapping    []LedMatrixPosition
	implementedEffects  []RgbAnimation
	currentEffect       RgbAnimation
	CurrentSpeed        uint8
	CurrentHue          uint8
	CurrentSaturation   uint8
	CurrentValue        uint8
	LedMatrixDirectVals []LedMatrixDirectModeColor
	LedMatrixVals       []color.RGBA
	ledDriver           *ws2812.Device
}

type LedMatrixPosition struct {
	PhysicalX   uint8
	PhysicalY   uint8
	KbIndex     uint8
	MatrixIndex uint8
	LedFlags    uint8
}

type LedMatrixDirectModeColor struct {
	H uint8
	S uint8
	V uint8
}

type RgbAnimationFunc func(matrix *RGBMatrix)

type RgbAnimation struct {
	AnimationFunc RgbAnimationFunc
	AnimationType uint16
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

func (d *Device) AddRGBMatrix(brightness uint8, ledCount uint16, ledMatrixMapping []LedMatrixPosition, animations []RgbAnimation, ledDriver *ws2812.Device) {
	if int(ledCount) != len(ledMatrixMapping) {
		panic("LedMatrixMapping must have length equal to number of LedMatrixMapping")
	}
	effectOffAnimation := RgbAnimation{
		AnimationFunc: func(matrix *RGBMatrix) {
			matrix.ledDriver.WriteColors(matrix.LedMatrixVals)
		},
		AnimationType: VIALRGB_EFFECT_OFF,
	}
	rgbMatrix := RGBMatrix{
		maximumBrightness: brightness,
		ledCount:          ledCount,
		LedMatrixMapping:  ledMatrixMapping,
		implementedEffects: []RgbAnimation{
			effectOffAnimation,
		},
		currentEffect:     effectOffAnimation,
		CurrentSpeed:      0x00,
		CurrentHue:        0xFF,
		CurrentSaturation: 0xFF,
		CurrentValue:      brightness,
		LedMatrixVals:     make([]color.RGBA, ledCount),
		ledDriver:         ledDriver,
	}
	rgbMatrix.implementedEffects = append(rgbMatrix.implementedEffects, animations...)
	for _, animation := range animations {
		if animation.AnimationType == VIALRGB_EFFECT_DIRECT {
			rgbMatrix.LedMatrixDirectVals = make([]LedMatrixDirectModeColor, ledCount)
		}
	}
	d.rgbMat = &rgbMatrix
}

func (d *Device) GetRGBMatrixMaximumBrightness() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.maximumBrightness
}

func (d *Device) GetRGBMatrixLEDCount() uint16 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.ledCount
}

func (d *Device) GetRGBMatrixLEDMapping(ledIndex uint16) LedMatrixPosition {
	invalidPosition := LedMatrixPosition{
		KbIndex:     0xFF,
		MatrixIndex: 0xFF,
	}
	if !d.IsRGBMatrixEnabled() || ledIndex >= d.rgbMat.ledCount {
		return invalidPosition
	}
	return d.rgbMat.LedMatrixMapping[ledIndex]
}

func (d *Device) GetSupportedRGBModes() []RgbAnimation {
	if !d.IsRGBMatrixEnabled() {
		return []RgbAnimation{}
	}
	return d.rgbMat.implementedEffects
}

func (d *Device) GetCurrentRGBMode() uint16 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.currentEffect.AnimationType
}

func (d *Device) SetCurrentRGBMode(mode uint16) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	for _, e := range d.rgbMat.implementedEffects {
		if e.AnimationType == mode {
			d.rgbMat.currentEffect = e
		}
	}
}

func (d *Device) GetCurrentSpeed() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.CurrentSpeed
}

func (d *Device) SetCurrentSpeed(speed uint8) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	d.rgbMat.CurrentSpeed = speed
}

func (d *Device) GetCurrentHue() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.CurrentHue
}

func (d *Device) GetCurrentSaturation() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.CurrentSaturation
}

func (d *Device) GetCurrentValue() uint8 {
	if !d.IsRGBMatrixEnabled() {
		return 0
	}
	return d.rgbMat.CurrentValue
}

func (d *Device) SetCurrentHSV(hue uint8, saturation uint8, value uint8) {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	d.rgbMat.CurrentHue = hue
	d.rgbMat.CurrentSaturation = saturation
	d.rgbMat.CurrentValue = value
}

func (d *Device) SetDirectHSV(hue uint8, saturation uint8, value uint8, ledIndex uint16) {
	if !d.IsDirectModeEnabled() {
		return
	}
	rgb := d.rgbMat
	var actualValue uint8
	if value > rgb.maximumBrightness {
		actualValue = rgb.maximumBrightness
	} else {
		actualValue = value
	}
	rgb.LedMatrixDirectVals[ledIndex] = LedMatrixDirectModeColor{
		H: hue,
		S: saturation,
		V: actualValue,
	}
}

func (d *Device) updateRGBTask() {
	if !d.IsRGBMatrixEnabled() {
		return
	}
	rgb := d.rgbMat
	for {
		time.Sleep(time.Millisecond * time.Duration(0x100-uint16(rgb.CurrentSpeed)))
		rgb.currentEffect.AnimationFunc(rgb)
		_ = rgb.ledDriver.WriteColors(rgb.LedMatrixVals)
	}
}

func (d *Device) IsRGBMatrixEnabled() bool {
	return d.rgbMat != nil
}

func (d *Device) IsDirectModeEnabled() bool {
	return d.IsRGBMatrixEnabled() && d.rgbMat.LedMatrixDirectVals != nil
}
