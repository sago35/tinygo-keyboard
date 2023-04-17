package keyboard

type Options struct {
	InvertButtonState bool
}

type Option func(*Options)

func InvertButtonState(b bool) Option {
	return func(o *Options) {
		o.InvertButtonState = b
	}
}
