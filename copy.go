package markov

func copyFragment(data []*Fragment) []*Fragment {
	out := make([]*Fragment, 0, len(data))

	for _, v := range data {
		out = append(out, &Fragment{
			Words: v.Words,
			Refs:  v.Refs,
		})
	}

	return out
}

func copyCorpus(data Corpus) Corpus {
	out := make(Corpus)

	for k, v := range data {
		out[k] = copyFragment(v)
	}

	return out
}
