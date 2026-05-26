package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Mode 0 keeps the printed 2{p} — the optional cost is skipped.
func TestPunchAboveYourWeight_Mode0PrintedAttack(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PunchAboveYourWeightRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 2 {
		t.Errorf("mode 0 Value = %d, want 2 (printed)", ge.Value())
	}
}

// Mode 1 pays the {r}{r}{r} for +5{p}, giving the full 7{p} swing.
func TestPunchAboveYourWeight_Mode1AddsFive(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PunchAboveYourWeightRed{}, Mode: 1}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 7 {
		t.Errorf("mode 1 Value = %d, want 7 (2 printed + 5 bonus)", ge.Value())
	}
}

// ModalCost reports 0 for the skip mode and {r}{r}{r} = 3 for the buff mode.
func TestPunchAboveYourWeight_ModalCostsTrackMode(t *testing.T) {
	c := cards.PunchAboveYourWeightRed{}
	if got := c.ModalCost(0); got != 0 {
		t.Errorf("ModalCost(0) = %d, want 0", got)
	}
	if got := c.ModalCost(1); got != 3 {
		t.Errorf("ModalCost(1) = %d, want 3", got)
	}
}
