package registry

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that v2/deck.Random builds a legal deck against the production card pool through
// Registry.
func TestRegistry_DrivesDeckRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, Registry{})
	if d.Size() != 40 {
		t.Errorf("len(Cards) = %d, want 40", d.Size())
	}
	if n := len(d.Weapons); n == 0 || n > 2 {
		t.Errorf("len(Weapons) = %d, want 1 or 2", n)
	}
}
