# Fix: los tests no corrían (`go test ./...` fallaba con "build failed")

## Síntoma

Al correr los tests:

```
go test ./...
```

Salía:

```
# github.com/Tomas-Paoletti/roulette-movies/cmd/api
cmd\api\main.go:5:2: "strings" imported and not used
cmd\api\main.go:10:23: undefined: movie.Greeting
cmd\api\main.go:11:12: undefined: New
cmd\api\main.go:12:14: undefined: errors
cmd\api\main.go:12:29: undefined: ErrInvalidYear
FAIL	github.com/Tomas-Paoletti/roulette-movies/cmd/api [build failed]
ok  	github.com/Tomas-Paoletti/roulette-movies/internal/movie	0.211s
FAIL
```

## Por qué pasaba

Los tests del paquete `internal/movie` **estaban bien y pasaban** (`ok ... internal/movie`).

El problema es que `go test ./...` **compila todos los paquetes del módulo**, incluido
`cmd/api`. Como `cmd/api/main.go` no compilaba, Go marcaba `[build failed]` y el
comando completo devolvía `FAIL`, aunque los tests reales estuvieran en verde.

Errores concretos en `cmd/api/main.go`:

1. `"strings"` se importaba pero no se usaba (en Go eso es error de compilación, no warning).
2. `movie.Greeting()` — esa función no existe en el paquete `movie`.
3. `New(...)` se llamaba sin el prefijo del paquete; debía ser `movie.New(...)`.
4. `errors` se usaba (`errors.Is`) pero nunca se importaba.
5. `ErrInvalidYear` se usaba sin el prefijo del paquete; debía ser `movie.ErrInvalidYear`.

## Cómo se solucionó

Se reescribió `cmd/api/main.go` para que compile:

- Se quitó el import `"strings"` que no se usaba.
- Se agregó el import `"errors"`.
- Se eliminó la llamada a `movie.Greeting()` (función inexistente).
- Se usaron los identificadores con su prefijo de paquete: `movie.New` y `movie.ErrInvalidYear`.

`main.go` final:

```go
package main

import (
	"errors"
	"fmt"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
)

func main() {
	_, err := movie.New(1, "Matrix", 1500)
	fmt.Println(errors.Is(err, movie.ErrInvalidYear)) // true
	fmt.Println(err)                                  // "invalid year: 1500"
}
```

## Verificación

```
go test ./...
```

```
?   	github.com/Tomas-Paoletti/roulette-movies/cmd/api	[no test files]
ok  	github.com/Tomas-Paoletti/roulette-movies/internal/movie
```

## Nota para la próxima vez

Si `go test ./...` dice `[build failed]` en un paquete, **no** es un test que falla:
es un paquete que no compila. Fijate en el nombre del paquete que aparece arriba del
`FAIL` (ahí era `cmd/api`) — el error está en el código de ese paquete, no en los tests.
Para correr solo los tests de un paquete sin depender del resto:
`go test ./internal/movie/`.
