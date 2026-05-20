package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that Arcane Cussing, popped by an attack hitting during our turn, creates its N
// Runechants — the printed "when this leaves the arena during your turn" clause.
func TestArcaneCussing_HitOnOurTurnCreatesRunechants(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.ArcaneCussingRed{}, 3},
		{cards.ArcaneCussingYellow{}, 2},
		{cards.ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			CreateAuraFromCard(tc.c).
			Build()}
		ge.SetIsMyTurn(true)
		ge.FireTriggers(triggertype.Hit, testutils.AttackWithPower{Power: 4})
		if got := ge.RunechantCount(); got != tc.n {
			t.Errorf("%s: RunechantCount = %d, want %d (popped by a hit during our turn)", tc.c.Name(), got, tc.n)
		}
		if got := len(ge.Auras()); got != 1 {
			t.Errorf("%s: Auras = %d, want 1 (Arcane Cussing destroyed, the Runechant aura created)", tc.c.Name(), got)
		}
	}
}

// Tests that Arcane Cussing popped by DamageTaken outside our turn — the defense phase — is
// destroyed without creating Runechants. The payoff is "during your turn" only.
func TestArcaneCussing_DamageTakenOffOurTurnCreatesNoRunechants(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		Build()}
	ge.SetIsMyTurn(false)
	ge.FireTriggers(triggertype.DamageTaken, nil)
	if got := ge.RunechantCount(); got != 0 {
		t.Errorf("RunechantCount = %d, want 0 (popped off our turn — no payoff)", got)
	}
	if got := len(ge.Auras()); got != 0 {
		t.Errorf("Auras = %d, want 0 (Arcane Cussing destroyed, no Runechant aura)", got)
	}
}

// Tests the defense-phase DamageTaken fire end to end: a carried-over Arcane Cussing aura is
// destroyed — without paying out, since it isn't our turn — when incoming damage gets
// through unblocked.
func TestArcaneCussing_DamageTakenInDefensePhase(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.ArcaneCussingRed{}).
		SetIncomingDamage(7).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{testutils.BluePitch{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if got := summary.State.RunechantCount(); got != 0 {
		t.Fatalf("end-of-turn RunechantCount = %d, want 0 (destroyed by DamageTaken, not our turn)", got)
	}
	if got := len(summary.State.Auras()); got != 0 {
		t.Fatalf("end-of-turn Auras = %d, want 0 (Arcane Cussing destroyed in the defense phase)", got)
	}
}
