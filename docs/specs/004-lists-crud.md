# 004 — CRUD de listas y sus películas

## Objetivo

Permitir armar listas de películas: crear/listar/ver/borrar listas y agregar/quitar
películas de una lista. Es la base sobre la que después gira la ruleta (spec 005).

## Contrato (HTTP)

Base: `/api/v1`.

### `POST /lists`  — crear
Body: `{ "name": "Para el finde", "description": "opcional" }` (`name` requerido).
- `201` → `listResponse`. `400` si falta `name`.

### `GET /lists`  — listar
- `200` → `{ "lists": [ listResponse... ] }` (orden: más nuevas primero).

### `GET /lists/:id`  — ver una lista con sus películas
- `200` → `{ "list": listResponse, "movies": [ movieResponse... ] }`.
- `400` si el id no es un UUID válido. `404` si no existe.

### `DELETE /lists/:id`
- `204` sin body. `404` si no existe. (Borra en cascada `list_movies`.)

### `POST /lists/:id/movies`  — agregar película
Body: `{ "movieId": 123 }`.
- `204`. Idempotente: agregar dos veces no falla (`ON CONFLICT DO NOTHING`).

### `DELETE /lists/:id/movies/:movieID`
- `204`.

`listResponse`:
```json
{ "id": "uuid", "name": "...", "description": "...", "createdAt": "2026-08-22T..." }
```

## Modelo de datos

- `lists` — se inserta/borra la fila.
- `list_movies` — pivote (list_id, movie_id). FK a `lists` con `ON DELETE CASCADE`.
- Las películas de la lista se resuelven trayendo los `movie_id` y pidiéndolos al catálogo
  (`CatalogRepo.GetMany`).

## Criterios de aceptación

- [ ] Crear lista → `201` con `id` UUID y `createdAt`.
- [ ] Agregar movie 1 dos veces → sigue `204`, no duplica.
- [ ] `GET /lists/:id` devuelve la lista con la película agregada.
- [ ] Quitar la película → `GET` ya no la incluye.
- [ ] Borrar la lista → `204`; volver a `GET` → `404`.

## Plan de test

- Manual con `curl` sobre la DB local.
- (Futuro) test de `ListRepo` con testcontainers, afirmando el CASCADE.

## Notas de aprendizaje

- **UUID como PK de listas:** generado por Postgres (`gen_random_uuid()`, extensión pgcrypto).
  A diferencia de `movies.id` (que viene de MovieLens), acá el id lo pone la DB.
- **Idempotencia con `ON CONFLICT DO NOTHING`:** agregar una película ya presente no es un error;
  es la semántica natural de "que esté en la lista". Se traduce del `clause.OnConflict` de GORM.
- **La lista no guarda las películas, guarda referencias:** `list_movies` solo tiene ids; los
  datos de cada film salen del catálogo. Evita duplicar información y desincronizarla.
