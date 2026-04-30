package sim_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that printed Go again short-circuits to true without invoking the probe.
func TestHasGoAgainHeuristic_PrintedGoAgain(t *testing.T) {
	c := testutils.NewStubCard("printed").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeAction)).
		WithGoAgain()
	if !HasGoAgainHeuristic(c) {
		t.Errorf("HasGoAgainHeuristic(printed go-again) = false, want true")
	}
}

// Tests that a card whose Play never grants Go again returns false from the probe.
func TestHasGoAgainHeuristic_PrintedNoGoAgainNoConditional(t *testing.T) {
	// Plain stub: no GoAgain printed, no conditional grant in Play (zero-method body).
	c := testutils.NewStubCard("plain").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeAction))
	if HasGoAgainHeuristic(c) {
		t.Errorf("HasGoAgainHeuristic(plain) = true, want false")
	}
}

// Tests that Runerager Swarm — printed GoAgain() == false, but reliably granted via
// "if you've played or created an aura this turn" in a Viserai deck — registers as
// having Go again under the probe (which seeds AuraCreated = true).
func TestHasGoAgainHeuristic_RuneragerSwarmConditionalGrant(t *testing.T) {
	c := cards.RuneragerSwarmRed{}
	if c.GoAgain() {
		t.Skip("Runerager Swarm gained printed Go again — switch this test to a different conditional-grant card or drop it")
	}
	if !HasGoAgainHeuristic(c) {
		t.Errorf("HasGoAgainHeuristic(RuneragerSwarmRed) = false, want true (aura conditional should fire under probe)")
	}
}
