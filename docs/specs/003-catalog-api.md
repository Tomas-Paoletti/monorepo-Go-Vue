# 003 — API de catálogo (búsqueda de películas) + CORS

## Objetivo

Exponer el catálogo por HTTP: buscar películas con filtros y traer una por id.
Es el primer endpoint que consume el frontend. Incluye CORS para el dev server de Vite.

## Contrato (HTTP)

Base: `/api/v1`. Respuestas JSON.

### `GET /movies`
Query params (todos opcionales):
- `q` — substring del título (case-insensitive, `ILIKE`).
- `genre` — género exacto (overlap sobre `text[]`).
- `yearFrom`, `yearTo` — rango de año.
- `minRating` — rating mínimo.
- `limit` (default 20, tope 100), `offset` (paginado).

Respuesta `200`:
```json
{ "movies": [ { "id": 1, "title": "Toy Story", "year": 1995,
  "genres": ["Adventure","Animation"], "rating": 3.9, "ratingCount": 0 } ] }
```

### `GET /movies/:id`
- `200` → un `movieResponse`.
- `400` si el id no es numérico.
- `404` si no existe (`movie.ErrNotFound`).

### `GET /healthz`
- `200 {"status":"ok"}` — para probar que el server está vivo.

### CORS
- Origen permitido por default: `http://localhost:5173` (Vite). Configurable con
  `CORS_ORIGINS` (coma-separado). Métodos: GET, POST, DELETE, OPTIONS.

## Modelo de datos

Solo lectura sobre `movies`. La búsqueda por título aprovecha el índice GIN trgm;
el filtro por género usa el operador `&&` (overlap) sobre el índice GIN de `genres`.

## Criterios de aceptación

- [ ] `GET /healthz` → `200 {"status":"ok"}`.
- [ ] `GET /movies?q=toy` devuelve Toy Story.
- [ ] `GET /movies?genre=Comedy&limit=5` devuelve ≤5 películas, todas de Comedy.
- [ ] `GET /movies/1` → Toy Story; `GET /movies/999999` → `404`.
- [ ] Un fetch desde `localhost:5173` no es bloqueado por CORS.

## Plan de test

- Manual con `curl` (ver README) una vez migrada+importada la DB.
- (Futuro) test de `CatalogRepo` con testcontainers.

## Notas de aprendizaje

- **Interfaces del lado del consumidor:** `internal/api` define `Catalog` (lo que *necesita*),
  y `gormdb.CatalogRepo` la implementa. La API no importa `gormdb`: el wiring en `cmd/api`
  inyecta la implementación concreta. Así podríamos cambiar Postgres por otra cosa sin tocar
  los handlers. Es el patrón que el tutorial marca en el paso 12.
- **DTO ≠ dominio:** `movieResponse` (con tags `json:`) separa la forma HTTP de `movie.Movie`.
  Cambiar el JSON no obliga a tocar el dominio y viceversa.
- **CORS existe por el navegador:** en dev, front (`:5173`) y back (`:8080`) son orígenes
  distintos; sin CORS el browser bloquea la respuesta. En fase 2 (embed) desaparece.
