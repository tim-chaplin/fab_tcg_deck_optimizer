package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests the full activated-ability stack: blue pitch funds both the Gold ability
// ({2}) and a Reaping Blade swing (cost {1}, power 3). Spending the gold draws a
// card that post-hoc-promotes into arsenal; the swing credits 3 damage. Verifies
// Gold is consumed, arsenal carries the drawn card, and Value reflects the swing.
// Tests that an item created mid-chain is playable in the same chain. Hand =
// Flying High [Y], Strike Gold [R], blue pitch. Best line: Flying High (cost 0,
// go-again, grants go-again to Strike Gold) → Strike Gold [R] (cost 0, power 4 hits
// → Gold token) → spend Gold (cost {2}, paid by the blue pitch, has go-again, draws
// a card). Damage = 4 from Strike Gold; the drawn card promotes into arsenal post-
// chain so the next turn opens with a full hand + arsenal and zero gold tokens.
func TestGoldAbility_PopsItemCreatedSameTurn(t *testing.T) {
	deck := []sim.Card{
		// Five fillers for the gold-spend draw + next-turn 4-card deal.
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deck)
	hand := []sim.Card{
		cards.FlyingHighYellow{},
		cards.StrikeGoldRed{},
		testutils.BluePitch{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Strike Gold Red 4 power)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if state.Gold() != 0 {
		t.Fatalf("Gold at start of next turn = %d, want 0 (popped same-turn)", state.Gold())
	}
	if state.StartOfNextTurnArsenal == nil {
		t.Fatalf("StartOfNextTurnArsenal = nil, want the drawn card promoted")
	}
	if got := len(state.StartOfNextTurnHand); got != d.Hero.Intelligence() {
		t.Fatalf("StartOfNextTurnHand size = %d, want %d (full hand)", got, d.Hero.Intelligence())
	}
}

func TestGoldAbility_SpendsToFillArsenalAndSwings(t *testing.T) {
	deck := []sim.Card{
		// Five fillers covers the gold-spend draw plus next-turn's 4 dealt cards.
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, deck)
	hand := []sim.Card{testutils.BluePitch{}}
	priorItems := []sim.Item{sim.NewGoldItem(1)}
	got := sim.BestWithTriggers(heroes.Viserai{}, []sim.Weapon{weapons.ReapingBlade{}}, hand, sim.Matchup{IncomingDamage: 0}, deck, nil, nil, priorItems)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Reaping Blade swing power 3)", got.Value)
	}
	if got.State.Gold() != 0 {
		t.Fatalf("Gold after turn = %d, want 0 (the only token spent)", got.State.Gold())
	}
	if got.State.Arsenal == nil {
		t.Fatalf("Arsenal = nil, want the drawn card promoted into the slot")
	}
	if len(got.State.Deck) < d.Hero.Intelligence() {
		t.Fatalf("State.Deck has %d cards, want >= %d so next turn deals a full hand",
			len(got.State.Deck), d.Hero.Intelligence())
	}
	_ = d
}
