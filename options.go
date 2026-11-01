package keyboard

// Options holds the settings that the Option functions change.
type Options struct {
	// InvertButtonState selects an active low scan. A pressed key reads Low.
	InvertButtonState bool
	// InvertDiode tells that the diode anode is on the row side.
	InvertDiode bool
}

// Option changes one field of Options.
type Option func(*Options)

// InvertButtonState selects an active low scan. The input pins use a pull up
// and a pressed key reads Low. The default scan is active high.
func InvertButtonState(b bool) Option {
	return func(o *Options) {
		o.InvertButtonState = b
	}
}

// InvertDiode tells that the diode anode is on the row side and the cathode
// is on the column side. By default the anode is on the column side.
func InvertDiode(b bool) Option {
	return func(o *Options) {
		o.InvertDiode = b
	}
}
