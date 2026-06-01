package sim

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Tests that Best.Value is invariant to the input order of equal-(ID, FromArsenal)
// attackers, by comparing forward vs reversed Viserai hands.
func TestBestSequence_PermutationOrderInvariance(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		hands     = 200
	)
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for h := 0; h < hands; h++ {
		setupSeed := gen.Int63()
		baseline := deck.Random(heroes.Viserai, format.SilverAge, deckSize, maxCopies, rand.New(rand.NewSource(setupSeed)), registry.Registry{})
		d := baseline.Copy()
		d.Shuffle(rand.New(rand.NewSource(setupSeed ^ 0xdeadbeef)))

		intel := heroes.Viserai.Intelligence()
		if d.Size() < intel {
			continue
		}
		drawn := d.Draw(intel)
		hand := make([]card.Card, len(drawn))
		for i, c := range drawn {
			hand[i] = c.(card.Card)
		}
		reversed := make([]card.Card, len(hand))
		for i := range hand {
			reversed[len(hand)-1-i] = hand[i]
		}

		makeState := func() *gameengine.GameState {
			gs := gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build()
			gs.SetIncomingPhysicalDamage(incoming)
			return gs
		}
		forward := NewEvaluatorWithoutCache().Best(nil, hand, d.Copy(), makeState()).Value
		backward := NewEvaluatorWithoutCache().Best(nil, reversed, d.Copy(), makeState()).Value
		if forward != backward {
			t.Fatalf("setupSeed=%d: Best.Value depends on input hand order — forward=%d backward=%d hand=%v",
				setupSeed, forward, backward, displayHand(hand))
		}
	}
}

func displayHand(hand []card.Card) []string {
	out := make([]string, len(hand))
	for i, c := range hand {
		out[i] = c.DisplayName()
	}
	return out
}
