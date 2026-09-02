package gormdb

import "github.com/Tomas-Paoletti/roulette-movies/internal/movie"

// toDomain convierte una fila de persistencia en el dominio puro. Los NULL de
// la DB (punteros nil) se vuelven los "sin dato" del dominio (0 / "").
func (r Movie) toDomain() movie.Movie {
	m := movie.Movie{
		ID:          r.ID,
		Title:       r.Title,
		Genres:      []string(r.Genres),
		Rating:      r.Rating,
		RatingCount: r.RatingCount,
	}
	if r.Year != nil {
		m.Year = *r.Year
	}
	if r.AltTitle != nil {
		m.AltTitle = *r.AltTitle
	}
	if r.TMDBID != nil {
		m.TMDBID = *r.TMDBID
	}
	if r.IMDBID != nil {
		m.IMDBID = *r.IMDBID
	}
	return m
}

// toDomainSlice mapea un lote de filas a dominio.
func toDomainSlice(rows []Movie) []movie.Movie {
	out := make([]movie.Movie, len(rows))
	for i, r := range rows {
		out[i] = r.toDomain()
	}
	return out
}
