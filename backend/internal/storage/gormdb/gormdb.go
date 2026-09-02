package gormdb

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open abre la conexión a Postgres con GORM.
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}
	return db, nil
}

// Migrate deja el esquema al día. Es el equivalente de `php artisan migrate`,
// pero declarativo: GORM compara los modelos contra las tablas y aplica la
// diferencia. Idempotente: correrlo dos veces no rompe nada.
func Migrate(db *gorm.DB) error {
	// 1. Extensiones que necesitan las tablas ANTES de crearlas:
	//    - pgcrypto -> gen_random_uuid() para las PK de tipo uuid
	//    - pg_trgm  -> el índice GIN trgm sobre movies.title
	for _, ext := range []string{"pgcrypto", "pg_trgm"} {
		if err := db.Exec(fmt.Sprintf(`CREATE EXTENSION IF NOT EXISTS %s`, ext)).Error; err != nil {
			return fmt.Errorf("creating extension %s: %w", ext, err)
		}
	}

	// 2. AutoMigrate crea/actualiza tablas, columnas, FKs e índices declarados
	//    en los tags de los modelos.
	if err := db.AutoMigrate(Models()...); err != nil {
		return fmt.Errorf("auto-migrating: %w", err)
	}

	return nil
}
