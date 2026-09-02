// Package gormdb contiene los modelos de PERSISTENCIA (GORM) y la migración
// del esquema. Ojo: NO son el dominio. El dominio puro sigue siendo
// movie.Movie; estos structs solo describen cómo se guardan las filas.
package gormdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Movie mapea la tabla `movies`. El id es el movieId de MovieLens, así que
// NO es autoincremental: lo traemos nosotros.
type Movie struct {
	ID          int            `gorm:"primaryKey;autoIncrement:false"`
	TMDBID      *int           `gorm:"column:tmdb_id"`
	IMDBID      *string        `gorm:"column:imdb_id"`
	Title       string         `gorm:"not null;index:idx_movies_title_trgm,type:gin,option:gin_trgm_ops"`
	AltTitle    *string        `gorm:"column:alt_title"`
	Year        *int           `gorm:"index:idx_movies_year"`
	Genres      pq.StringArray `gorm:"type:text[];not null;default:'{}';index:idx_movies_genres,type:gin"`
	Rating      float64        `gorm:"type:real;not null;default:0;index:idx_movies_rating,sort:desc"`
	RatingCount int            `gorm:"column:rating_count;not null;default:0"`
}

func (Movie) TableName() string { return "movies" }

// List mapea la tabla `lists`.
type List struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Description *string
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:now()"`
}

func (List) TableName() string { return "lists" }

// ListMovie es la tabla pivote `list_movies` con clave primaria compuesta.
type ListMovie struct {
	ListID    uuid.UUID  `gorm:"column:list_id;type:uuid;primaryKey"`
	MovieID   int        `gorm:"column:movie_id;primaryKey"`
	AddedAt   time.Time  `gorm:"column:added_at;not null;default:now()"`
	Watched   bool       `gorm:"not null;default:false"`
	WatchedAt *time.Time `gorm:"column:watched_at"`

	List  List  `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE"`
	Movie Movie `gorm:"foreignKey:MovieID"`
}

func (ListMovie) TableName() string { return "list_movies" }

// Favorite mapea la tabla `favorites`.
type Favorite struct {
	MovieID int       `gorm:"column:movie_id;primaryKey"`
	AddedAt time.Time `gorm:"column:added_at;not null;default:now()"`

	Movie Movie `gorm:"foreignKey:MovieID"`
}

func (Favorite) TableName() string { return "favorites" }

// Spin mapea la tabla `spins`: el historial de la ruleta.
type Spin struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ListID   *uuid.UUID `gorm:"column:list_id;type:uuid;index:idx_spins_list,priority:1"`
	MovieID  *int       `gorm:"column:movie_id"`
	SpunAt   time.Time  `gorm:"column:spun_at;not null;default:now();index:idx_spins_list,priority:2,sort:desc"`
	Accepted *bool

	List  *List  `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE"`
	Movie *Movie `gorm:"foreignKey:MovieID"`
}

func (Spin) TableName() string { return "spins" }

// Models es el orden en que se migran (respeta las foreign keys: primero las
// tablas referenciadas, después las que referencian).
func Models() []any {
	return []any{
		&Movie{},
		&List{},
		&ListMovie{},
		&Favorite{},
		&Spin{},
	}
}
