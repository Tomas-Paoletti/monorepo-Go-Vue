package movie

// Query son los filtros del buscador de catálogo. El cero de cada campo
// significa "no filtrar por eso". Vive en el dominio (no en storage ni en la
// API) para que ambos lados lo compartan sin acoplarse entre sí.
type Query struct {
	Text      string  // búsqueda por título (substring, case-insensitive)
	Genre     string  // género exacto
	YearFrom  int
	YearTo    int
	MinRating float64
	Limit     int // 0 => default del repo
	Offset    int
}
