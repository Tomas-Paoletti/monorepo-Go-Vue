// Package config centraliza la lectura de configuración desde el entorno.
// En Go no hay un archivo tipo config/database.php: se leen variables de
// entorno con os.Getenv. Load() falla temprano si falta algo obligatorio.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string // DSN de Postgres
	Port        string // puerto HTTP de la API
}

// Load lee la configuración del entorno. DATABASE_URL es obligatoria; PORT
// tiene default 8080.
func Load() (Config, error) {
	c := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        getEnv("PORT", "8080"),
	}
	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("missing required env var DATABASE_URL")
	}
	return c, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
