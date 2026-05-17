package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
)

// Moon Wish's on-hit tutor shuffles the deck, so assertions only check Value, total copies
// across next-turn Hand/Deck/Arsenal, and BestLine roles — never specific deck positions.

// Tests Moon Wish alt-cost: hand DR returned to deck, on-hit tutors Sun Kiss; without
// go-again, Sun Kiss is post-hoc-promoted to Arsenal as the only Held candidate.
func TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck(t *testing.T) {
	deckCards := []deck.Card{
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{
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

// Tests Moon Wish alt-cost when the deck has no Sun Kiss: tutor fizzles, alt cost still
// recycles the DR, and arsenal stays empty.
func TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	deckCards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{
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

// Tests the go-again branch: Flying High grants Moon Wish go again so the tutored Sun Kiss
// plays the same turn and lands in graveyard. Sun Kiss is at deck index 2 (not top) so the
// buf-removal path is exercised — a blind head-advance would consume the wrong slot.
func TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss(t *testing.T) {
	deckCards := []deck.Card{
		testutils.RedAttack{}, testutils.RedAttack{},
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{
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

// countAcrossSurfaces totals occurrences of the printing across the start-of-next-turn
// Hand, Deck, and Arsenal — keyed by DisplayName since that's what NameCounts surfaces.
// Asserts "exists / doesn't exist" without pinning a specific position.
func countAcrossSurfaces(state sim.TurnStartState, id ids.CardID) int {
	name := registry.GetCard(id).DisplayName()
	n := state.StartOfNextTurnDeck.NameCounts()[name]
	for _, c := range state.StartOfNextTurnHand {
		if c.DisplayName() == name {
			n++
		}
	}
	if state.StartOfNextTurnArsenal != nil && state.StartOfNextTurnArsenal.DisplayName() == name {
		n++
	}
	return n
}
