package keyboard

const (
	LayerCount = 6

	// This value must also be set in `matrix.cols` in vial.json.
	// Keeping this value as a multiple of 14 is more efficient for vial communication.
	MaxKeyCount = 100
)
