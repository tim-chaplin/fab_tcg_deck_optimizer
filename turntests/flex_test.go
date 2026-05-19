package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Mode 0 keeps the printed 4{p} — the optional cost is skipped.
func TestFlex_Mode0PrintedAttack(t *testing.T) {
	ge := gameengine.New()
	self := &card.CardState{Card: cards.FlexRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	if ge.Value() != 4 {
		t.Errorf("mode 0 Value = %d, want 4 (printed)", ge.Value())
	}
}

// Mode 1 pays {r}{r} for +2{p}, giving the full 6{p} swing.
func TestFlex_Mode1AddsTwo(t *testing.T) {
	ge := gameengine.New()
	self := &card.CardState{Card: cards.FlexRed{}, Mode: 1}
	ge.ResolveChainStep(ge.Logger(), self)
	if ge.Value() != 6 {
		t.Errorf("mode 1 Value = %d, want 6 (4 printed + 2 bonus)", ge.Value())
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
