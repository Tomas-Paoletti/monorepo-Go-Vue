package gormdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tomas-Paoletti/roulette-movies/internal/list"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ListRepo maneja las listas y su tabla pivote list_movies.
type ListRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) *ListRepo {
	return &ListRepo{db: db}
}

// toDomain mapea la fila de persistencia al dominio list.List.
func listToDomain(row List) list.List {
	l := list.List{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
	}
	if row.Description != nil {
		l.Description = *row.Description
	}
	return l
}

func (r *ListRepo) Create(ctx context.Context, name, description string) (list.List, error) {
	row := List{Name: name}
	if description != "" {
		row.Description = &description
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return list.List{}, fmt.Errorf("creating list: %w", err)
	}
	return listToDomain(row), nil
}

func (r *ListRepo) All(ctx context.Context) ([]list.List, error) {
	var rows []List
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("listing lists: %w", err)
	}
	out := make([]list.List, len(rows))
	for i, row := range rows {
		out[i] = listToDomain(row)
	}
	return out, nil
}

func (r *ListRepo) Get(ctx context.Context, id uuid.UUID) (list.List, error) {
	var row List
	err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return list.List{}, list.ErrNotFound
	}
	if err != nil {
		return list.List{}, fmt.Errorf("getting list: %w", err)
	}
	return listToDomain(row), nil
}

func (r *ListRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&List{}, "id = ?", id)
	if tx.Error != nil {
		return fmt.Errorf("deleting list: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return list.ErrNotFound
	}
	return nil
}

// AddMovie agrega una película a la lista. Idempotente: si ya está, no falla.
func (r *ListRepo) AddMovie(ctx context.Context, listID uuid.UUID, movieID int) error {
	row := ListMovie{ListID: listID, MovieID: movieID}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&row).Error
	if err != nil {
		return fmt.Errorf("adding movie to list: %w", err)
	}
	return nil
}

func (r *ListRepo) RemoveMovie(ctx context.Context, listID uuid.UUID, movieID int) error {
	err := r.db.WithContext(ctx).
		Delete(&ListMovie{}, "list_id = ? AND movie_id = ?", listID, movieID).Error
	if err != nil {
		return fmt.Errorf("removing movie from list: %w", err)
	}
	return nil
}

// MovieIDs devuelve los ids de películas de una lista (para armar candidatas).
func (r *ListRepo) MovieIDs(ctx context.Context, listID uuid.UUID) ([]int, error) {
	var ids []int
	err := r.db.WithContext(ctx).
		Model(&ListMovie{}).
		Where("list_id = ?", listID).
		Pluck("movie_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("listing movie ids: %w", err)
	}
	return ids, nil
}
