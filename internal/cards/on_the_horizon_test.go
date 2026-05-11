package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that On the Horizon is typed as Block — not Defense Reaction — so the chain
// runner routes it through the plain-block path rather than calling its Play during the
// defender DR loop.
func TestOnTheHorizon_TypedAsBlock(t *testing.T) {
	cases := []card.Card{OnTheHorizonRed{}, OnTheHorizonYellow{}, OnTheHorizonBlue{}}
	for _, c := range cases {
		ts := c.Types()
		if !ts.Has(card.TypeBlock) {
			t.Errorf("%s: missing TypeBlock", c.Name())
		}
		if ts.IsDefenseReaction() {
			t.Errorf("%s: should not have TypeDefenseReaction", c.Name())
		}
	}
}
