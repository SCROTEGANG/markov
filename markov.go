package markov

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

const (
	defaultStateSize = 2
	defaultMaxTries  = 10
)

// Result contains data pertaining to a generation.
type Result struct {
	String string
	Score  int
	Refs   []string
	Tries  int
}

// Fragment contains fragmented data used to form a sentence.
type Fragment struct {
	Words string
	Refs  []string
}

// Corpus contains data that makes up the Markov chain.
type Corpus map[string][]*Fragment

// Structure represents the internal structure of a Markov chain
// and is used to import & export data.
type Structure struct {
	Corpus     Corpus
	StartWords []*Fragment
	EndWords   []*Fragment
}

// Markov represents a Markov chain which can be trained
// and used to generate sentences.
type Markov struct {
	data      []string
	stateSize int

	store Store
}

// New creates a new Markov chain with the given options.
func New(opts ...Option) *Markov {
	m := &Markov{}

	for _, o := range opts {
		o(m)
	}

	if m.stateSize == 0 {
		m.stateSize = defaultStateSize
	}

	if m.store == nil {
		m.store = newMemoryStore()
	}

	return m
}

// AddData adds the given data to the corpus.
func (m *Markov) AddData(ctx context.Context, data []string) {
	m.buildCorpus(ctx, data)
	m.data = append(m.data, data...)
}

func (m *Markov) buildCorpus(ctx context.Context, data []string) {
	startWords := m.store.StartWords(ctx)
	endWords := m.store.EndWords(ctx)
	corpus := m.store.Corpus(ctx)

	lo.ForEach(data, func(item string, index int) {
		words := strings.Split(item, " ")
		stateSize := m.stateSize

		start := strings.Join(
			lo.Slice(words, 0, stateSize),
			" ",
		)
		oldStartObj, idx, found := lo.FindIndexOf(startWords, func(f *Fragment) bool {
			return f.Words == start
		})

		if found {
			if !lo.Contains(oldStartObj.Refs, item) {
				oldStartObj.Refs = append(oldStartObj.Refs, item)
				startWords[idx] = oldStartObj
			}
		} else {
			startWords = append(startWords, &Fragment{
				Words: start,
				Refs:  []string{item},
			})
		}

		end := strings.Join(
			lo.Slice(words, len(words)-stateSize, len(words)),
			" ",
		)
		oldEndObj, idx, found := lo.FindIndexOf(endWords, func(f *Fragment) bool {
			return f.Words == end
		})

		if found {
			if !lo.Contains(oldEndObj.Refs, item) {
				oldEndObj.Refs = append(oldEndObj.Refs, item)
				endWords[idx] = oldEndObj
			}
		} else {
			endWords = append(endWords, &Fragment{
				Words: end,
				Refs:  []string{item},
			})
		}

		for i := 0; i < len(words)-1; i++ {
			curr := strings.Join(
				lo.Slice(words, i, i+stateSize),
				" ",
			)
			next := strings.Join(
				lo.Slice(words, i+stateSize, i+stateSize*2),
				" ",
			)

			if next == "" || len(strings.Split(next, " ")) != stateSize {
				continue
			}

			if block, ok := corpus[curr]; ok {
				oldObj, idx, found := lo.FindIndexOf(block, func(f *Fragment) bool {
					return f.Words == next
				})

				if found {
					oldObj.Refs = append(oldObj.Refs, item)
					corpus[curr][idx] = oldObj
				} else {
					corpus[curr] = append(corpus[curr], &Fragment{
						Words: next,
						Refs:  []string{item},
					})
				}
			} else {
				corpus[curr] = []*Fragment{
					{Words: next, Refs: []string{item}},
				}
			}
		}
	})

	m.store.SetStartWords(ctx, startWords)
	m.store.SetEndWords(ctx, endWords)
	m.store.SetCorpus(ctx, corpus)
}

// Generate generates a sentence using the data it's been trained on.
func (m *Markov) Generate(ctx context.Context, opts ...GenerateOption) (*Result, error) {
	startWords := m.store.StartWords(ctx)
	endWords := m.store.EndWords(ctx)
	corpus := m.store.Corpus(ctx)

	if len(corpus) == 0 {
		return nil, fmt.Errorf("markov: corpus is empty")
	}

	g := &generateOptions{}

	for _, o := range opts {
		o(g)
	}

	maxTries := lo.Ternary(g.MaxTries > 0, g.MaxTries, defaultMaxTries)

	var tries int
	for tries = 1; tries <= maxTries; tries++ {
		ended := false

		arr := []*Fragment{lo.Sample(startWords)}
		score := 0

		for innerTries := 0; innerTries < maxTries; innerTries++ {
			block := lo.LastOrEmpty(arr)
			state := lo.Sample(corpus[block.Words])

			if state == nil {
				break
			}

			arr = append(arr, state)
			score += len(corpus[block.Words]) - 1

			if _, found := lo.Find(endWords, func(f *Fragment) bool {
				return f.Words == state.Words
			}); found {
				ended = true
				break
			}
		}

		sentence := strings.TrimSpace(
			strings.Join(
				lo.Map(arr, func(f *Fragment, index int) string {
					return f.Words
				}),
				" ",
			),
		)

		uniqRefs := lo.Uniq(
			lo.Flatten(
				lo.Map(arr, func(f *Fragment, index int) []string {
					return f.Refs
				}),
			),
		)

		result := &Result{
			String: sentence,
			Score:  score,
			Refs:   uniqRefs,
			Tries:  tries,
		}

		if !ended || (g.Filter != nil && g.Filter(result)) {
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("markov: failed to build sentence after %d tries", tries-1)
}
