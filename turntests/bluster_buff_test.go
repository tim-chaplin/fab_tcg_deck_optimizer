package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that mode 0 fires the -1{p} self-debuff before crediting attack damage.
func TestBlusterBuff_Mode0DebuffsByOne(t *testing.T) {
	ge := gameengine.New()
	self := &card.CardState{Card: cards.BlusterBuffRed{}}
	ge.ResolveChainStep(ge.Logger(), self)
	if ge.Value() != 5 {
		t.Errorf("mode 0 Value = %d, want 5 (printed 6 - 1)", ge.Value())
	}
}

// Tests that mode 1 keeps the printed power.
func TestBlusterBuff_Mode1KeepsPrintedPower(t *testing.T) {
	ge := gameengine.New()
	self := &card.CardState{Card: cards.BlusterBuffRed{}, Mode: 1}
	ge.ResolveChainStep(ge.Logger(), self)
	if ge.Value() != 6 {
		t.Errorf("mode 1 Value = %d, want 6 (printed)", ge.Value())
	}
}

// Tests that ModalCost reports printed-cost on mode 0 and printed+1 on mode 1.
func TestBlusterBuff_ModalCostsTrackMode(t *testing.T) {
	c := cards.BlusterBuffRed{}
	if got := c.ModalCost(0); got != 1 {
		t.Errorf("ModalCost(0) = %d, want 1", got)
	}
	if got := c.ModalCost(1); got != 2 {
		t.Errorf("ModalCost(1) = %d, want 2", got)
	}
}
