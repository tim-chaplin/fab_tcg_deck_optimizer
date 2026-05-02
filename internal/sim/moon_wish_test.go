package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Moon Wish's printed text shuffles the deck after the on-hit Sun Kiss tutor, so assertions
// here avoid pinning specific deck positions or which card lands in arsenal off a post-tutor
// DrawOne. Valid checks: Value, total copies of each card across next-turn Hand/Deck/Arsenal,
// and which BestLine roles got assigned.

// TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck is the canonical
// no-go-again scenario for Moon Wish run end-to-end through one turn:
//
//   - Hand: Moon Wish [Y], Weeping Battleground (Red — a red Defense Reaction).
//   - Deck: Sun Kiss [R] plus filler.
//
// Expected line: Moon Wish plays via its alt cost — Weeping Battleground is the consumed
// hand card and gets returned to the deck. Moon Wish hits, tutors Sun Kiss out of the deck,
// and (without go-again) appends Sun Kiss to s.Hand for the post-hoc Arsenal promotion.
//
// Cross-turn assertions:
//   - turn-1 Value = 4 (Moon Wish base only; Sun Kiss tutored, not played).
//   - Sun Kiss appears in NO turn-2 surface (Hand / Deck / Arsenal) directly — but is
//     promoted to Arsenal as the only Held candidate after alt cost consumed the DR.
//   - Weeping Battleground exists exactly once across turn-2 surfaces — the alt-cost path
//     must not drop the card on the floor.
func TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck(t *testing.T) {
	deckCards := []Card{
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := d.EvalOneTurnForTesting(Matchup{IncomingDamage: 0}, nil, []Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.Value != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; Sun Kiss tutored, not played)",
			state.Value)
	}
	if state.StartOfNextTurnArsenal == nil || state.StartOfNextTurnArsenal.ID() != ids.SunKissRed {
		t.Errorf("StartOfNextTurnArsenal = %v, want Sun Kiss [R] (post-hoc promoted from State.Hand)",
			state.StartOfNextTurnArsenal)
	}
	if got := countAcrossSurfaces(state, ids.SunKissRed); got != 1 {
		t.Errorf("Sun Kiss [R] total across turn-2 Hand/Deck/Arsenal = %d, want 1 (in Arsenal)",
			got)
	}
	if got := countAcrossSurfaces(state, ids.WeepingBattlegroundRed); got != 1 {
		t.Errorf("Weeping Battleground [R] total across turn-2 surfaces = %d, want 1 "+
			"(alt cost returned it to deck — should still exist somewhere)",
			got)
	}
}

// TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss is the negative-tutor variant:
// alt cost still returns Weeping Battleground to the deck, but with no Sun Kiss in the deck
// the tutor finds nothing and bails. No drawn card lands; arsenal stays empty (the only
// Held candidate was consumed by alt cost).
func TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	deckCards := []Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := d.EvalOneTurnForTesting(Matchup{IncomingDamage: 0}, nil, []Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.Value != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; tutor fizzles)",
			state.Value)
	}
	if state.StartOfNextTurnArsenal != nil {
		t.Errorf("StartOfNextTurnArsenal = %v, want nil (DR was the only Held; alt cost consumed it)",
			state.StartOfNextTurnArsenal)
	}
	if got := countAcrossSurfaces(state, ids.WeepingBattlegroundRed); got != 1 {
		t.Errorf("Weeping Battleground [R] total across turn-2 surfaces = %d, want 1 "+
			"(alt cost returned it to deck even when the tutor fizzled)",
			got)
	}
}

// TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss covers the go-again branch:
// Flying High [R] plays first in the chain and grants Moon Wish go-again. Moon Wish then
// alt-costs the Held DR, hits, tutors Sun Kiss, sees self.EffectiveGoAgain() = true, and
// plays Sun Kiss immediately. Sun Kiss heads to the graveyard.
//
// Sun Kiss sits at deck position 2 (not the top) so the buf-mutation path is exercised: a
// blind head-advance would consume buf[0] instead of patching out the specific Sun Kiss
// slot, leaving Sun Kiss to surface in turn-2 Hand or Deck.
//
// Cross-turn assertions:
//   - turn-1 Value = 7 (Moon Wish 4 + Sun Kiss 3) — confirms Sun Kiss actually played.
//   - Graveyard contains Sun Kiss — proves it was both tutored AND played, not stuck in the
//     deck or some intermediate state.
//   - StartOfNextTurnArsenal is non-nil — Sun Kiss's DrawOne pulled some card off the (per
//     printed rules, shuffled) deck top; the specific identity is random so we don't pin it.
func TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss(t *testing.T) {
	// Sun Kiss at index 2 (not 0) so a blind head++ would consume the wrong slot; verifies
	// the buf-removal logic patches out the specific tutored card.
	deckCards := []Card{
		testutils.RedAttack{}, testutils.RedAttack{},
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := d.EvalOneTurnForTesting(Matchup{IncomingDamage: 0}, nil, []Card{
		cards.FlyingHighRed{},
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.Value != 7 {
		t.Errorf("turn-1 Value = %d, want 7 (Moon Wish 4 + Sun Kiss 3 via Flying High go-again)",
			state.Value)
	}
	skInGraveyard := false
	for _, c := range state.Graveyard {
		if c.ID() == ids.SunKissRed {
			skInGraveyard = true
			break
		}
	}
	if !skInGraveyard {
		t.Errorf("Sun Kiss [R] not in turn-1 Graveyard %v; want it there (tutored and played)",
			testutils.CardNames(state.Graveyard))
	}
	if got := countAcrossSurfaces(state, ids.SunKissRed); got != 0 {
		t.Errorf("Sun Kiss [R] total across turn-2 surfaces = %d, want 0 (it's in the graveyard)",
			got)
	}
	if state.StartOfNextTurnArsenal == nil {
		t.Error("StartOfNextTurnArsenal = nil; want any card (Sun Kiss's DrawOne pulled one card into " +
			"State.Hand → Arsenal promotion is the only candidate)")
	}
}

// countAcrossSurfaces totals occurrences of id across the start-of-next-turn Hand, Deck, and
// Arsenal. Asserts "exists / doesn't exist" without pinning a specific position.
func countAcrossSurfaces(state TurnStartState, id ids.CardID) int {
	n := 0
	for _, c := range state.StartOfNextTurnHand {
		if c.ID() == id {
			n++
		}
	}
	for _, c := range state.StartOfNextTurnDeck {
		if c.ID() == id {
			n++
		}
	}
	if state.StartOfNextTurnArsenal != nil && state.StartOfNextTurnArsenal.ID() == id {
		n++
	}
	return n
}
