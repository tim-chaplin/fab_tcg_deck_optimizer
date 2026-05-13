package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that mode 0 fires the -1{p} self-debuff before crediting attack damage.
func TestBlusterBuff_Mode0DebuffsByOne(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	self := &card.CardState{Card: cards.BlusterBuffRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if s.Value() != 5 {
		t.Errorf("mode 0 Value = %d, want 5 (printed 6 - 1)", s.Value())
	}
}

// Tests that mode 1 keeps the printed power.
func TestBlusterBuff_Mode1KeepsPrintedPower(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	self := &card.CardState{Card: cards.BlusterBuffRed{}, Mode: 1}
	s.ResolveChainStep(s.Logger(), self)
	if s.Value() != 6 {
		t.Errorf("mode 1 Value = %d, want 6 (printed)", s.Value())
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
