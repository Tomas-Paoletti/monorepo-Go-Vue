# Movie Roulette en Go — Tutorial paso a paso (para quien viene de Laravel)

Un solo documento: construís el proyecto de principio a fin, y en cada paso te digo
**qué estás haciendo, cómo se hace en Laravel, y por qué en Go es distinto**.

**Convención:** todo el código, nombres de variables, funciones, tipos, paquetes, columnas
de la base y campos JSON van en **inglés**. Los comentarios y esta guía, en español.

**Lo que vas a tener al final:** una app donde buscás películas, armás listas
("terror pendiente", "para ver con mi novia") y le das a una ruleta que elige una al azar
con filtros. Datos locales (MovieLens), API en Gin, frontend Vue, todo en un binario.

**Cómo leerlo:** cada paso tiene
- 🎯 **Qué hacés**
- 💻 **El código**
- 🐘 **En Laravel** — la equivalencia
- 🧠 **Por qué** — la razón de fondo

**Regla de oro:** pasos 1-10 a mano, sin agente. Después lo abrís.

---

# ETAPA 0 — Preparar el terreno

## Paso 0.1 — Instalar Go

**Windows 11:**
```powershell
winget install GoLang.Go
```
Cerrá y reabrí la terminal.

**Linux Mint / WSL:**
```bash
sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

```bash
go version
```

🐘 **En Laravel:** equivale a instalar PHP + Composer. La diferencia es que Go trae compilador,
formateador, gestor de dependencias, test runner y linter básico en un solo binario. No hay
`phpunit`, `php-cs-fixer`, `phpstan` por separado.

---

## Paso 0.2 — El proyecto

```bash
mkdir movie-roulette && cd movie-roulette
go mod init github.com/tomas/movie-roulette
```

Esto crea `go.mod`:
```
module github.com/tomas/movie-roulette

go 1.24
```

🐘 **En Laravel:** es `composer.json`. El `module` es como el `name` del paquete, pero acá
**es también la ruta de import**: `import "github.com/tomas/movie-roulette/internal/movie"`.

🧠 **Por qué:** Go no tiene registry central tipo Packagist. Las dependencias se resuelven por
URL del repo. Suena raro y funciona bien: publicar un paquete es hacer push con un tag.

---

## Paso 0.3 — El dataset

```bash
mkdir data && cd data
wget https://files.grouplens.org/datasets/movielens/ml-latest-small.zip
unzip ml-latest-small.zip
cd ..
```

**Windows:**
```powershell
mkdir data; cd data
Invoke-WebRequest "https://files.grouplens.org/datasets/movielens/ml-latest-small.zip" -OutFile ml.zip
Expand-Archive ml.zip -DestinationPath .
cd ..
```

Mirá qué hay:
```bash
head -3 data/ml-latest-small/movies.csv
# movieId,title,genres
# 1,Toy Story (1995),Adventure|Animation|Children|Comedy|Fantasy
```

9.742 películas, 1 MB. Sin API key, sin rate limit, sin red.

---

## Paso 0.4 — Postgres

```bash
docker start mi-postgres
docker exec -it mi-postgres psql -U postgres -c "CREATE DATABASE movie_roulette;"
export DATABASE_URL="postgres://postgres:TU_PASS@localhost:5432/movie_roulette?sslmode=disable"
```

🐘 **En Laravel:** el `.env` con `DB_CONNECTION`, `DB_HOST`, etc. En Go es una sola connection
string que leés con `os.Getenv`. No hay `config/database.php`.

---

# ETAPA 1 — El dominio (pasos 1-6, sin agente)

Acá construís el corazón de la app **sin base de datos, sin HTTP, sin nada**. Solo structs,
funciones y tests.

🧠 **Por qué así:** en Laravel arrancarías con `php artisan make:model Movie -mcr` y ya tenés
modelo, migración, controller y rutas. En Go se hace al revés: primero la lógica pura, después
las capas de afuera. Al principio se siente lento. La ventaja es que todo esto se testea en
microsegundos, sin levantar nada.

---

## Paso 1 — Tu primer paquete

🎯 Entender paquetes, visibilidad y ejecutables.

```bash
mkdir -p cmd/api internal/movie
```

`cmd/api/main.go`:
```go
package main

import (
	"fmt"

	"github.com/tomas/movie-roulette/internal/movie"
)

func main() {
	fmt.Println(movie.Greeting())
}
```

`internal/movie/movie.go`:
```go
package movie

func Greeting() string {
	return "movie package ready"
}
```

```bash
go run ./cmd/api
```

🐘 **En Laravel:** `cmd/api/main.go` es tu `public/index.php` — el entrypoint.
`internal/movie` es un namespace `App\Movie`.

🧠 **Por qué `cmd/` e `internal/`:**
- **`cmd/`** = ejecutables. Cada subcarpeta es un binario. Vas a tener tres: `api`, `import`, `seed`.
  En Laravel el entrypoint siempre es uno solo y los comandos son clases de Artisan; en Go cada
  comando es su propio programa.
- **`internal/`** = palabra reservada del compilador. Nadie fuera de tu módulo puede importar de
  ahí. Es `private` a nivel de proyecto. Laravel no tiene equivalente.

**🔬 Experimento (hacelo):** cambiá `Greeting` por `greeting` (minúscula). Compilá.

```
cannot refer to unexported name movie.greeting
```

🐘 **En Laravel** la visibilidad es por palabra clave (`public`/`private`). En Go es por
**mayúscula inicial**. Mayúscula = exportado, minúscula = privado del paquete. Sin excepciones,
sin keywords.

---

## Paso 2 — El struct Movie

🎯 Structs, métodos, y la trampa de valor vs puntero.

`internal/movie/movie.go`:
```go
package movie

import "strings"

type Movie struct {
	ID          int      // movieId de MovieLens
	TMDBID      int      // por si después querés posters
	IMDBID      string
	Title       string
	AltTitle    string   // título original en otro idioma
	Year        int      // 0 si no se conoce
	Genres      []string
	Rating      float64  // promedio de ratings
	RatingCount int
}

func (m Movie) HasGenre(g string) bool {
	for _, x := range m.Genres {
		if strings.EqualFold(x, g) {
			return true
		}
	}
	return false
}

func (m Movie) IsClassic() bool {
	return m.Year > 0 && m.Year < 1980
}

func (m *Movie) AddGenre(g string) {
	if !m.HasGenre(g) {
		m.Genres = append(m.Genres, g)
	}
}
```

🐘 **En Laravel** esto sería un modelo Eloquent:
```php
class Movie extends Model {
    protected $fillable = ['title', 'year', 'genres'];
    public function hasGenre(string $g): bool { ... }
}
```

🧠 **La diferencia grande:** el modelo de Eloquent es *dos cosas a la vez* — el struct de datos y
el acceso a la base. `Movie::find(1)` y `$movie->title` viven en la misma clase.

En Go se parte en dos: `movie.Movie` es **solo datos y reglas**. Guardar y traer es
responsabilidad de otro paquete (`storage`). Al principio parece más trabajo. La ganancia: podés
testear todas las reglas de negocio sin base de datos.

🧠 **Receiver por valor vs puntero:** `func (m Movie)` recibe una **copia**. `func (m *Movie)`
recibe la dirección y puede mutar. Si un método muta, va con puntero.

**🔬 Experimento (te va a pasar en producción si no lo hacés ahora):**
```go
movies := []Movie{{Title: "Matrix"}, {Title: "Alien"}}
for _, m := range movies {
	m.AddGenre("sci-fi")
}
fmt.Printf("%+v\n", movies)   // ← sin géneros
```

El `range` te da una **copia** en `m`. Modificás la copia. Arreglo:
```go
for i := range movies {
	movies[i].AddGenre("sci-fi")
}
```

🐘 **En PHP** los objetos se pasan por referencia y esto "simplemente funciona". Es el error #1 al
pasar a Go. Cometelo ahora, en un test, no en producción.

---

## Paso 3 — Errores

🎯 Entender que en Go el error es un valor de retorno, no una excepción.

`internal/movie/errors.go`:
```go
package movie

import "errors"

var (
	ErrEmptyTitle  = errors.New("empty title")
	ErrInvalidYear = errors.New("invalid year")
	ErrNotFound    = errors.New("movie not found")
)
```

`internal/movie/movie.go` (agregá):
```go
func New(id int, title string, year int) (Movie, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Movie{}, ErrEmptyTitle
	}
	if year != 0 && (year < 1888 || year > 2100) {
		return Movie{}, fmt.Errorf("%w: %d", ErrInvalidYear, year)
	}
	return Movie{ID: id, Title: title, Year: year}, nil
}
```

🐘 **En Laravel:**
```php
public static function create(string $title, int $year): self {
    if (trim($title) === '') {
        throw new EmptyTitleException();
    }
}

try {
    $m = Movie::create($t, $y);
} catch (EmptyTitleException $e) { }
```

**En Go:**
```go
m, err := movie.New(1, "", 1999)
if err != nil {
	if errors.Is(err, movie.ErrEmptyTitle) {
		// manejar
	}
	return err
}
```

🧠 **Por qué es mejor de lo que parece:** en PHP no sabés qué excepciones puede tirar una función
sin leer su código (o confiar en el docblock). En Go el error está en la **firma**: si devuelve
`(Movie, error)`, sabés que puede fallar. El compilador no te obliga a manejarlo, pero el linter
sí, y es imposible que se te escape una excepción que no viste venir.

🧠 **`%w` es la clave.** Envuelve el error manteniendo la cadena. Es como
`new Exception($msg, 0, $previous)` de PHP, pero se usa siempre, no solo a veces.

**🔬 Experimento:**
```go
_, err := New(1, "Matrix", 1500)
fmt.Println(errors.Is(err, ErrInvalidYear))  // true
fmt.Println(err)                             // "invalid year: 1500"
```
Cambiá `%w` por `%v` y volvé a correr: `errors.Is` da `false`. Ese es el punto de `%w`.

🧠 **Nota de diseño:** permito `year == 0`. MovieLens tiene películas sin año. Modelá el dato
ausente desde el principio en vez de descubrirlo cuando explota el importador.

---

## Paso 4 — Tests

🎯 Aprender el patrón table-driven, que es *el* estilo de Go.

`internal/movie/movie_test.go`:
```go
package movie

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name    string
		title   string
		year    int
		wantErr error
	}{
		{"valid", "Matrix", 1999, nil},
		{"no year allowed", "Babylon 5", 0, nil},
		{"empty title", "", 1999, ErrEmptyTitle},
		{"whitespace only", "   ", 1999, ErrEmptyTitle},
		{"year too old", "Something", 1500, ErrInvalidYear},
		{"year in future", "Something", 2200, ErrInvalidYear},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := New(1, c.title, c.year)
			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, c.wantErr) {
				t.Errorf("got %v, want %v", err, c.wantErr)
			}
		})
	}
}
```

```bash
go test -v ./internal/movie
```

🐘 **En Laravel** usarías Pest o PHPUnit con un `dataProvider`.

🧠 **Tres diferencias:**
1. **El test vive al lado del código**, no en `tests/`. `movie.go` y `movie_test.go` en la misma
   carpeta. Suena desordenado y termina siendo cómodo.
2. **Está en la stdlib.** No instalás nada. `go test` y listo.
3. **No hay asserts.** Comparás con `if` y llamás a `t.Errorf`. Al principio molesta; después
   notás que los mensajes de error son más claros porque los escribís vos. (Si extrañás los
   asserts, `testify` existe. Yo empezaría sin él.)

🧠 **`t.Errorf` vs `t.Fatalf`:** `Errorf` marca fallo y sigue. `Fatalf` corta el subtest. Usá
`Fatalf` cuando lo que sigue no tiene sentido sin eso.

**📝 Tarea:** escribí `TestHasGenre` cubriendo mayúsculas/minúsculas, slice vacío y género
inexistente.

---

## Paso 5 — Parsear los títulos

🎯 Primer contacto con datos reales sucios. `regexp`, `strings`, y por qué el orden importa.

MovieLens guarda los títulos así:
```
Toy Story (1995)
Matrix, The (1999)                                     ← artículo al final
Mermaid, The (Mei ren yu) (2016)                       ← + título alternativo
Seven (a.k.a. Se7en) (1995)
Babylon 5                                              ← sin año
```

`internal/movie/parse.go`:
```go
package movie

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	yearRe    = regexp.MustCompile(`^(.*?)\s*\((\d{4})\)\s*$`)
	altRe     = regexp.MustCompile(`^(.*?)\s*\(([^()]*)\)\s*$`)
	articleRe = regexp.MustCompile(`^(.*),\s+(The|A|An|Le|La|Les|Il|El|Los|Las|Der|Die|Das)$`)
)

// ParseTitle convierte "Matrix, The (1999)" en ("The Matrix", 1999, "").
func ParseTitle(raw string) (title string, year int, alt string) {
	rest := strings.TrimSpace(raw)

	if m := yearRe.FindStringSubmatch(rest); m != nil {      // 1. el año
		rest = m[1]
		year, _ = strconv.Atoi(m[2])
	}
	if m := altRe.FindStringSubmatch(rest); m != nil {       // 2. título alternativo
		rest = m[1]
		alt = m[2]
	}
	if m := articleRe.FindStringSubmatch(rest); m != nil {   // 3. el artículo
		rest = m[2] + " " + m[1]
	}

	return strings.TrimSpace(rest), year, alt
}
```

🧠 **El orden es la lección.** Si movés el artículo antes de sacar el paréntesis alternativo,
`"Mermaid, The (Mei ren yu) (2016)"` se rompe. Probalo al revés para verlo fallar.

🧠 **`regexp.MustCompile` a nivel de paquete:** se compila una sola vez al arrancar. Si lo ponés
dentro de la función, recompilás la regex en cada llamada — 9.742 veces en el importador. Es el
error de performance más común de Go.

🐘 **En PHP** `preg_match` cachea internamente, así que no te preocupabas por esto. En Go la
compilación es explícita.

Test (los casos están verificados contra el dataset real):
```go
func TestParseTitle(t *testing.T) {
	cases := []struct {
		raw, wantTitle, wantAlt string
		wantYear                int
	}{
		{"Toy Story (1995)", "Toy Story", "", 1995},
		{"Matrix, The (1999)", "The Matrix", "", 1999},
		{"Fantasia 2000 (1999)", "Fantasia 2000", "", 1999},
		{"Babylon 5", "Babylon 5", "", 0},
		{"Mermaid, The (Mei ren yu) (2016)", "The Mermaid", "Mei ren yu", 2016},
		{"Seven (a.k.a. Se7en) (1995)", "Seven", "a.k.a. Se7en", 1995},
		{"Amityville 3-D (1983)", "Amityville 3-D", "", 1983},
		{"City of Lost Children, The (Cité des enfants perdus, La) (1995)",
			"The City of Lost Children", "Cité des enfants perdus, La", 1995},
	}

	for _, c := range cases {
		t.Run(c.raw, func(t *testing.T) {
			title, year, alt := ParseTitle(c.raw)
			if title != c.wantTitle || year != c.wantYear || alt != c.wantAlt {
				t.Errorf("got (%q,%d,%q) want (%q,%d,%q)",
					title, year, alt, c.wantTitle, c.wantYear, c.wantAlt)
			}
		})
	}
}
```

**📝 Tarea:** `grep` en `movies.csv` buscando otros formatos raros. Agregá casos.

---

## Paso 6 — La ruleta

🎯 La lógica central de la app. Y tu primera interfaz, con un motivo real.

```bash
mkdir internal/roulette
```

`internal/roulette/filter.go`:
```go
package roulette

import "github.com/tomas/movie-roulette/internal/movie"

type Filter struct {
	Genre       string
	YearFrom    int
	YearTo      int
	MinRating   float64
	MinVotes    int
	ExcludeIDs  []int
}

func (f Filter) Matches(m movie.Movie) bool {
	if f.Genre != "" && !m.HasGenre(f.Genre) {
		return false
	}
	if f.YearFrom > 0 && (m.Year == 0 || m.Year < f.YearFrom) {
		return false
	}
	if f.YearTo > 0 && (m.Year == 0 || m.Year > f.YearTo) {
		return false
	}
	if f.MinRating > 0 && m.Rating < f.MinRating {
		return false
	}
	if f.MinVotes > 0 && m.RatingCount < f.MinVotes {
		return false
	}
	for _, id := range f.ExcludeIDs {
		if m.ID == id {
			return false
		}
	}
	return true
}
```

`internal/roulette/roulette.go`:
```go
package roulette

import (
	"errors"
	"math/rand/v2"

	"github.com/tomas/movie-roulette/internal/movie"
)

var ErrNoCandidates = errors.New("no movie matches the filters")

// Randomizer existe SOLO para poder testear de forma determinista.
type Randomizer interface {
	IntN(n int) int
}

type realRand struct{}

func (realRand) IntN(n int) int { return rand.IntN(n) }

func Spin(candidates []movie.Movie, f Filter) (movie.Movie, error) {
	return SpinWith(candidates, f, realRand{})
}

func SpinWith(candidates []movie.Movie, f Filter, r Randomizer) (movie.Movie, error) {
	var valid []movie.Movie
	for _, m := range candidates {
		if f.Matches(m) {
			valid = append(valid, m)
		}
	}
	if len(valid) == 0 {
		return movie.Movie{}, ErrNoCandidates
	}
	return valid[r.IntN(len(valid))], nil
}
```

Test determinista:
```go
type fakeRand struct{ value int }

func (r fakeRand) IntN(n int) int { return r.value }

func TestSpinWith(t *testing.T) {
	movies := []movie.Movie{
		{ID: 1, Title: "A", Year: 1990, Genres: []string{"Horror"}},
		{ID: 2, Title: "B", Year: 2020, Genres: []string{"Comedy"}},
		{ID: 3, Title: "C", Year: 2021, Genres: []string{"Horror"}},
	}

	got, err := SpinWith(movies, Filter{Genre: "Horror"}, fakeRand{value: 1})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 3 {   // índice 1 de las filtradas [A, C]
		t.Errorf("got %d want 3", got.ID)
	}
}
```

🐘 **En Laravel** para testear algo con `rand()` tendrías que mockear una fachada o inyectar un
servicio, con `Mockery::mock(RandomService::class)->shouldReceive('int')->andReturn(1)`.

🧠 **En Go:** definís una interfaz de **un método** y escribís un struct de 3 líneas. Sin librería
de mocks, sin configurar expectativas. Y `fakeRand` nunca declara que implementa `Randomizer` —
tiene el método `IntN(int) int`, con eso alcanza.

🧠 **Esto es lo más distinto de Go:** las interfaces son **implícitas**. No hay `implements`.
Cualquier tipo que tenga los métodos, cumple. Eso cambia dónde conviene definirlas — lo vemos en
el paso 12.

🧠 **Y fijate dónde está la interfaz:** en `roulette`, que es quien la **usa**. No en un paquete
`random/` junto a la implementación. Retené esto.

**📝 Tarea:** versión con pesos, donde las mejor puntuadas tienen más chance. Sumás pesos, sacás
random en ese rango, recorrés acumulando hasta pasarlo. Testeala con el fake.

**🎉 Punto de control:** tenés el corazón de la app funcionando y testeado. Sin base de datos, sin
HTTP, sin framework. `go test ./...` corre en milisegundos.

---

# ETAPA 2 — Los datos (pasos 7-10)

## Paso 7 — Leer CSV

🎯 `encoding/csv`, manejo de archivos, errores con contexto.

```bash
mkdir internal/importer
```

`internal/importer/csv.go`:
```go
package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tomas/movie-roulette/internal/movie"
)

func ReadMovies(path string) ([]movie.Movie, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 3

	if _, err := r.Read(); err != nil {          // descartar header
		return nil, fmt.Errorf("reading header: %w", err)
	}

	var out []movie.Movie
	line := 1
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		line++
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}

		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid movieId %q", line, rec[0])
		}

		title, year, alt := movie.ParseTitle(rec[1])

		var genres []string
		if rec[2] != "(no genres listed)" && rec[2] != "" {
			genres = strings.Split(rec[2], "|")
		}

		out = append(out, movie.Movie{
			ID: id, Title: title, AltTitle: alt, Year: year, Genres: genres,
		})
	}
	return out, nil
}
```

🐘 **En Laravel** harías `League\Csv` o `fgetcsv()` a mano, y el manejo de archivos lo cerraría el
garbage collector.

🧠 **`defer f.Close()`** es el `finally` de Go. Se ejecuta al salir de la función, pase lo que
pase. Ponelo **inmediatamente después** de abrir, no al final.

🧠 **El número de línea en el error** parece un detalle. Cuando falle en la fila 4.821 de 9.742 vas
a agradecerlo. Regla general: cada capa agrega contexto con `%w`.

🧠 **`"(no genres listed)"`** es un valor real del dataset. Los datos siempre tienen sorpresas que
no están en el README.

**📝 Tarea:** escribí `ReadLinks` (ojo: `tmdbId` puede venir vacío) y `ReadRatings` que devuelva
`map[int]Aggregate`.

---

## Paso 8 — Concurrencia (el paso importante)

🎯 Goroutines, channels, y el detector de race conditions.

🐘 **En Laravel** esto directamente no existe. Para procesar en paralelo usarías queues con
múltiples workers, que son **procesos** separados. Go tiene concurrencia **dentro** del proceso, y
es barata: una goroutine arranca con 2KB de stack. Podés tener cien mil.

`internal/importer/ratings.go`:
```go
package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type Aggregate struct {
	Sum   float64
	Count int
}

// Pipeline: 1 lector → N workers → 1 agregador.
func ReadRatingsConcurrent(path string, workers int) (map[int]Aggregate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 4
	r.ReuseRecord = true           // reusa el slice interno: menos allocations
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	type entry struct {
		movieID int
		rating  float64
	}

	lines := make(chan []string, 1000)
	entries := make(chan entry, 1000)

	// 1. lector
	var readErr error
	go func() {
		defer close(lines)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				return
			}
			if err != nil {
				readErr = err
				return
			}
			cp := make([]string, len(rec))   // ReuseRecord OBLIGA a copiar
			copy(cp, rec)
			lines <- cp
		}
	}()

	// 2. workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rec := range lines {
				id, err1 := strconv.Atoi(rec[1])
				rating, err2 := strconv.ParseFloat(rec[2], 64)
				if err1 != nil || err2 != nil {
					continue
				}
				entries <- entry{movieID: id, rating: rating}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(entries)
	}()

	// 3. agregador: UNA sola goroutine escribe el map → sin mutex
	out := make(map[int]Aggregate)
	for e := range entries {
		a := out[e.movieID]
		a.Sum += e.rating
		a.Count++
		out[e.movieID] = a
	}

	return out, readErr
}
```

**🔬 Los tres experimentos que valen más que cualquier tutorial:**

**1.** Sacá la copia (`lines <- rec` directo) y corré `go test -race ./...`. Vas a ver el data race
reportado. Ese es el bug sutil que te come una tarde.

**2.** Hacé que los workers escriban `out` directamente en vez de mandar al channel. `-race` lo
detecta al instante. **Los maps de Go no son thread-safe** — a diferencia de los arrays de PHP,
donde el tema ni existe porque no hay concurrencia real.

**3.** Medí con 1, 4 y 16 workers. Vas a descubrir que más no siempre es mejor: acá el cuello de
botella es el lector, no el parseo.

🧠 **La pregunta honesta:** ¿vale la pena la concurrencia acá? Con 100k líneas la versión
secuencial tarda milisegundos. **Medí las dos.** Aprender *cuándo no* usar goroutines es tan
importante como saber usarlas — y es lo que separa a alguien que sabe Go de alguien que lo usa.

---

## Paso 9 — El esquema y las migraciones

🎯 SQL a mano, índices que importan.

```bash
mkdir migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

`migrations/000001_init.up.sql`:
```sql
CREATE TABLE movies (
    id            INT PRIMARY KEY,
    tmdb_id       INT,
    imdb_id       TEXT,
    title         TEXT NOT NULL,
    alt_title     TEXT,
    year          INT,
    genres        TEXT[] NOT NULL DEFAULT '{}',
    rating        REAL   NOT NULL DEFAULT 0,
    rating_count  INT    NOT NULL DEFAULT 0
);

CREATE INDEX idx_movies_genres ON movies USING GIN (genres);
CREATE INDEX idx_movies_year   ON movies (year);
CREATE INDEX idx_movies_rating ON movies (rating DESC);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_movies_title_trgm ON movies USING GIN (title gin_trgm_ops);

CREATE TABLE lists (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE list_movies (
    list_id    UUID REFERENCES lists(id) ON DELETE CASCADE,
    movie_id   INT  REFERENCES movies(id),
    added_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    watched    BOOLEAN NOT NULL DEFAULT false,
    watched_at TIMESTAMPTZ,
    PRIMARY KEY (list_id, movie_id)
);

CREATE TABLE favorites (
    movie_id INT PRIMARY KEY REFERENCES movies(id),
    added_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE spins (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id    UUID REFERENCES lists(id) ON DELETE CASCADE,
    movie_id   INT  REFERENCES movies(id),
    spun_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted   BOOLEAN
);

CREATE INDEX idx_spins_list ON spins (list_id, spun_at DESC);
```

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

🐘 **En Laravel** tendrías `Schema::create('movies', function (Blueprint $t) { ... })` en PHP, con
`php artisan migrate`. Acá es SQL puro.

🧠 **Se siente un retroceso y no lo es.** El schema builder de Laravel te abstrae de Postgres, y
por eso nunca usaste `TEXT[]`, ni un índice GIN, ni `pg_trgm`. Con SQL directo usás la base al 100%.

- **`genres TEXT[]`** — array nativo. Filtrás con `genres && ARRAY['Horror']` y el índice GIN lo
  hace instantáneo. Sin tabla pivote.
- **`pg_trgm`** — búsqueda por similitud. "matrix" encuentra "The Matrix" sin full-text search.
- **`spins`** — el historial es lo que hace la app interesante: podés no repetir lo de los últimos 5.

🧠 **Sobre no tener ORM:** el equivalente de Eloquent en Go es GORM, y existe. Pero la comunidad
prefiere SQL explícito. Vamos a usar **sqlc**, que es distinto: escribís SQL, te genera structs y
funciones tipadas. Tenés el mapeo automático de un ORM con la transparencia del SQL. El SQL malo
falla en `sqlc generate`, no en runtime.

---

## Paso 10 — El importador

🎯 Tu primer binario útil. `flag`, `pgx`, carga masiva.

```bash
mkdir cmd/import
go get github.com/jackc/pgx/v5
```

`cmd/import/main.go`:
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tomas/movie-roulette/internal/importer"
	"github.com/tomas/movie-roulette/internal/movie"
)

func main() {
	dir := flag.String("dir", "./data/ml-latest-small", "folder with MovieLens CSVs")
	dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "postgres connection string")
	truncate := flag.Bool("truncate", false, "truncate movies before importing")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("missing DATABASE_URL or -db")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := run(ctx, *dir, *dbURL, *truncate); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(ctx context.Context, dir, dbURL string, truncate bool) error {
	start := time.Now()

	log.Println("reading movies.csv...")
	movies, err := importer.ReadMovies(filepath.Join(dir, "movies.csv"))
	if err != nil {
		return err
	}
	log.Printf("  %d movies", len(movies))

	log.Println("reading links.csv...")
	links, err := importer.ReadLinks(filepath.Join(dir, "links.csv"))
	if err != nil {
		return err
	}

	log.Println("reading ratings.csv...")
	ratings, err := importer.ReadRatings(filepath.Join(dir, "ratings.csv"))
	if err != nil {
		return err
	}
	log.Printf("  %d movies with ratings", len(ratings))

	// merge
	for i := range movies {
		if l, ok := links[movies[i].ID]; ok {
			movies[i].TMDBID = l.TMDBID
			movies[i].IMDBID = l.IMDBID
		}
		if a, ok := ratings[movies[i].ID]; ok && a.Count > 0 {
			movies[i].Rating = a.Sum / float64(a.Count)
			movies[i].RatingCount = a.Count
		}
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("connecting: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	if truncate {
		if _, err := pool.Exec(ctx, `TRUNCATE movies CASCADE`); err != nil {
			return fmt.Errorf("truncate: %w", err)
		}
	}

	n, err := insert(ctx, pool, movies)
	if err != nil {
		return err
	}

	log.Printf("done: %d rows in %s", n, time.Since(start).Round(time.Millisecond))
	return nil
}

func insert(ctx context.Context, pool *pgxpool.Pool, movies []movie.Movie) (int64, error) {
	rows := make([][]any, 0, len(movies))
	for _, m := range movies {
		var year, tmdb any
		if m.Year > 0 {
			year = m.Year
		}
		if m.TMDBID > 0 {
			tmdb = m.TMDBID
		}

		genres := m.Genres
		if genres == nil {
			genres = []string{}
		}

		rows = append(rows, []any{
			m.ID, tmdb, nullIfEmpty(m.IMDBID), m.Title,
			nullIfEmpty(m.AltTitle), year, genres, m.Rating, m.RatingCount,
		})
	}

	return pool.CopyFrom(ctx,
		pgx.Identifier{"movies"},
		[]string{"id", "tmdb_id", "imdb_id", "title", "alt_title",
			"year", "genres", "rating", "rating_count"},
		pgx.CopyFromRows(rows))
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
```

```bash
go run ./cmd/import -truncate
```

```
reading movies.csv...
  9742 movies
reading ratings.csv...
  9724 movies with ratings
done: 9742 rows in 380ms
```

🐘 **En Laravel** esto sería `php artisan make:command ImportMovies` con un `handle()`, y
probablemente `Movie::insert($chunk)` en chunks de 500.

🧠 **`CopyFrom` es el `COPY` de Postgres**, no un `INSERT` múltiple. 9.742 filas en 380ms. Con
`INSERT` fila por fila serían minutos. Eloquent no tiene equivalente directo.

🧠 **`any` para representar NULL:** en un `[]any`, `nil` se traduce a NULL. Si dejás el `int` en
cero, se guarda 0, no NULL. Por eso los chequeos. Go no tiene tipos nullables como `?int` de C#, y
esta es la forma idiomática de manejarlo en el borde con la DB.

🧠 **`[]string` → `TEXT[]`** directo, sin serializar. pgx lo maneja solo.

**Verificá:**
```sql
SELECT title, year, rating, rating_count, genres
FROM movies
WHERE genres && ARRAY['Horror'] AND rating_count > 50
ORDER BY rating DESC LIMIT 10;
```

**🎉 Punto de control:** tenés casi 10.000 películas en tu base, con puntajes reales.

---

# ETAPA 3 — La API (pasos 11-16)

**A partir de acá podés abrir el agente**, con el `CLAUDE.md` puesto y una spec por paso.

## Paso 11 — sqlc

🎯 SQL tipado. Tu reemplazo de Eloquent.

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
mkdir -p internal/storage/db
```

`sqlc.yaml`:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/query.sql"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "internal/storage/db"
        sql_package: "pgx/v5"
```

`db/query.sql`:
```sql
-- name: GetMovie :one
SELECT * FROM movies WHERE id = $1;

-- name: SearchMoviesByTitle :many
SELECT * FROM movies
WHERE similarity(title, $1) > 0.2
ORDER BY similarity(title, $1) DESC
LIMIT $2 OFFSET $3;

-- name: ListCandidates :many
SELECT m.* FROM movies m
JOIN list_movies lm ON lm.movie_id = m.id
WHERE lm.list_id = $1
  AND (sqlc.narg('genre')::text IS NULL OR m.genres && ARRAY[sqlc.narg('genre')::text])
  AND (sqlc.narg('year_from')::int IS NULL OR m.year >= sqlc.narg('year_from'))
  AND lm.watched = false;
```

```bash
sqlc generate
```

🐘 **La comparación con Eloquent:**
```php
Movie::where('id', $id)->first();     // Laravel
```
```go
q.GetMovie(ctx, id)                   // Go + sqlc
```

Casi igual de corto. La diferencia: en Laravel el SQL lo genera el ORM y no lo ves (hasta que
hacés `->toSql()` y descubrís el N+1). Acá el SQL lo escribiste vos, y sqlc te dio los tipos gratis.

🧠 **Si el SQL está mal, `sqlc generate` falla.** Un typo en un nombre de columna se detecta en tu
terminal, no en producción.

---

## Paso 12 — Servicios e interfaces

🎯 La pieza arquitectónica clave: **dónde vive la interfaz**.

```bash
mkdir internal/catalog
```

`internal/catalog/service.go`:
```go
package catalog

import (
	"context"
	"fmt"

	"github.com/tomas/movie-roulette/internal/movie"
)

// ⭐ La interfaz vive ACÁ, donde se CONSUME. No junto a la implementación.
type Repository interface {
	Get(ctx context.Context, id int) (movie.Movie, error)
	SearchByTitle(ctx context.Context, q string, limit, offset int) ([]movie.Movie, int, error)
}

type Service struct {
	repo Repository
}

func New(r Repository) *Service { return &Service{repo: r} }

func (s *Service) Get(ctx context.Context, id int) (movie.Movie, error) {
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		return movie.Movie{}, fmt.Errorf("getting movie %d: %w", id, err)
	}
	return m, nil
}
```

`internal/storage/postgres/catalog.go`:
```go
package postgres

type CatalogRepo struct{ q *db.Queries }

func NewCatalogRepo(q *db.Queries) *CatalogRepo { return &CatalogRepo{q: q} }

// No declara implementar nada. Tiene los métodos, alcanza.
func (r *CatalogRepo) Get(ctx context.Context, id int) (movie.Movie, error) {
	row, err := r.q.GetMovie(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return movie.Movie{}, movie.ErrNotFound
		}
		return movie.Movie{}, err
	}
	return toDomain(row), nil
}

// verificación en tiempo de compilación
var _ catalog.Repository = (*CatalogRepo)(nil)
```

🐘 **En Laravel** ponés la interfaz en `App\Repositories\MovieRepositoryInterface` y la
implementación en `App\Repositories\EloquentMovieRepository`, y las bindeás en el
`AppServiceProvider`. Las dos viven en la misma carpeta.

🧠 **En Go va al revés y hay una razón:** si la interfaz vive en `storage`, entonces `catalog`
tiene que importar `storage`. La dependencia apunta **hacia afuera** y rompiste la arquitectura.
Poniéndola en `catalog`:

- `catalog` no importa `storage`. Podés borrar todo `storage/` y `catalog` sigue compilando.
- La interfaz tiene exactamente los métodos que `catalog` necesita, ni uno más.
- Si `lists` necesita otros métodos, define **su propia** interfaz. Sí, aunque se repitan. Eso es
  correcto en Go, no es duplicación.

🧠 **`var _ catalog.Repository = (*CatalogRepo)(nil)`** es un chequeo en tiempo de compilación: si
`CatalogRepo` no cumple la interfaz, no compila.

🧠 **`toDomain(row)`** traduce el struct generado por sqlc a tu `movie.Movie`. Parece burocracia.
La razón: si mañana renombrás una columna, el cambio queda contenido en esa función.

---

## Paso 13 — Gin y el manejo de errores

🎯 Entrar al framework, y resolver de entrada su punto débil.

```bash
go get github.com/gin-gonic/gin
mkdir internal/apierr internal/api
```

`internal/apierr/apierr.go`:
```go
package apierr

import (
	"errors"
	"net/http"
)

type Error struct {
	Code    int
	Message string
	Field   string
	Err     error   // el original, para el log. NUNCA se expone.
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Err }

func BadRequest(msg string, err error) *Error {
	return &Error{Code: http.StatusBadRequest, Message: msg, Err: err}
}

func NotFound(msg string) *Error {
	return &Error{Code: http.StatusNotFound, Message: msg}
}

func Validation(field, msg string) *Error {
	return &Error{Code: http.StatusUnprocessableEntity, Message: msg, Field: field}
}

func Internal(err error) *Error {
	return &Error{Code: http.StatusInternalServerError, Message: "internal error", Err: err}
}

// From traduce errores de dominio a errores HTTP.
func From(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	switch {
	case errors.Is(err, movie.ErrNotFound):
		return NotFound("movie not found")
	case errors.Is(err, roulette.ErrNoCandidates):
		return &Error{Code: http.StatusNotFound, Message: "no movie matches the filters"}
	default:
		return Internal(err)
	}
}
```

`internal/api/middleware.go`:
```go
func ErrorHandler(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()                    // ← primero corre el handler

		if len(c.Errors) == 0 {
			return
		}

		e := apierr.From(c.Errors[0].Err)

		if e.Code >= 500 {
			log.Error("request failed",
				"path", c.Request.URL.Path,
				"err", e.Err,                          // detalle al log
				"request_id", c.GetString("request_id"))
		}

		body := gin.H{"error": e.Message}              // mensaje limpio al cliente
		if e.Field != "" {
			body["field"] = e.Field
		}
		c.JSON(e.Code, body)
	}
}
```

🐘 **Esto es tu `app/Exceptions/Handler.php`.** En Laravel el `render()` traduce excepciones a
respuestas HTTP. Acá `From()` hace lo mismo y `ErrorHandler` lo aplica.

🧠 **Por qué hace falta:** a diferencia de Laravel (o de Echo), en Gin el handler no devuelve nada.
Si escribís `c.JSON(400, ...)` y olvidás el `return`, el código sigue y escribís la respuesta dos
veces. Con este patrón los handlers hacen `c.Error(err); return` y todo el mapeo vive en un lugar.

🧠 **`c.Next()` divide el middleware en dos:** lo de arriba corre a la ida, lo de abajo a la
vuelta. Es como el `$next($request)` de un middleware de Laravel — mismo concepto exacto.

---

## Paso 14 — El servidor y las rutas

`internal/api/server.go`:
```go
func (s *Server) setupRouter() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()                    // NO gin.Default()
	r.Use(RequestID())
	r.Use(Logger(s.log))
	r.Use(gin.Recovery())
	r.Use(ErrorHandler(s.log))

	r.GET("/healthz", s.health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/movies",     s.searchMovies)
		v1.GET("/movies/:id", s.getMovie)

		v1.GET("/lists",        s.listLists)
		v1.POST("/lists",       s.createList)
		v1.GET("/lists/:id",    s.getList)
		v1.DELETE("/lists/:id", s.deleteList)

		v1.POST("/lists/:id/movies",                 s.addToList)
		v1.DELETE("/lists/:id/movies/:movieID",      s.removeFromList)
		v1.PATCH("/lists/:id/movies/:movieID/watch", s.markWatched)

		v1.POST("/lists/:id/spin", s.spin)
		v1.PATCH("/spins/:spinID", s.respondToSpin)
	}

	s.router = r
}

func (s *Server) Handler() http.Handler { return s.router }
```

🐘 **Es tu `routes/api.php`.** `r.Group("/api/v1")` es `Route::prefix('api/v1')->group(...)`. Los
middleware globales son el `$middleware` del Kernel.

🧠 **`gin.New()` y no `gin.Default()`:** `Default()` trae su propio logger que no controlás.
Preferís middleware explícito.

🧠 **`Handler()` devuelve `http.Handler`:** Gin es *solo el router*. El servidor es `http.Server`
de la stdlib. Por eso `httptest` funciona sin adaptadores, y por eso migrar a Echo después sería
medio día.

---

## Paso 15 — Handlers, DTOs y validación

`internal/api/dto.go`:
```go
type CreateListRequest struct {
	Name        string `json:"name"        binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type SpinRequest struct {
	Genre     string  `json:"genre"      binding:"omitempty,max=50"`
	YearFrom  int     `json:"year_from"  binding:"omitempty,gte=1888,lte=2100"`
	YearTo    int     `json:"year_to"    binding:"omitempty,gte=1888,lte=2100"`
	MinRating float64 `json:"min_rating" binding:"omitempty,gte=0,lte=5"`
	SkipLast  int     `json:"skip_last"  binding:"omitempty,gte=0,lte=50"`
	Dramatic  bool    `json:"dramatic"`
}
```

🐘 **Esto es tu FormRequest.** Los tags `binding:` son las `rules()`:
```php
public function rules(): array {
    return ['name' => 'required|min:1|max:100'];
}
```
Prácticamente la misma sintaxis. Es lo que más rápido te va a resultar familiar de Gin.

`internal/api/handlers_list.go`:
```go
func (s *Server) createList(c *gin.Context) {
	var req CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(translateValidation(err))
		return
	}

	l, err := s.lists.Create(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, l)
}
```

🐘 **Comparación directa:**
```php
public function store(CreateListRequest $r) {          // Laravel
    $l = $this->service->create($r->name, $r->description);
    return response()->json($l, 201);
}
```

Casi idéntico. Las diferencias: el binding es explícito (`ShouldBindJSON`) en vez de inyectado por
el container, y el error se registra en vez de tirarse.

🧠 **La regla:** el handler solo traduce HTTP ↔ dominio. Máximo ~15 líneas. Toda la lógica va al
servicio. Es la misma disciplina de "controllers finos" de Laravel, pero acá no hay Eloquent
tentándote a hacer queries en el controller.

⚠️ **Trampa que te va a pasar:** si un slice es `nil`, Gin lo serializa como `null`, no `[]`. Tu
Vue explota. Inicializá siempre: `results := []movie.Movie{}` en el repositorio.

---

## Paso 16 — El wiring

`cmd/api/main.go`:
```go
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		log.Error("invalid config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// infraestructura
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// adaptadores de salida
	queries     := db.New(pool)
	catalogRepo := postgres.NewCatalogRepo(queries)
	listRepo    := postgres.NewListRepo(queries)

	// casos de uso
	catalogSvc := catalog.New(catalogRepo)
	listSvc    := lists.New(listRepo, catalogRepo)

	// adaptador de entrada
	srv := api.New(catalogSvc, listSvc, log)

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("listening", "port", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	log.Info("bye")
}
```

🐘 **Esto reemplaza tu `AppServiceProvider`.** En Laravel:
```php
$this->app->bind(MovieRepositoryInterface::class, EloquentMovieRepository::class);
```
y el container resuelve todo por type-hint.

🧠 **En Go no hay container y no lo necesitás.** Todo el grafo de dependencias se lee de arriba
abajo en 20 líneas. Nada es implícito, nada se resuelve por reflexión, nada falla en runtime porque
olvidaste un binding.

🧠 **Existen `wire`, `fx` y `dig`.** No los uses. Cuando `main.go` pase de ~60 líneas, extraé un
`internal/app/app.go` con `Build(cfg) (*Server, func(), error)`. Nada más.

🧠 **Graceful shutdown** no tiene equivalente en Laravel (php-fpm mata el proceso). Acá esperás a
que terminen los requests en curso. Es lo que separa un juguete de algo que podés deployar.

**🎉 Punto de control:** `go run ./cmd/api` y tenés API funcionando contra 9.742 películas.

---

# ETAPA 4 — La app real (pasos 17-20)

## Paso 17 — La ruleta completa

`internal/lists/service.go` — acá va la lógica de verdad:

1. Traer candidatas de la lista con filtros aplicados **en SQL** (que filtre Postgres, no Go)
2. Traer los últimos N spins y excluir esas películas
3. `roulette.SpinWith(candidates, filter, rand)`
4. Guardar el spin en la tabla
5. Si `dramatic`, elegir 6-8 más para la animación

Respuesta:
```json
{
  "spin_id": "uuid",
  "winner": { "id": 296, "title": "Pulp Fiction" },
  "reel": [ {...}, {...}, {...} ],
  "total_candidates": 47
}
```

🧠 **Por qué filtrar en SQL y no en Go:** traer 500 películas para descartar 490 es desperdicio de
red y memoria. Que Postgres use sus índices. Esta decisión es la misma en Laravel — la diferencia
es que Eloquent te lo hace fácil de hacer mal con `->get()->filter()`.

---

## Paso 18 — Tests de la API

```go
func TestCreateList(t *testing.T) {
	srv := setupTestServer(t)      // con repos en memoria, sin Postgres

	req := httptest.NewRequest("POST", "/api/v1/lists",
		strings.NewReader(`{"name":"Classic horror"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("got %d want 201: %s", w.Code, w.Body.String())
	}
}
```

🐘 **En Laravel** sería `$this->postJson('/api/v1/lists', [...])->assertStatus(201)`, y por defecto
usa una base de datos real con `RefreshDatabase`.

🧠 **Acá el test no toca la base.** Los repos son fakes en memoria, gracias a las interfaces del
paso 12. Corre en milisegundos. Para los tests que sí necesitan Postgres real, usás
`testcontainers-go` — pero deberían ser pocos.

**La regla:** 80% de tus tests sin base de datos. Si necesitás Postgres para testear una regla de
negocio, la regla está en el lugar equivocado.

---

## Paso 19 — Frontend Vue embebido

```go
//go:embed all:web/dist
var frontend embed.FS

func (s *Server) setupFrontend(r *gin.Engine) {
	dist, _ := fs.Sub(frontend, "web/dist")
	r.StaticFS("/assets", http.FS(dist))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found"})
			return
		}
		c.FileFromFS("index.html", http.FS(dist))   // SPA fallback
	})
}
```

🐘 **En Laravel** servirías desde `public/` con Vite y nginx apuntando ahí.

🧠 **Acá el frontend queda DENTRO del binario.** `go build` produce **un archivo** con la API, el
frontend y todo. Lo copiás al VPS y anda. Sin nginx, sin `npm install` en el servidor, sin versión
de Node que no coincide.

Esto es probablemente lo que más vas a valorar de Go para proyectos personales.

---

## Paso 20 — Docker, CI y deploy

```dockerfile
FROM node:22 AS front
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

FROM golang:1.24 AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=front /web/dist ./web/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/api

FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
ENTRYPOINT ["/app"]
```

🧠 Imagen final: ~20 MB, con frontend incluido. Compará con tu imagen de Laravel (PHP-FPM +
extensiones + nginx + vendor, fácil 400 MB).

GitHub Actions:
```yaml
- run: go vet ./...
- run: golangci-lint run
- run: go test -race ./...
```

---

# Estructura final

```
movie-roulette/
├── cmd/
│   ├── api/main.go              # el servidor          ← public/index.php
│   ├── import/main.go           # importa MovieLens    ← artisan command
│   └── seed/main.go             # datos de prueba      ← DatabaseSeeder
├── internal/
│   ├── movie/                   # dominio puro         ← Model (solo la parte de datos)
│   ├── roulette/                # dominio puro         ← ninguno, es lógica pura
│   ├── catalog/                 # caso de uso          ← app/Services
│   ├── lists/                   # caso de uso          ← app/Services
│   ├── storage/
│   │   ├── postgres/            # repos               ← app/Repositories
│   │   ├── memory/              # fakes para tests    ← ninguno
│   │   └── db/                  # generado por sqlc   ← Eloquent (la parte de queries)
│   ├── api/
│   │   ├── server.go            # rutas               ← routes/api.php
│   │   ├── handlers_*.go        # handlers            ← app/Http/Controllers
│   │   ├── dto.go               # requests            ← app/Http/Requests
│   │   └── middleware.go        # middleware          ← app/Http/Middleware
│   ├── apierr/                  # errores HTTP        ← app/Exceptions/Handler
│   ├── importer/                # lectura de CSV
│   └── config/                  # config              ← config/ + .env
├── migrations/                  # SQL                 ← database/migrations
├── db/query.sql                 # queries para sqlc
├── web/                         # Vue                 ← resources/js
├── docs/specs/                  # specs del flujo
└── CLAUDE.md
```

**Diferencia clave de organización:** Laravel organiza **por capa técnica** (todos los models
juntos, todos los controllers juntos). Go organiza **por dominio**: todo lo de lists vive en
`internal/lists/`. Es lo más contraintuitivo del cambio y es lo correcto acá.

---

# Los 10 errores que vas a cometer (viniendo de Laravel)

1. **Modificar dentro de un `for _, x := range`** — es una copia. Usá el índice.
2. **Ignorar un error** porque "no va a fallar". El linter te va a insistir. Hacele caso.
3. **Crear `internal/models/`** con todo junto. Organizá por dominio.
4. **Poner la interfaz junto a la implementación.** Va junto al consumidor.
5. **Devolver el struct de sqlc como respuesta HTTP.** Usá un DTO.
6. **Buscar un container de DI.** No lo necesitás.
7. **`c.JSON()` sin `return`** después. El bug #1 de Gin.
8. **Slice `nil` serializado como `null`.** Inicializá vacío.
9. **Lógica de negocio en el handler.** Máximo 15 líneas.
10. **Copiar el "Standard Go Project Layout"** de 15 carpetas. No es oficial, empezá plano.

---

# Convención de nombres

Todo el código en inglés. Los comentarios pueden ir en español si te resulta más natural,
pero identificadores, columnas de DB y campos JSON siempre en inglés.

| Concepto | Identificador | Columna DB | Campo JSON |
|---|---|---|---|
| Película | `Movie`, `movies` | `movies` | `movie` |
| Lista | `List`, `lists` | `lists` | `list` |
| Giro | `Spin`, `spins` | `spins` | `spin` |
| Favorito | `Favorite` | `favorites` | `favorite` |
| Puntaje | `Rating` | `rating` | `rating` |
| Visto | `Watched` | `watched` | `watched` |
| Filtro | `Filter` | — | `filter` |

🧠 **Por qué importa en Go específicamente:** los tags JSON y los nombres de columna tienen que
coincidir con lo que espera el frontend y la base. Mezclar idiomas ahí genera bugs silenciosos
(`json:"pelicula"` contra un Vue que espera `movie`). Y sqlc genera los nombres de los campos
Go a partir de las columnas SQL, así que si la columna se llama `pelicula`, el struct generado
va a tener un campo `Pelicula`. La consistencia se propaga sola.

---

# Cómo saber que aprendiste

- [ ] Escribís un handler de Gin de memoria
- [ ] Explicás por qué el error es un valor y no una excepción
- [ ] Decidís puntero vs valor en un receiver sin dudar
- [ ] Detectás un data race leyendo código
- [ ] Sabés cuándo **no** usar goroutines, y lo justificás con una medición
- [ ] Definís una interfaz del lado correcto sin pensarlo
- [ ] Explicás por qué `main.go` reemplaza al Service Container

---

# Notas finales

- **Los pasos 1-8 son el 80% del aprendizaje** aunque sean el 20% del código. No los apures.
- **El paso 8 es EL paso.** Los tres experimentos con `-race` valen más que cualquier tutorial.
- **El paso 12 es la clave arquitectónica.** Si entendés dónde va la interfaz, entendiste Go.
- **Cada paso es un commit.** Desde el 11, cada paso es una spec + un PR contra vos mismo, así el
  proyecto también es el experimento de spec-driven.
- **Cuando no entiendas algo, escribí un test que lo demuestre.** Go compila en segundos:
  experimentar sale más barato que teorizar.
- **MovieLens es libre para uso personal, no comercial.** Si lo publicás como producto, migrá a
  TMDB usando los `tmdb_id` que ya guardaste.
