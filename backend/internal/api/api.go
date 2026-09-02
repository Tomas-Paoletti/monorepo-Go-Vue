// Package api es la capa HTTP (Gin). Define las INTERFACES que necesita del
// lado del consumidor: no sabe que por debajo hay GORM/Postgres, solo pide
// "algo que sepa buscar películas" y "algo que sepa manejar listas". El wiring
// (qué implementación concreta se inyecta) vive en cmd/api.
package api

import (
	"context"

	"github.com/Tomas-Paoletti/roulette-movies/internal/list"
	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Catalog es lo que la API necesita para servir el catálogo.
type Catalog interface {
	Search(ctx context.Context, q movie.Query) ([]movie.Movie, error)
	Get(ctx context.Context, id int) (movie.Movie, error)
	GetMany(ctx context.Context, ids []int) ([]movie.Movie, error)
}

// Lists es lo que la API necesita para manejar listas.
type Lists interface {
	Create(ctx context.Context, name, description string) (list.List, error)
	All(ctx context.Context) ([]list.List, error)
	Get(ctx context.Context, id uuid.UUID) (list.List, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddMovie(ctx context.Context, listID uuid.UUID, movieID int) error
	RemoveMovie(ctx context.Context, listID uuid.UUID, movieID int) error
	MovieIDs(ctx context.Context, listID uuid.UUID) ([]int, error)
}

// Server agrupa las dependencias de los handlers.
type Server struct {
	catalog Catalog
	lists   Lists
}

func NewServer(catalog Catalog, lists Lists) *Server {
	return &Server{catalog: catalog, lists: lists}
}

// Router arma el motor Gin con las rutas y el middleware. allowedOrigins son
// los orígenes CORS permitidos (ej. el dev server de Vite en :5173).
func (s *Server) Router(allowedOrigins []string) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware(allowedOrigins))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/movies", s.searchMovies)
		v1.GET("/movies/:id", s.getMovie)

		v1.GET("/lists", s.getLists)
		v1.POST("/lists", s.createList)
		v1.GET("/lists/:id", s.getList)
		v1.DELETE("/lists/:id", s.deleteList)
		v1.POST("/lists/:id/movies", s.addListMovie)
		v1.DELETE("/lists/:id/movies/:movieID", s.removeListMovie)
		v1.POST("/lists/:id/spin", s.spinList)
	}

	return r
}
