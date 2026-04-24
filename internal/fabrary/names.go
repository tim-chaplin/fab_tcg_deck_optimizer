package fabrary

// Card-name conversion between the optimizer's canonical pitch suffix ("(Red)") and
// fabrary's lowercase form ("(red)"). Shared by both marshalling (which lowercases on the
// way out) and unmarshalling (which re-canonicalises on the way in).

import "strings"

// pitchColors pairs the optimizer's canonical suffix ("(Red)") with fabrary's lowercase form
// ("(red)"). One entry per color is enough — suffixes don't overlap.
var pitchColors = []struct{ canon, fabrary string }{
	{"(Red)", "(red)"},
	{"(Yellow)", "(yellow)"},
	{"(Blue)", "(blue)"},
}

func toFabraryCardName(s string) string {
	for _, p := range pitchColors {
		if strings.HasSuffix(s, p.canon) {
			return strings.TrimSuffix(s, p.canon) + p.fabrary
		}
	}
	return s
}

func fromFabraryCardName(s string) string {
	for _, p := range pitchColors {
		if strings.HasSuffix(s, p.fabrary) {
			return strings.TrimSuffix(s, p.fabrary) + p.canon
		}
	}
	return s
}
