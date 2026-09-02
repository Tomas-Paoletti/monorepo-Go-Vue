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

