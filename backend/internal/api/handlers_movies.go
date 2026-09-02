package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
	"github.com/gin-gonic/gin"
)

// GET /api/v1/movies?q=&genre=&yearFrom=&yearTo=&minRating=&limit=&offset=
func (s *Server) searchMovies(c *gin.Context) {
	q := movie.Query{
		Text:      c.Query("q"),
		Genre:     c.Query("genre"),
		YearFrom:  atoiOr(c.Query("yearFrom"), 0),
		YearTo:    atoiOr(c.Query("yearTo"), 0),
		MinRating: atofOr(c.Query("minRating"), 0),
		Limit:     atoiOr(c.Query("limit"), 0),
		Offset:    atoiOr(c.Query("offset"), 0),
	}

	movies, err := s.catalog.Search(c.Request.Context(), q)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not search movies")
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": toMovieResponses(movies)})
}

// GET /api/v1/movies/:id
func (s *Server) getMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid movie id")
		return
	}

	m, err := s.catalog.Get(c.Request.Context(), id)
	if errors.Is(err, movie.ErrNotFound) {
		notFound(c, "movie not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get movie")
		return
	}
	c.JSON(http.StatusOK, toMovieResponse(m))
}

// atoiOr parsea un entero; si falla devuelve el default.
func atoiOr(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}

// atofOr parsea un float; si falla devuelve el default.
func atofOr(s string, def float64) float64 {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v
	}
	return def
}
