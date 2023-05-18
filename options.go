package keyboard

type Options struct {
	InvertButtonState bool
	InvertDiode       bool
}

type Option func(*Options)

func InvertButtonState(b bool) Option {
	return func(o *Options) {
		o.InvertButtonState = b
	}
}

func InvertDiode(b bool) Option {
	return func(o *Options) {
		o.InvertDiode = b
	}
}
