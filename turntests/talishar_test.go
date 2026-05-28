package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests Talishar's rust-counter lifecycle over four turns: swing for 4 on turns 1-3
// (accumulating a rust counter each swing), self-destruct at 3 counters in turn 3's end
// phase, then deal 0 on turn 4 with no weapon left to swing. The deck is all blue pitch so
// the only value-generating play each turn is the Talishar swing (pitch one blue to fund the
// 2-cost {r}{r} activation) — resources can't arsenal or attack, so the line is forced.
func TestTalishar_SelfDestructsAfterThreeSwings(t *testing.T) {
	filler := testutils.FakeBlueResource()
	deckCards := make([]deck.Card, 40)
	for i := range deckCards {
		deckCards[i] = filler
	}
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.Talishar{}}, deckCards)
	hand1 := []card.Card{filler, filler, filler, filler}

	s1, s2, s3, s4 := sim.EvalFourTurnsForTesting(d,
		gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand1)

	for i, s := range []sim.TurnSummary{s1, s2, s3} {
		if s.Value != 4 {
			t.Errorf("turn %d Value = %d, want 4 (Talishar swing)\nBestLine: %s",
				i+1, s.Value, formatBestLine(s.BestLine))
		}
	}
	if s4.Value != 0 {
		t.Errorf("turn 4 Value = %d, want 0 (Talishar destroyed at 3 rust counters)\nBestLine: %s",
			s4.Value, formatBestLine(s4.BestLine))
	}
	if n := len(s3.State.Weapons()); n != 0 {
		t.Errorf("turn 3 equipped weapons = %d, want 0 (Talishar self-destructed end of turn 3)", n)
	}
	if !graveyardContains(s3.State.Graveyard(), weapons.Talishar{}) {
		t.Errorf("turn 3 graveyard = %v, want Talishar's card present after self-destruct", s3.State.Graveyard())
	}
}
