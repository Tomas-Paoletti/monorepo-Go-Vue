package api

import (
	"errors"
	"net/http"

	"github.com/Tomas-Paoletti/roulette-movies/internal/list"
	"github.com/Tomas-Paoletti/roulette-movies/internal/roulette"
	"github.com/gin-gonic/gin"
)

// POST /api/v1/lists/:id/spin
// Body opcional con filtros (spinRequest). Elige una película al azar entre
// las de la lista que pasen el filtro. Acá se juntan las tres capas:
// storage (traer candidatas) + dominio (roulette) + HTTP (esta función).
func (s *Server) spinList(c *gin.Context) {
	id, ok := parseListID(c)
	if !ok {
		return
	}
	ctx := c.Request.Context()

	// El filtro es opcional: si no mandan body, va vacío.
	var req spinRequest
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			respondError(c, http.StatusBadRequest, "invalid filter body")
			return
		}
	}

	// La lista tiene que existir.
	if _, err := s.lists.Get(ctx, id); err != nil {
		if errors.Is(err, list.ErrNotFound) {
			notFound(c, "list not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "could not get list")
		return
	}

	// Candidatas = películas de la lista.
	movies, err := s.listMovies(ctx, id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get list movies")
		return
	}

	filter := roulette.Filter{
		Genre:      req.Genre,
		YearFrom:   req.YearFrom,
		YearTo:     req.YearTo,
		MinRating:  req.MinRating,
		ExcludeIDs: req.ExcludeIDs,
	}

	winner, err := roulette.Spin(movies, filter)
	if errors.Is(err, roulette.ErrNoCandidates) {
		respondError(c, http.StatusUnprocessableEntity, "no movie matches the filters")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not spin")
		return
	}

	c.JSON(http.StatusOK, gin.H{"winner": toMovieResponse(winner)})
}
