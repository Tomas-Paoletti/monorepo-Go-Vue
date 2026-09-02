// Comando migrate: crea/actualiza el esquema de la base con GORM AutoMigrate.
// Es el equivalente de `php artisan migrate`. Idempotente: se puede correr
// varias veces sin romper nada.
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Tomas-Paoletti/roulette-movies/internal/storage/gormdb"
)

func main() {
	dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "postgres connection string")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("missing DATABASE_URL or -db")
	}

	start := time.Now()

	db, err := gormdb.Open(*dbURL)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("running migrations...")
	if err := gormdb.Migrate(db); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("done: schema up to date in %s", time.Since(start).Round(time.Millisecond))
}