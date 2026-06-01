package registry

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
)

// Tests that internal/deck.Random builds a legal deck against the production card pool through
// Registry.
func TestRegistry_DrivesDeckRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(heroes.Viserai, format.SilverAge, 40, 2, rng, Registry{})
	if d.Size() != 40 {
		t.Errorf("len(Cards) = %d, want 40", d.Size())
	}
	if n := len(d.Weapons); n == 0 || n > 2 {
		t.Errorf("len(Weapons) = %d, want 1 or 2", n)
	}
}
