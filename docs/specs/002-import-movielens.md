# 002 — Importador del dataset MovieLens (CSV → Postgres)

## Objetivo

Cargar el catálogo de películas del CSV de MovieLens (`data/movies.csv`) a la tabla `movies`
usando el ORM. Es el equivalente de un seeder de Laravel, pero como binario CLI: lee el CSV,
mapea el dominio (`movie.Movie`) al modelo de persistencia (`gormdb.Movie`) e inserta en lotes.

## Contrato (interfaz de código)

- Binario `cmd/import`:
  - `-db`   → DSN de Postgres (default `DATABASE_URL`).
  - `-data` → ruta al CSV (default `data/movies.csv`).
  - `-batch` → tamaño de lote para `CreateInBatches` (default 500).
  - `-truncate` → si está, vacía `movies` antes de importar (import limpio).
  - Sin `-truncate`, la inserción es **idempotente**: `ON CONFLICT (id) DO NOTHING`.
- Reusa `importer.ReadMovies(path) ([]movie.Movie, error)` (ya existente).

## Modelo de datos

Escribe solo en **movies**. Mapeo dominio → persistencia:

| `movie.Movie` (dominio) | `gormdb.Movie` (persistencia) | Nota |
|---|---|---|
| `ID int`            | `ID int`                | movieId de MovieLens (PK, no autoincremental) |
| `Title string`      | `Title string`          | |
| `AltTitle string`   | `AltTitle *string`      | `nil` si vacío |
| `Year int` (0=desc) | `Year *int`             | `nil` si 0 |
| `Genres []string`   | `Genres pq.StringArray` | `text[]` en Postgres |
| `TMDBID int` (0)    | `TMDBID *int`           | `nil` si 0 (el CSV base no lo trae) |
| `IMDBID string`     | `IMDBID *string`        | `nil` si vacío |

Rating/RatingCount quedan en su default (0): el CSV base no los trae; se completarían con
`ratings.csv` en una iteración futura.

## Criterios de aceptación

- [ ] `go run ./cmd/import` contra la DB migrada inserta ~9.742 filas sin error.
- [ ] `SELECT count(*) FROM movies;` ≈ 9.742.
- [ ] Correrlo dos veces (sin `-truncate`) no duplica ni falla (idempotente por `ON CONFLICT`).
- [ ] `-truncate` deja la tabla con exactamente las filas del CSV.
- [ ] Una fila conocida (`id=1`) queda con `title='Toy Story'`, `year=1995`, `genres` con 5 géneros.
- [ ] `go build ./...` y `go vet ./...` en verde.

## Plan de test

- Verificación manual contra Postgres local: correr, contar filas, `SELECT * FROM movies WHERE id=1`.
- El parsing de CSV ya está cubierto por los tests de `internal/movie` (ParseTitle) y el importer.

## Notas de aprendizaje

- **`CreateInBatches` (ORM) vs `COPY`/`pgx.CopyFrom` (lo que hacía el tutorial):** el tutorial
  usaba `COPY`, la vía más rápida de Postgres para carga masiva. `CreateInBatches` arma
  `INSERT ... VALUES (...), (...)` por lote: más lento y con más overhead, pero sin escribir SQL
  ni mapear columnas a mano. Con ~9.7k filas la diferencia es imperceptible; a millones sí importa.
  Es el mismo trade-off que la spec 001: el ORM simplifica lo común y encarece lo especializado.
- **`clause.OnConflict{DoNothing: true}`:** hace el import repetible sin borrar. GORM lo traduce
  a `ON CONFLICT (id) DO NOTHING`. Es la forma idiomática de un "insert if not exists" en batch.
- **Dominio vs persistencia:** el dominio usa ceros/strings vacíos como "sin dato"; la DB prefiere
  `NULL`. El mapeo (a punteros) vive en el borde (el importador), no en el dominio ni en el modelo.
