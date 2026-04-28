package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// heroesByName resolves a Hero's display name to the concrete Hero value. Built once at init.
// Serialization packages (deckio, fabrary) look up heroes here so the list is maintained in
// one place as new heroes come online.
var heroesByName = func() map[string]sim.Hero {
	all := []sim.Hero{
		heroes.Viserai{},
	}
	m := make(map[string]sim.Hero, len(all))
	for _, h := range all {
		m[h.Name()] = h
	}
	return m
}()

// HeroByName returns the Hero for the given display name. Returns (nil, false) when the name
// isn't registered — callers surface that to the user rather than falling back to a default.
func HeroByName(name string) (sim.Hero, bool) {
	h, ok := heroesByName[name]
	return h, ok
}
