package deck

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestEvalOneTurn_MidTurnDrawCarriesIntoNextTurn pins the Phase 0 invariant that the card
// consumed by a mid-turn DrawOne lands in the next turn's dealt hand.
//
// Layout (consumed in source order — no shuffle):
//   - positions 0..3: turn 1's hand — Snatch Red (cost 0, attack 4, on-hit DrawOne) plus three
//     Blue attacks that the optimiser spends to reach Value = 6 (pitch 1 Blue → chain Blue +
//     Blue + Snatch for 1 + 1 + 4 damage).
//   - position 4: the "beacon" — a fake.RedAttack whose ID appears nowhere else in the deck.
//     Snatch's on-hit DrawOne consumes this mid-turn 1.
//   - positions 5..7: three Blue attacks that make up turn 2's refill.
//   - positions 8..9: Yellow attacks. Turn 2's refill should never reach these — if they show
//     up in turn 2's hand the sim over-drew past the Phase 0 budget.
//
// Turn 1 ends with the three Blues partitioned as 1 Pitch + 2 Attack (the pitched Blue recycles
// to the bottom of the deck). No card is Held, so the arsenal stays empty. At the start of
// turn 2 the sim's refill draws the top handSize cards from what remains — the beacon first
// (since the mid-turn draw displaced an end-of-turn refill draw, per Phase 0), then the three
// filler Blues. The recycled pitched Blue now sits at the bottom of the deck, after the two
// Yellows.
func TestEvalOneTurn_MidTurnDrawCarriesIntoNextTurn(t *testing.T) {
	beacon := fake.RedAttack{}
	deckCards := []card.Card{
		generic.SnatchRed{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0)

	wantHand := []card.Card{
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (beacon should land at slot 0 — it's the card Snatch drew mid-turn 1 — followed by three fresh Blues; a Yellow here means the sim over-drew)", state.Hand, wantHand)
	}

	if state.ArsenalCard != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (turn 1 has no Held card to promote)", state.ArsenalCard)
	}

	// Remaining deck: two untouched Yellows from source positions 8 and 9, then the pitched
	// Blue recycled to the bottom on turn 1.
	wantDeck := []card.Card{
		fake.YellowAttack{},
		fake.YellowAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (nothing on turn 1 creates runechants)", state.RunechantCarryover)
	}
}
