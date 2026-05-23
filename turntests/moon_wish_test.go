package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Moon Wish's on-hit tutor shuffles the deck, so assertions only check Value, total copies
// across next-turn Hand/Deck/Arsenal, and BestLine roles — never specific deck positions.

// Tests Moon Wish alt-cost: hand DR returned to deck, on-hit tutors Sun Kiss; without
// go-again, Sun Kiss is post-hoc-promoted to Arsenal as the only Held candidate.
func TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck(t *testing.T) {
	deckCards := []deck.Card{
		cards.SunKissRed{},
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if summary.Value != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; Sun Kiss tutored, not played)",
			summary.Value)
	}
	if summary.State.Arsenal() == nil || summary.State.Arsenal().ID() != ids.SunKissRed {
		t.Errorf("Arsenal() = %v, want Sun Kiss [R] (post-hoc promoted from hand)",
			summary.State.Arsenal())
	}
	if got := countAcrossSurfaces(summary.State, ids.SunKissRed); got != 1 {
		t.Errorf("Sun Kiss [R] total across turn-2 Hand/Deck/Arsenal = %d, want 1 (in Arsenal)",
			got)
	}
	if got := countAcrossSurfaces(summary.State, ids.WeepingBattlegroundRed); got != 1 {
		t.Errorf("Weeping Battleground [R] total across turn-2 surfaces = %d, want 1 "+
			"(alt cost returned it to deck — should still exist somewhere)",
			got)
	}
}

// Tests Moon Wish alt-cost when the deck has no Sun Kiss: tutor fizzles, alt cost still
// recycles the DR, and arsenal stays empty.
func TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	deckCards := []deck.Card{
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if summary.Value != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; tutor fizzles)",
			summary.Value)
	}
	if summary.State.Arsenal() != nil {
		t.Errorf("Arsenal() = %v, want nil (DR was the only Held; alt cost consumed it)",
			summary.State.Arsenal())
	}
	if got := countAcrossSurfaces(summary.State, ids.WeepingBattlegroundRed); got != 1 {
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
		testutils.FakeRedAttack(), testutils.FakeRedAttack(),
		cards.SunKissRed{},
		testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(),
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{
		cards.FlyingHighRed{},
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if summary.Value != 7 {
		t.Errorf("turn-1 Value = %d, want 7 (Moon Wish 4 + Sun Kiss 3 via Flying High go-again)",
			summary.Value)
	}
	skInGraveyard := false
	for _, c := range summary.State.Graveyard() {
		if c.ID() == ids.SunKissRed {
			skInGraveyard = true
			break
		}
	}
	if !skInGraveyard {
		t.Errorf("Sun Kiss [R] not in turn-1 Graveyard %v; want it there (tutored and played)",
			testutils.CardNamesSim(summary.State.Graveyard()))
	}
	if got := countAcrossSurfaces(summary.State, ids.SunKissRed); got != 0 {
		t.Errorf("Sun Kiss [R] total across turn-2 surfaces = %d, want 0 (it's in the graveyard)",
			got)
	}
	if summary.State.Arsenal() == nil {
		t.Error("Arsenal() = nil; want any card (Sun Kiss's DrawOne pulled one card; " +
			"arsenal promotion is the only candidate)")
	}
}

// countAcrossSurfaces totals occurrences of the printing across the start-of-next-turn
// Hand, Deck, and Arsenal — keyed by DisplayName since that's what NameCounts surfaces.
// Asserts "exists / doesn't exist" without pinning a specific position.
func countAcrossSurfaces(state *gameengine.GameState, id ids.CardID) int {
	name := registry.GetCard(id).DisplayName()
	n := state.Deck().NameCounts()[name]
	for _, c := range state.Hand() {
		if c.DisplayName() == name {
			n++
		}
	}
	if state.Arsenal() != nil && state.Arsenal().DisplayName() == name {
		n++
	}
	return n
}
