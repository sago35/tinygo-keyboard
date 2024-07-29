package keyboard

type RGBMatrix struct {
	maximumBrightness uint8
	ledCount          uint16
	ledMatrixMapping  []LedMatrixPosition
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

func (d *Device) AddRGBMatrix(brightness uint8, ledCount uint16, ledMatrixMapping []LedMatrixPosition) {
	if int(ledCount) != len(ledMatrixMapping) {
		panic("ledMatrixMapping must have length equal to number of ledMatrixMapping")
	}
	rgbMatrix := RGBMatrix{
		maximumBrightness: brightness,
		ledCount:          ledCount,
		ledMatrixMapping:  ledMatrixMapping,
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

func (d *Device) IsRGBMatrixEnabled() bool {
	return len(d.rgbMat) > 0
}
