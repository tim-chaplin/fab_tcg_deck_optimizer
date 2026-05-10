package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"

// heroesByName resolves a Hero's display name to the concrete Hero value. Single source of
// truth for the hero roster; serialization packages look up heroes here.
var heroesByName = func() map[string]Hero {
	all := []Hero{
		heroes.Viserai{},
	}
	m := make(map[string]Hero, len(all))
	for _, h := range all {
		m[h.Name()] = h
	}
	return m
}()

// HeroByName returns the Hero for the given display name. Returns (nil, false) when the
// name isn't registered; callers surface that to the user rather than defaulting.
func HeroByName(name string) (Hero, bool) {
	h, ok := heroesByName[name]
	return h, ok
}
