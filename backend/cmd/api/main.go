// Comando api: levanta el servidor HTTP.
// Acá vive el "wiring": config -> conexión -> repos -> server -> rutas.
// Es el único lugar que conoce las implementaciones concretas (gormdb);
// el paquete api solo conoce interfaces.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Tomas-Paoletti/roulette-movies/internal/api"
	"github.com/Tomas-Paoletti/roulette-movies/internal/config"
	"github.com/Tomas-Paoletti/roulette-movies/internal/storage/gormdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := gormdb.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// Repos concretos -> satisfacen las interfaces api.Catalog / api.Lists.
	catalog := gormdb.NewCatalogRepo(db)
	lists := gormdb.NewListRepo(db)

	server := api.NewServer(catalog, lists)
	handler := server.Router(corsOrigins())

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Arranca el server en una goroutine para poder escuchar señales.
	go func() {
		log.Printf("api listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	// Graceful shutdown: espera Ctrl+C / SIGTERM y cierra con timeout.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("bye")
}

// corsOrigins lee CORS_ORIGINS (coma-separado) o usa el dev server de Vite.
func corsOrigins() []string {
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		return strings.Split(v, ",")
	}
	return []string{"http://localhost:5173"}
}
