package api

import (
	"time"

	"github.com/Tomas-Paoletti/roulette-movies/internal/list"
	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
)

// --- Responses (dominio -> JSON) -------------------------------------------
// Los DTO son la frontera: separan la forma interna (dominio) de la forma que
// viaja por HTTP. Así podemos cambiar una sin romper la otra.

type movieResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	AltTitle    string   `json:"altTitle,omitempty"`
	Year        int      `json:"year,omitempty"`
	Genres      []string `json:"genres"`
	Rating      float64  `json:"rating"`
	RatingCount int      `json:"ratingCount"`
}

func toMovieResponse(m movie.Movie) movieResponse {
	genres := m.Genres
	if genres == nil {
		genres = []string{} // JSON: [] en vez de null
	}
	return movieResponse{
		ID:          m.ID,
		Title:       m.Title,
		AltTitle:    m.AltTitle,
		Year:        m.Year,
		Genres:      genres,
		Rating:      m.Rating,
		RatingCount: m.RatingCount,
	}
}

func toMovieResponses(ms []movie.Movie) []movieResponse {
	out := make([]movieResponse, len(ms))
	for i, m := range ms {
		out[i] = toMovieResponse(m)
	}
	return out
}

type listResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func toListResponse(l list.List) listResponse {
	return listResponse{
		ID:          l.ID.String(),
		Name:        l.Name,
		Description: l.Description,
		CreatedAt:   l.CreatedAt,
	}
}

// --- Requests (JSON -> dominio) --------------------------------------------

type createListRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type addMovieRequest struct {
	MovieID int `json:"movieId" binding:"required"`
}

// spinRequest son los filtros opcionales de la ruleta.
type spinRequest struct {
	Genre     string  `json:"genre"`
	YearFrom  int     `json:"yearFrom"`
	YearTo    int     `json:"yearTo"`
	MinRating float64 `json:"minRating"`
	ExcludeIDs []int  `json:"excludeIds"`
}
