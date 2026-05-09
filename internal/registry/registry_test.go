package registry_test

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that v2/deck.Random builds a legal deck against the production card pool through
// registry.Registry. Confirms the package's LegalCards + LegalWeapons methods satisfy
// deck.Registry directly — no separate adapter type required.
func TestRegistry_DrivesDeckRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(heroes.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	if len(d.Cards) != 40 {
		t.Errorf("len(Cards) = %d, want 40", len(d.Cards))
	}
	if n := len(d.Weapons); n == 0 || n > 2 {
		t.Errorf("len(Weapons) = %d, want 1 or 2", n)
	}
}
