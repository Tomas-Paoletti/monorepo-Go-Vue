package movie

import "testing"

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