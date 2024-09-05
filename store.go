package markov

import "context"

type Store interface {
	StartWords(context.Context) []*Fragment
	SetStartWords(context.Context, []*Fragment)
	EndWords(context.Context) []*Fragment
	SetEndWords(context.Context, []*Fragment)
	Corpus(context.Context) Corpus
	SetCorpus(context.Context, Corpus)

	Import(context.Context, *Structure)
	Export(context.Context) *Structure
}

type memoryStore struct {
	startWords []*Fragment
	endWords   []*Fragment
	corpus     Corpus
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		startWords: make([]*Fragment, 0),
		endWords:   make([]*Fragment, 0),
		corpus:     make(Corpus),
	}
}

func (s *memoryStore) StartWords(_ context.Context) []*Fragment {
	return copyFragment(s.startWords)
}

func (s *memoryStore) SetStartWords(_ context.Context, f []*Fragment) {
	s.startWords = f
}

func (s *memoryStore) EndWords(_ context.Context) []*Fragment {
	return copyFragment(s.endWords)
}

func (s *memoryStore) SetEndWords(_ context.Context, f []*Fragment) {
	s.endWords = f
}

func (s *memoryStore) Corpus(_ context.Context) Corpus {
	return copyCorpus(s.corpus)
}

func (s *memoryStore) SetCorpus(_ context.Context, c Corpus) {
	s.corpus = c
}

func (s *memoryStore) Import(_ context.Context, data *Structure) {
	s.startWords = data.StartWords
	s.endWords = data.EndWords
	s.corpus = data.Corpus
}

func (s *memoryStore) Export(ctx context.Context) *Structure {
	return &Structure{
		StartWords: copyFragment(s.startWords),
		EndWords:   copyFragment(s.endWords),
		Corpus:     copyCorpus(s.corpus),
	}
}
