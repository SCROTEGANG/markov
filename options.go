package markov

// Option is used to configure a Markov chain.
type Option func(*Markov)

// WithStateSize sets the size of a Markov chain's state.
func WithStateSize(ss int) Option {
	return func(m *Markov) {
		m.stateSize = ss
	}
}

func WithStore(s Store) Option {
	return func(m *Markov) {
		m.store = s
	}
}
