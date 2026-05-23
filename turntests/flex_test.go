package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// With no extra pitch the runner can't afford Mode 1's {r}{r} cost, so it falls back to
// Mode 0 — the printed 4{p} swing.
func TestFlex_Mode0PrintedAttack(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.FlexRed{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 4 {
		t.Errorf("Value = %d, want 4 (Mode 0 fallback when no pitch funds Mode 1)", summary.Value)
	}
}

// With 3 pitch available the runner pays {r}{r} for +2{p}, giving the full 6{p} swing.
func TestFlex_Mode1AddsTwo(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.FlexRed{}, testutils.BluePitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 6 {
		t.Errorf("Value = %d, want 6 (4 printed + 2 Mode 1 bonus)", summary.Value)
	}
}

// ModalCost reports 0 for the skip mode and {r}{r} = 2 for the buff mode.
func TestFlex_ModalCostsTrackMode(t *testing.T) {
	c := cards.FlexRed{}
	if got := c.ModalCost(0); got != 0 {
		t.Errorf("ModalCost(0) = %d, want 0", got)
	}
	if got := c.ModalCost(1); got != 2 {
		t.Errorf("ModalCost(1) = %d, want 2", got)
	}
}
