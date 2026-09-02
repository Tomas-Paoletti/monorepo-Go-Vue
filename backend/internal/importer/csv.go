package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Tomas-Paoletti/roulette-movies/internal/movie"
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