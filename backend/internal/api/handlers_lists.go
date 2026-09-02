package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Tomas-Paoletti/roulette-movies/internal/list"
	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GET /api/v1/lists
func (s *Server) getLists(c *gin.Context) {
	lists, err := s.lists.All(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not list lists")
		return
	}
	out := make([]listResponse, len(lists))
	for i, l := range lists {
		out[i] = toListResponse(l)
	}
	c.JSON(http.StatusOK, gin.H{"lists": out})
}

// POST /api/v1/lists
func (s *Server) createList(c *gin.Context) {
	var req createListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "name is required")
		return
	}
	l, err := s.lists.Create(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not create list")
		return
	}
	c.JSON(http.StatusCreated, toListResponse(l))
}

// GET /api/v1/lists/:id  → la lista + sus películas.
func (s *Server) getList(c *gin.Context) {
	id, ok := parseListID(c)
	if !ok {
		return
	}
	ctx := c.Request.Context()

	l, err := s.lists.Get(ctx, id)
	if errors.Is(err, list.ErrNotFound) {
		notFound(c, "list not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get list")
		return
	}

	movies, err := s.listMovies(ctx, id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not get list movies")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":   toListResponse(l),
		"movies": toMovieResponses(movies),
	})
}

// DELETE /api/v1/lists/:id
func (s *Server) deleteList(c *gin.Context) {
	id, ok := parseListID(c)
	if !ok {
		return
	}
	err := s.lists.Delete(c.Request.Context(), id)
	if errors.Is(err, list.ErrNotFound) {
		notFound(c, "list not found")
		return
	}
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not delete list")
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/v1/lists/:id/movies  { "movieId": 123 }
func (s *Server) addListMovie(c *gin.Context) {
	id, ok := parseListID(c)
	if !ok {
		return
	}
	var req addMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "movieId is required")
		return
	}
	if err := s.lists.AddMovie(c.Request.Context(), id, req.MovieID); err != nil {
		respondError(c, http.StatusInternalServerError, "could not add movie")
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /api/v1/lists/:id/movies/:movieID
func (s *Server) removeListMovie(c *gin.Context) {
	id, ok := parseListID(c)
	if !ok {
		return
	}
	movieID, err := strconv.Atoi(c.Param("movieID"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid movie id")
		return
	}
	if err := s.lists.RemoveMovie(c.Request.Context(), id, movieID); err != nil {
		respondError(c, http.StatusInternalServerError, "could not remove movie")
		return
	}
	c.Status(http.StatusNoContent)
}

// listMovies trae las películas de una lista resolviendo sus ids contra el catálogo.
func (s *Server) listMovies(ctx context.Context, id uuid.UUID) ([]movie.Movie, error) {
	ids, err := s.lists.MovieIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.catalog.GetMany(ctx, ids)
}

// parseListID lee y valida el uuid de la URL. Si es inválido responde 400 y
// devuelve ok=false.
func parseListID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid list id")
		return uuid.Nil, false
	}
	return id, true
}
