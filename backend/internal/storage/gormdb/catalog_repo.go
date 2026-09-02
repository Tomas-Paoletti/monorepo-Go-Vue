package gormdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// CatalogRepo lee el catálogo de películas. Devuelve SIEMPRE el dominio puro
// (movie.Movie), no los structs de persistencia: el resto de la app no sabe
// que por debajo hay GORM.
type CatalogRepo struct {
	db *gorm.DB
}

func NewCatalogRepo(db *gorm.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}

// Search devuelve las películas que matchean, ordenadas por rating desc.
func (r *CatalogRepo) Search(ctx context.Context, q movie.Query) ([]movie.Movie, error) {
	tx := r.db.WithContext(ctx).Model(&Movie{})

	if q.Text != "" {
		tx = tx.Where("title ILIKE ?", "%"+q.Text+"%")
	}
	if q.Genre != "" {
		// overlap: la fila tiene el género si el array se cruza con {genre}.
		tx = tx.Where("genres && ?", pq.StringArray{q.Genre})
	}
	if q.YearFrom > 0 {
		tx = tx.Where("year >= ?", q.YearFrom)
	}
	if q.YearTo > 0 {
		tx = tx.Where("year <= ?", q.YearTo)
	}
	if q.MinRating > 0 {
		tx = tx.Where("rating >= ?", q.MinRating)
	}

	limit := q.Limit
	if limit <= 0 || limit > 100 {
		limit = 20 // default y tope, para no traer todo
	}

	var rows []Movie
	err := tx.Order("rating DESC, rating_count DESC").
		Limit(limit).
		Offset(q.Offset).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("searching movies: %w", err)
	}
	return toDomainSlice(rows), nil
}

// Get busca una película por id. Devuelve movie.ErrNotFound si no existe.
func (r *CatalogRepo) Get(ctx context.Context, id int) (movie.Movie, error) {
	var row Movie
	err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return movie.Movie{}, movie.ErrNotFound
	}
	if err != nil {
		return movie.Movie{}, fmt.Errorf("getting movie %d: %w", id, err)
	}
	return row.toDomain(), nil
}

// GetMany trae varias películas por id (para candidatas de una lista/ruleta).
func (r *CatalogRepo) GetMany(ctx context.Context, ids []int) ([]movie.Movie, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []Movie
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("getting movies: %w", err)
	}
	return toDomainSlice(rows), nil
}
