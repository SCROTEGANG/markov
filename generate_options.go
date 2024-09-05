package markov

type generateOptions struct {
	MaxTries int
	PRNG     func() int
	Filter   func(*Result) bool
}

// GenerateOption is used to configure a Markov chain's generation settings.
type GenerateOption func(*generateOptions)

// WithMaxTries sets the allotted number of attempts for a generation.
func WithMaxTries(mt int) GenerateOption {
	return func(g *generateOptions) {
		g.MaxTries = mt
	}
}

// WithPRNG sets the pseudo-random number generator used for generation.
func WithPRNG(fn func() int) GenerateOption {
	return func(g *generateOptions) {
		g.PRNG = fn
	}
}

// WithFilter sets the filter used for generation.
func WithFilter(fn func(*Result) bool) GenerateOption {
	return func(g *generateOptions) {
		g.Filter = fn
	}
}
