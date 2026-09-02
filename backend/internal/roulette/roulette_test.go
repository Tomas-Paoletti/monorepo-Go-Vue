package roulette

import (
	"testing"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
)

type fakeRand struct{ value int }

func (r fakeRand) IntN(n int) int { return r.value }

func TestSpinWith(t *testing.T) {
	movies := []movie.Movie{
		{ID: 1, Title: "A", Year: 1990, Genres: []string{"Horror"}},
		{ID: 2, Title: "B", Year: 2020, Genres: []string{"Comedy"}},
		{ID: 3, Title: "C", Year: 2021, Genres: []string{"Horror"}},
	}

	got, err := SpinWith(movies, Filter{Genre: "Horror"}, fakeRand{value: 1})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 3 {   // índice 1 de las filtradas [A, C]
		t.Errorf("got %d want 3", got.ID)
	}
}