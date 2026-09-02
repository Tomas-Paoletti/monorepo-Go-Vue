// Comando import: carga el CSV de MovieLens a la tabla `movies`.
// Es el equivalente de un seeder de Laravel, pero como binario CLI.
// Reusa importer.ReadMovies (dominio) y mapea a gormdb.Movie (persistencia).
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Tomas-Paoletti/roulette-movies/internal/importer"
	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
	"github.com/Tomas-Paoletti/roulette-movies/internal/storage/gormdb"
	"github.com/lib/pq"
	"gorm.io/gorm/clause"
)

func main() {
	dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "postgres connection string")
	dataPath := flag.String("data", "data/movies.csv", "ruta al CSV de MovieLens")
	batchSize := flag.Int("batch", 500, "tamaño de lote para CreateInBatches")
	truncate := flag.Bool("truncate", false, "vaciar movies antes de importar")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("missing DATABASE_URL or -db")
	}

	// 1. Leer el CSV a dominio.
	movies, err := importer.ReadMovies(*dataPath)
	if err != nil {
		log.Fatalf("reading csv: %v", err)
	}
	log.Printf("read %d movies from %s", len(movies), *dataPath)

	// 2. Mapear dominio -> persistencia.
	rows := make([]gormdb.Movie, 0, len(movies))
	for _, m := range movies {
		rows = append(rows, toRow(m))
	}

	// 3. Conectar.
	db, err := gormdb.Open(*dbURL)
	if err != nil {
		log.Fatalf("connecting: %v", err)
	}

	start := time.Now()

	if *truncate {
		if err := db.Exec(`TRUNCATE movies RESTART IDENTITY CASCADE`).Error; err != nil {
			log.Fatalf("truncating: %v", err)
		}
		log.Println("truncated movies")
	}

	// 4. Insertar en lotes, idempotente (ON CONFLICT (id) DO NOTHING).
	tx := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).CreateInBatches(rows, *batchSize)
	if tx.Error != nil {
		log.Fatalf("inserting: %v", tx.Error)
	}

	var total int64
	db.Model(&gormdb.Movie{}).Count(&total)
	log.Printf("done: inserted %d rows (table now has %d) in %s",
		tx.RowsAffected, total, time.Since(start).Round(time.Millisecond))
}

// toRow mapea el dominio a la fila de persistencia: los "sin dato" del dominio
// (0 / "") se vuelven NULL (punteros nil) en la DB.
func toRow(m movie.Movie) gormdb.Movie {
	r := gormdb.Movie{
		ID:     m.ID,
		Title:  m.Title,
		Genres: pq.StringArray(m.Genres),
	}
	if m.Year != 0 {
		r.Year = &m.Year
	}
	if m.AltTitle != "" {
		r.AltTitle = &m.AltTitle
	}
	if m.TMDBID != 0 {
		r.TMDBID = &m.TMDBID
	}
	if m.IMDBID != "" {
		r.IMDBID = &m.IMDBID
	}
	return r
}
