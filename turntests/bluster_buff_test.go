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

// Tests that mode 0 fires the -1{p} self-debuff before crediting attack damage. With only
// 1 pitch resource available the runner can't afford Mode 1's extra {r}, so it falls back
// to the printed-mode swing minus the self-debuff.
func TestBlusterBuff_Mode0DebuffsByOne(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.BlusterBuffRed{}, testutils.FakeRedAttack()}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	if summary.Value != 5 {
		t.Errorf("Value = %d, want 5 (Mode 0: printed 6 - 1 self-debuff)", summary.Value)
	}
}

// Tests that mode 1 keeps the printed power when enough pitch is available to fund the
// extra {r}.
func TestBlusterBuff_Mode1KeepsPrintedPower(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.BlusterBuffRed{}, testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	if summary.Value != 6 {
		t.Errorf("Value = %d, want 6 (Mode 1: printed)", summary.Value)
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
