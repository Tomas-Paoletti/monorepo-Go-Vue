package roulette

import (
	"errors"
	"math/rand/v2"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
)

var ErrNoCandidates = errors.New("no movie matches the filters")

// Randomizer existe SOLO para poder testear de forma determinista.
type Randomizer interface {
	IntN(n int) int
}

type realRand struct{}

func (realRand) IntN(n int) int { return rand.IntN(n) }

func Spin(candidates []movie.Movie, f Filter) (movie.Movie, error) {
	return SpinWith(candidates, f, realRand{})
}

func SpinWith(candidates []movie.Movie, f Filter, r Randomizer) (movie.Movie, error) {
	var valid []movie.Movie
	for _, m := range candidates {
		if f.Matches(m) {
			valid = append(valid, m)
		}
	}
	if len(valid) == 0 {
		return movie.Movie{}, ErrNoCandidates
	}
	return valid[r.IntN(len(valid))], nil
}