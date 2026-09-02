package roulette

import "github.com/Tomas-Paoletti/roulette-movies/internal/movie"

type Filter struct {
	Genre       string
	YearFrom    int
	YearTo      int
	MinRating   float64
	ExcludeIDs  []int
}

func (f Filter) Matches(m movie.Movie) bool {
	if f.Genre != "" && !m.HasGenre(f.Genre) {
		return false
	}
	if f.YearFrom > 0 && (m.Year == 0 || m.Year < f.YearFrom) {
		return false
	}
	if f.YearTo > 0 && (m.Year == 0 || m.Year > f.YearTo) {
		return false
	}
	if f.MinRating > 0 && m.Rating < f.MinRating {
		return false
	}
	for _, id := range f.ExcludeIDs {
		if m.ID == id {
			return false
		}
	}
	return true
}