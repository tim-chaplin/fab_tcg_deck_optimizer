package fabrary

// Card-name conversion between the optimizer's canonical pitch suffix ("[R]") and
// fabrary's lowercase parenthesised form ("(red)"). Shared by both marshalling (which
// translates outbound) and unmarshalling (which translates inbound).

import "strings"

// pitchColors pairs the optimizer's canonical suffix ("[R]") with fabrary's lowercase form
// ("(red)"). One entry per color is enough — suffixes don't overlap.
var pitchColors = []struct{ canon, fabrary string }{
	{"[R]", "(red)"},
	{"[Y]", "(yellow)"},
	{"[B]", "(blue)"},
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
