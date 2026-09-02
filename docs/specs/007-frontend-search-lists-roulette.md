# 007 — Features del frontend: búsqueda, listas y ruleta

## Objetivo

Las tres vistas que consumen la API y hacen la app usable: buscar películas, armar listas y
girar la ruleta sobre una lista.

## Vistas y flujo

### `features/movies/` — buscar
- `MoviesView.vue`: input de texto → `GET /movies?q=`. Muestra `MovieCard` en grilla.
- Botón "+ Lista" en cada card → agrega a la lista default (`ensureDefaultList`).
- Store `movies`: `results`, `loading`, `error`, `search(params)`.

### `features/lists/` — listas
- `ListsView.vue`: crear lista (`POST /lists`), listar (`GET /lists`), borrar (`DELETE`).
- `ListDetailView.vue`: `GET /lists/:id` (lista + películas). Quitar película (`DELETE .../movies/:id`).
- Store `lists`: `lists`, `fetchAll`, `create`, `remove`, `addMovie`, `ensureDefaultList`.

### Ruleta (dentro del detalle)
- Botón "🎲 Girar" → `POST /lists/:id/spin`. Muestra la ganadora.
- `422` (lista vacía / sin match) → mensaje "no hay película para elegir".

## Contrato (tipos TS ↔ JSON del backend)

- `Movie` ↔ `movieResponse` (camelCase: `altTitle`, `ratingCount`).
- `List` ↔ `listResponse` (`id` string uuid, `createdAt` ISO).
- `SpinFilters` ↔ `spinRequest` (`genre`, `yearFrom`, `yearTo`, `minRating`, `excludeIds`).

## Criterios de aceptación

- [ ] Buscar "toy" muestra Toy Story (con la DB importada).
- [ ] "+ Lista" desde el buscador crea/usa una lista y agrega la película.
- [ ] La vista de la lista muestra las películas agregadas.
- [ ] "Girar" sobre una lista con películas devuelve y muestra una ganadora.
- [ ] Girar sobre una lista vacía muestra el mensaje de error (no rompe).
- [ ] `npm run build` sigue pasando el type-check.

## Notas de aprendizaje

- **El componente no sabe de HTTP:** `MovieCard` solo recibe un `movie` y emite `action`. Quién
  llama a la API es la vista/estore. Componentes "tontos" = reusables (la misma card sirve para
  "agregar" en el buscador y "quitar" en el detalle, cambiando el label y el handler).
- **Estado derivado vs servidor:** tras `addMovie`/`removeMovie` se recarga desde el backend
  (`load()`) en vez de mutar a mano el array local. Más simple y siempre consistente con la DB;
  el trade-off (una request extra) no importa a esta escala.
- **Errores por código HTTP:** el `422` del spin se traduce a un mensaje amable en `catch`. El
  front no reimplementa la regla ("no hay candidatas"); confía en el status del backend.
