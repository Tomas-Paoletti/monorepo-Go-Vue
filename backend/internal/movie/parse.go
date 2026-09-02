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

	if m := yearRe.FindStringSubmatch(rest); m != nil {
		rest = m[1]
		year, _ = strconv.Atoi(m[2])
	}
	if m := altRe.FindStringSubmatch(rest); m != nil {
		rest = m[1]
		alt = m[2]
	}
	if m := articleRe.FindStringSubmatch(rest); m != nil {
		rest = m[2] + " " + m[1]
	}

	return strings.TrimSpace(rest), year, alt
}