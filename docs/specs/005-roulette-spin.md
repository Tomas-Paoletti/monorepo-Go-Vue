# 005 — Ruleta: girar sobre una lista

## Objetivo

Elegir una película al azar entre las de una lista, con filtros opcionales. Es el corazón
de la app: junta las tres capas (storage trae candidatas, el dominio `roulette` decide,
HTTP expone el resultado).

## Contrato (HTTP)

### `POST /lists/:id/spin`
Body opcional (todos los campos opcionales):
```json
{ "genre": "Comedy", "yearFrom": 1990, "yearTo": 2010,
  "minRating": 3.5, "excludeIds": [1, 2] }
```
Respuestas:
- `200` → `{ "winner": movieResponse }`.
- `400` si el id no es UUID o el body no es JSON válido.
- `404` si la lista no existe.
- `422` si ninguna película de la lista pasa el filtro (`roulette.ErrNoCandidates`).

## Modelo de datos

Lee `list_movies` (ids) + `movies` (datos). No escribe (por ahora no se persiste el historial
de `spins`; queda como iteración futura documentada acá).

## Criterios de aceptación

- [ ] Lista con ≥1 película sin filtros → `200` con una `winner` de la lista.
- [ ] Filtro que no matchea nada → `422`.
- [ ] `excludeIds` con la única candidata → `422`.
- [ ] Lista inexistente → `404`.

## Plan de test

- El dominio `roulette` ya tiene tests unitarios (determinismo con `SpinWith` + `Randomizer` fake).
- El handler se prueba manual con `curl` sobre una lista con películas.

## Notas de aprendizaje

- **Separación dominio/infra:** el handler NO decide la película. Trae candidatas de storage,
  arma un `roulette.Filter` y delega en `roulette.Spin` (dominio puro, testeable sin DB ni HTTP).
  La aleatoriedad vive detrás de la interfaz `Randomizer` para poder testear determinista.
- **Códigos HTTP con sentido:** "no hay candidatas" no es un `500` (no fue un error del server)
  ni un `404` (la lista existe); es `422 Unprocessable Entity`: la petición es válida pero no se
  puede satisfacer con esos filtros.
- **Futuro — persistir `spins`:** guardar cada giro (lista, película, timestamp, si se aceptó)
  permitiría "no repetir" entre sesiones y estadísticas. La tabla ya existe (spec 001).
