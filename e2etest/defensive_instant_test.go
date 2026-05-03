package e2etest

// End-to-end tests for the DefensiveInstant marker. See sim.DefensiveInstant for the
// contract; docs/dev-standards.md for the rider wiring.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a free defensive instant prevents up to its Defense() value.
func TestDefensiveInstant_BrushOffRedAlone(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.BrushOffRed{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value; got != 3 {
		t.Fatalf("Value = %d, want 3 (Brush Off Red prevents 3 of 5)", got)
	}
}

// Tests the Yellow / Blue printings cap at their lower thresholds.
func TestDefensiveInstant_BrushOffYellowAndBlueAlone(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, []sim.Card{cards.BrushOffYellow{}}).Value; got != 2 {
		t.Fatalf("Yellow Value = %d, want 2", got)
	}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, []sim.Card{cards.BrushOffBlue{}}).Value; got != 1 {
		t.Fatalf("Blue Value = %d, want 1", got)
	}
}

// Tests Calming Breeze's collapsed-prevention model (3 events × 1 = 3 against the single
// IncomingDamage bucket).
func TestDefensiveInstant_CalmingBreezeAlone(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.CalmingBreezeRed{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value; got != 3 {
		t.Fatalf("Value = %d, want 3 (Calming Breeze 3 prevention)", got)
	}
}

// Tests that prevention caps at IncomingDamage — Oasis Respite Red has Defense 4 but only
// 1 incoming damage, so DealEffectiveDefense credits 1.
func TestDefensiveInstant_PreventionCapsAtIncoming(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.OasisRespiteRed{}, testutils.BluePitch{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 1}, nil, hand).Value; got != 1 {
		t.Fatalf("Value = %d, want 1 (Oasis Respite caps at IncomingDamage=1)", got)
	}
}

// Tests that a defensive instant with a real cost requires defense-phase pitch funding.
// Hand pitches Blue (3) into the defense phase; Peace of Mind costs 2 and prevents 4.
func TestDefensiveInstant_PeaceOfMindWithCost(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.PeaceOfMindRed{}, testutils.BluePitch{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value; got != 4 {
		t.Fatalf("Value = %d, want 4 (BluePitch funds 2 cost; Peace of Mind prevents 4)", got)
	}
}

// Tests that a defensive instant resolves alongside a printed DR — both go through the
// DR loop in defendersDamage and contribute their full prevention up to IncomingDamage.
func TestDefensiveInstant_StacksWithDR(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.BrushOffRed{}, cards.DodgeBlue{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 5}, nil, hand).Value; got != 5 {
		t.Fatalf("Value = %d, want 5 (Brush Off 3 + Dodge 2 fully prevent 5 incoming)", got)
	}
}

// Tests that a defensive instant in the arsenal slot can take Defend role.
func TestDefensiveInstant_DefendsFromArsenal(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{testutils.BluePitch{}}
	if got := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 3}, cards.BrushOffRed{}, hand).Value; got != 3 {
		t.Fatalf("Value = %d, want 3 (Brush Off plays from arsenal as defender)", got)
	}
}
