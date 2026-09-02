# 001 — Esquema y migraciones con ORM (GORM)

## Objetivo

Crear el esquema de la base (tablas, índices, extensiones) desde Go usando el ORM **GORM**
(`AutoMigrate`), en vez de SQL crudo + golang-migrate como planteaba el tutorial. Es el
equivalente de `php artisan migrate`: declarás el modelo, el ORM aplica la diferencia.

## Contrato (interfaz de código)

- `gormdb.Open(dsn string) (*gorm.DB, error)` — abre la conexión a Postgres.
- `gormdb.Migrate(db *gorm.DB) error` — deja el esquema al día. Idempotente.
- Binario `cmd/migrate`: lee `DATABASE_URL` (o flag `-db`), abre, migra, loguea y sale.
- `config.Load() (config.Config, error)` — lee `DATABASE_URL` y `PORT` del entorno.

## Modelo de datos

Tablas (nombres en inglés):

- **movies** — `id` (PK, no autoincremental: es el movieId de MovieLens), `tmdb_id`, `imdb_id`,
  `title`, `alt_title`, `year`, `genres text[]`, `rating real`, `rating_count`.
- **lists** — `id uuid` (PK, `gen_random_uuid()`), `name`, `description`, `created_at`.
- **list_movies** — PK compuesta (`list_id`, `movie_id`), `added_at`, `watched`, `watched_at`.
  FK a `lists` (ON DELETE CASCADE) y a `movies`.
- **favorites** — `movie_id` (PK, FK a movies), `added_at`.
- **spins** — `id uuid` (PK), `list_id` (FK CASCADE), `movie_id` (FK), `spun_at`, `accepted`.

Índices:
- GIN sobre `movies.genres` (filtrar por género con `&&`).
- GIN trgm sobre `movies.title` (búsqueda por similitud → extensión `pg_trgm`).
- btree en `movies.year` y `movies.rating DESC`.
- compuesto `spins (list_id, spun_at DESC)`.

Extensiones necesarias **antes** de crear tablas/índices: `pgcrypto` (para `gen_random_uuid()`)
y `pg_trgm` (para el índice trgm).

## Criterios de aceptación

- [ ] `go run ./cmd/migrate` contra una DB vacía crea las 5 tablas sin error.
- [ ] Correrlo dos veces seguidas no falla (idempotente).
- [ ] Existen los índices GIN sobre `genres` y `title`, y las extensiones `pgcrypto`/`pg_trgm`.
- [ ] Las FK de `list_movies` y `spins` a `lists` tienen `ON DELETE CASCADE`.
- [ ] `go build ./...` y `go test ./...` siguen en verde.

## Plan de test

- Verificación manual contra Postgres local (docker): `\d movies`, `\di` y `\dx` en psql.
- (Opcional, más adelante) test de integración con `testcontainers-go` que levante Postgres,
  corra `Migrate` y afirme que las tablas existen.

## Notas de aprendizaje

- **ORM vs SQL crudo:** GORM infiere el DDL de los tags del struct. Cómodo, pero cosas como
  el índice GIN con `gin_trgm_ops` o el tipo `text[]` **no** son obvias: hay que declararlas con
  tags específicos (`type:gin,option:gin_trgm_ops`, `type:text[]`) y crear la extensión a mano
  con `db.Exec(CREATE EXTENSION ...)` antes del `AutoMigrate`. Ese es justo el trade-off que el
  tutorial advertía: el ORM te ahorra el SQL simple pero te esconde (y complica) el SQL potente.
- **`AutoMigrate` no borra ni renombra columnas.** Solo agrega lo que falta. Para cambios
  destructivos en un proyecto real harían falta migraciones versionadas; acá alcanza para aprender.
- **DSN por entorno:** en Go no hay `config/database.php`; se lee `DATABASE_URL` con `os.Getenv`.
